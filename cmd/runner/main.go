package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventgrid/publisher"
	"github.com/RMI/pacta/async"
	"github.com/RMI/pacta/azure/azblob"
	"github.com/RMI/pacta/azure/azcreds"
	"github.com/RMI/pacta/azure/azlog"
	"github.com/RMI/pacta/task"
	"github.com/namsral/flag"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapfield"
	"go.uber.org/zap/zapcore"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("args cannot be empty")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	var (
		env = fs.String("env", "", "The environment we're running in.")

		benchmarkDir = fs.String("benchmark_dir", "", "The path to the benchmark data for report generation")
		pactaDataDir = fs.String("pacta_data_dir", "", "The path to the PACTA data for report generation")

		azEventTopic    = fs.String("azure_event_topic", "", "The EventGrid topic to send notifications when tasks have finished")
		azTopicLocation = fs.String("azure_topic_location", "", "The location (like 'centralus-1') where our EventGrid topics are hosted")

		azStorageAccount  = fs.String("azure_storage_account", "", "The storage account to authenticate against for blob operations")
		azReportContainer = fs.String("azure_report_container", "", "The container in the storage account where we write generated portfolio reports to")

		minLogLevel zapcore.Level = zapcore.DebugLevel
	)
	fs.Var(&minLogLevel, "min_log_level", "If set, retains logs at the given level and above. Options: 'debug', 'info', 'warn', 'error', 'dpanic', 'panic', 'fatal' - default warn.")

	// Allows for passing in configuration via a -config path/to/env-file.conf
	// flag, see https://pkg.go.dev/github.com/namsral/flag#readme-usage
	fs.String(flag.DefaultConfigFlagname, "", "path to config file")
	if err := fs.Parse(args[1:]); err != nil {
		return fmt.Errorf("failed to parse flags: %v", err)
	}

	logger, err := azlog.New(&azlog.Config{
		Local:       *env == "local",
		MinLogLevel: minLogLevel,
	})
	if err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}
	defer logger.Sync()

	creds, credType, err := azcreds.New()
	if err != nil {
		return fmt.Errorf("failed to load Azure credentials: %w", err)
	}
	logger.Info("authenticated with Azure", zapfield.Str("credential_type", credType))

	pubsubClient, err := publisher.NewClient(fmt.Sprintf("https://%s.%s.eventgrid.azure.net/api/events", *azEventTopic, *azTopicLocation), creds, nil)
	if err != nil {
		return fmt.Errorf("failed to init pub/sub client: %w", err)
	}

	blobClient, err := azblob.NewClient(creds, *azStorageAccount)
	if err != nil {
		return fmt.Errorf("failed to init blob client: %w", err)
	}

	h, err := async.New(&async.Config{
		Blob:         blobClient,
		PubSub:       pubsubClient,
		Logger:       logger,
		BenchmarkDir: *benchmarkDir,
		PACTADataDir: *pactaDataDir,
	})
	if err != nil {
		return fmt.Errorf("failed to init async biz logic handler: %w", err)
	}

	validTasks := map[task.Type]func(context.Context, task.ID) error{
		task.CreateReport: toRunFn(async.LoadCreateReportRequestFromEnv, func(ctx context.Context, id task.ID, req *task.CreateReportRequest) error {
			return h.CreateReport(ctx, id, req, *azReportContainer)
		}),
		task.CreateAudit: toRunFn(async.LoadCreateAuditRequestFromEnv, h.CreateAudit),
	}

	taskID := task.ID(os.Getenv("TASK_ID"))
	if taskID == "" {
		return errors.New("no TASK_ID given")
	}

	taskType := task.Type(os.Getenv("TASK_TYPE"))
	if taskType == "" {
		return errors.New("no TASK_TYPE given")
	}

	taskFn, ok := validTasks[taskType]
	if !ok {
		return fmt.Errorf("unknown task type %q", taskType)
	}

	logger.Info("running PACTA task", zap.String("task_id", string(taskID)), zap.String("task_type", string(taskType)))

	if err := taskFn(ctx, taskID); err != nil {
		return fmt.Errorf("error running task: %w", err)
	}

	logger.Info("ran PACTA task successfully", zap.String("task_id", string(taskID)), zap.String("task_type", string(taskType)))

	return nil
}

func toRunFn[T any](reqFn func() (T, error), runFn func(context.Context, task.ID, T) error) func(context.Context, task.ID) error {
	return func(ctx context.Context, taskID task.ID) error {
		req, err := reqFn()
		if err != nil {
			return fmt.Errorf("failed to format request: %w", err)
		}
		return runFn(ctx, taskID, req)
	}
}

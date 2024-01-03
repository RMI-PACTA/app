package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventgrid/publisher"
	"github.com/RMI/pacta/azure/azblob"
	"github.com/RMI/pacta/azure/azlog"
	"github.com/RMI/pacta/blob"
	"github.com/RMI/pacta/pacta"
	"github.com/RMI/pacta/task"
	"github.com/google/uuid"
	"github.com/namsral/flag"
	"go.uber.org/zap"
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

		azEventParsePortfolioCompleteTopic = fs.String("azure_event_parse_portfolio_complete_topic", "", "The EventGrid topic to send notifications when parsing of portfolio(s) has finished")
		azTopicLocation                    = fs.String("azure_topic_location", "", "The location (like 'centralus-1') where our EventGrid topics are hosted")

		azStorageAccount           = fs.String("azure_storage_account", "", "The storage account to authenticate against for blob operations")
		azSourcePortfolioContainer = fs.String("azure_source_portfolio_container", "", "The container in the storage account where we read raw portfolios from")
		azDestPortfolioContainer   = fs.String("azure_dest_portfolio_container", "", "The container in the storage account where we read/write processed portfolios")
		azReportContainer          = fs.String("azure_report_container", "", "The container in the storage account where we write generated portfolio reports to")

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

	var creds azcore.TokenCredential
	// Azure has 3-4 ways to authenticate as an identity to their APIs (KMS, storage, etc).
	//   - When running locally, we use the "Environment" approach, which means we provide AZURE_* environment variables that authenticate against a local-only service account.
	//   - When running in Azure Container Apps Jobs, we use the "ManagedIdentitiy" approach, meaning we pull ambiently from the infrastructure we're running on (via a metadata service).
	// See https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#readme-credential-types for more info
	if azClientSecret := os.Getenv("AZURE_CLIENT_SECRET"); azClientSecret != "" {
		if creds, err = azidentity.NewEnvironmentCredential(nil); err != nil {
			return fmt.Errorf("failed to load Azure credentials from environment: %w", err)
		}
	} else {
		// We use "ManagedIdentity" instead of just "Default" because the default
		// timeout is too low in azidentity.NewDefaultAzureCredentials, so it times out
		// and fails to run.
		azClientID := os.Getenv("AZURE_CLIENT_ID")
		logger.Info("Loading user managed credentials", zap.String("client_id", azClientID))
		if creds, err = azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{
			ID: azidentity.ClientID(azClientID),
		}); err != nil {
			return fmt.Errorf("failed to load Azure credentials: %w", err)
		}
	}

	pubsubClient, err := publisher.NewClient(fmt.Sprintf("https://%s.%s.eventgrid.azure.net/api/events", *azEventParsePortfolioCompleteTopic, *azTopicLocation), creds, nil)
	if err != nil {
		return fmt.Errorf("failed to init pub/sub client: %w", err)
	}

	blobClient, err := azblob.NewClient(creds, *azStorageAccount)
	if err != nil {
		return fmt.Errorf("failed to init blob client: %w", err)
	}

	h := handler{
		blob:   blobClient,
		pubsub: pubsubClient,
		logger: logger,

		sourcePortfolioContainer: *azSourcePortfolioContainer,
		destPortfolioContainer:   *azDestPortfolioContainer,
		reportContainer:          *azReportContainer,
	}

	validTasks := map[task.Type]func(context.Context, task.ID) error{
		task.ParsePortfolio: toRunFn(parsePortfolioReq, h.parsePortfolio),
		task.CreateReport:   toRunFn(createReportReq, h.createReport),
		task.CreateAudit:    toRunFn(createAuditReq, h.createAudit),
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

type Blob interface {
	ReadBlob(ctx context.Context, uri string) (io.ReadCloser, error)
	WriteBlob(ctx context.Context, uri string, r io.Reader) error
	Scheme() blob.Scheme
}

type handler struct {
	blob   Blob
	pubsub *publisher.Client
	logger *zap.Logger

	sourcePortfolioContainer string
	destPortfolioContainer   string
	reportContainer          string
}

func parsePortfolioReq() (*task.ParsePortfolioRequest, error) {
	taskStr := os.Getenv("PARSE_PORTFOLIO_REQUEST")
	if taskStr == "" {
		return nil, errors.New("no PARSE_PORTFOLIO_REQUEST given")
	}
	var task task.ParsePortfolioRequest
	if err := json.NewDecoder(strings.NewReader(taskStr)).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to load ParsePortfolioRequest: %w", err)
	}
	return &task, nil
}

func (h *handler) uploadDirectory(ctx context.Context, dirPath, container string) error {
	base := filepath.Base(dirPath)

	return filepath.WalkDir(dirPath, func(path string, info fs.DirEntry, err error) error {
		if info.IsDir() {
			return nil
		}

		// This is a file, let's upload it to the container
		uri := blob.Join(h.blob.Scheme(), container, base, strings.TrimPrefix(path, dirPath))
		if err := h.uploadBlob(ctx, path, uri); err != nil {
			return fmt.Errorf("failed to upload blob: %w", err)
		}
		return nil
	})
}

func (h *handler) uploadBlob(ctx context.Context, srcPath, destURI string) error {
	srcF, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open file for upload: %w", err)
	}
	defer srcF.Close() // Best-effort in case something fails

	if err := h.blob.WriteBlob(ctx, destURI, srcF); err != nil {
		return fmt.Errorf("failed to write file to blob storage: %w", err)
	}

	if err := srcF.Close(); err != nil {
		return fmt.Errorf("failed to close source file: %w", err)
	}

	return nil
}

func (h *handler) downloadBlob(ctx context.Context, srcURI, destPath string) error {
	// Make sure the destination exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0600); err != nil {
		return fmt.Errorf("failed to create directory to download blob to: %w", err)
	}

	destF, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create dest file: %w", err)
	}
	defer destF.Close() // Best-effort in case something fails

	br, err := h.blob.ReadBlob(ctx, srcURI)
	if err != nil {
		return fmt.Errorf("failed to read raw portfolio: %w", err)
	}
	defer br.Close() // Best-effort in case something fails

	if _, err := io.Copy(destF, br); err != nil {
		return fmt.Errorf("failed to load raw portfolio: %w", err)
	}

	if err := br.Close(); err != nil {
		return fmt.Errorf("failed to close blob reader: %w", err)
	}

	if err := destF.Close(); err != nil {
		return fmt.Errorf("failed to close dest file: %w", err)
	}

	return nil
}

// TODO: Send a notification when parsing fails.
func (h *handler) parsePortfolio(ctx context.Context, taskID task.ID, req *task.ParsePortfolioRequest) error {
	// Load the portfolio from blob storage, place it in /mnt/raw_portfolios, where
	// the `process_portfolios.R` script expects it to be.
	for _, srcURI := range req.BlobURIs {
		id := uuid.New().String()
		// TODO: Probably set the CSV extension in the signed upload URL instead.
		destPath := filepath.Join("/", "mnt", "raw_portfolios", fmt.Sprintf("%s.csv", id))
		if err := h.downloadBlob(ctx, string(srcURI), destPath); err != nil {
			return fmt.Errorf("failed to download raw portfolio blob: %w", err)
		}
	}

	processedDir := filepath.Join("/", "mnt", "processed_portfolios")
	if err := os.MkdirAll(processedDir, 0600); err != nil {
		return fmt.Errorf("failed to create directory to download blob to: %w", err)
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "/usr/local/bin/Rscript", "/app/process_portfolios.R")
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run process_portfolios script: %w", err)
	}

	sc := bufio.NewScanner(&stderr)

	// TODO: Load from the output database file (or similar, like reading the processed_portfolios dir) instead of parsing stderr
	var paths []string
	for sc.Scan() {
		line := sc.Text()
		idx := strings.Index(line, "writing to file: " /* 17 chars */)
		if idx == -1 {
			continue
		}
		paths = append(paths, strings.TrimSpace(line[idx+17:]))
	}

	// NOTE: This code could benefit from some concurrency, but I'm opting not to prematurely optimize.
	var out []*task.ParsePortfolioResponseItem
	for _, p := range paths {
		lineCount, err := countCSVLines(p)
		if err != nil {
			return fmt.Errorf("failed to count lines in file %q: %w", p, err)
		}
		fileName := filepath.Base(p)
		blobURI := pacta.BlobURI(blob.Join(h.blob.Scheme(), h.destPortfolioContainer, fileName))
		if err := h.uploadBlob(ctx, p, string(blobURI)); err != nil {
			return fmt.Errorf("failed to copy parsed portfolio from %q to %q: %w", p, blobURI, err)
		}
		extension := filepath.Ext(fileName)
		fileType, err := pacta.ParseFileType(extension)
		if err != nil {
			return fmt.Errorf("failed to parse file type from file name %q: %w", fileName, err)
		}
		out = append(out, &task.ParsePortfolioResponseItem{
			Blob: pacta.Blob{
				FileName: fileName,
				FileType: fileType,
				BlobURI:  blobURI,
			},
			LineCount: lineCount,
		})
	}

	events := []publisher.Event{
		{
			Data: task.ParsePortfolioResponse{
				TaskID:  taskID,
				Request: req,
				Outputs: out,
			},
			DataVersion: to.Ptr("1.0"),
			EventType:   to.Ptr("parse-portfolio-complete"),
			EventTime:   to.Ptr(time.Now()),
			ID:          to.Ptr(string(taskID)),
			Subject:     to.Ptr("subject"),
		},
	}

	if _, err := h.pubsub.PublishEvents(ctx, events, nil); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	h.logger.Info("parsed portfolio", zap.String("task_id", string(taskID)))

	return nil
}

// TODO(grady): Move this line counting into the image to prevent having our code do any read of the actual underlying data.
func countCSVLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("opening file failed: %w", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("scanner.error returned: %w", err)
	}
	// Subtract 1 for the header row
	return lineCount - 1, nil
}

func createAuditReq() (*task.CreateAuditRequest, error) {
	car := os.Getenv("CREATE_AUDIT_REQUEST")
	if car == "" {
		return nil, errors.New("no CREATE_AUDIT_REQUEST was given")
	}
	var task task.CreateAuditRequest
	if err := json.NewDecoder(strings.NewReader(car)).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to load CreateAuditRequest: %w", err)
	}
	return &task, nil
}

func createReportReq() (*task.CreateReportRequest, error) {
	crr := os.Getenv("CREATE_REPORT_REQUEST")
	if crr == "" {
		return nil, errors.New("no CREATE_REPORT_REQUEST was given")
	}
	var task task.CreateReportRequest
	if err := json.NewDecoder(strings.NewReader(crr)).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to load CreateReportRequest: %w", err)
	}
	return &task, nil
}

func (h *handler) createAudit(ctx context.Context, taskID task.ID, req *task.CreateAuditRequest) error {
	return errors.New("not implemented")
}

func (h *handler) createReport(ctx context.Context, taskID task.ID, req *task.CreateReportRequest) error {
	fileNames := []string{}
	for i, blobURI := range req.BlobURIs {
		// Load the parsed portfolio from blob storage, place it in /mnt/
		// processed_portfolios, where the `create_report.R` script expects it
		// to be.
		fileName := fmt.Sprintf("%d.json", i)
		fileNames = append(fileNames, fileName)
		destPath := filepath.Join("/", "mnt", "processed_portfolios", fileName)
		if err := h.downloadBlob(ctx, string(blobURI), destPath); err != nil {
			return fmt.Errorf("failed to download processed portfolio blob: %w", err)
		}
	}

	reportDir := filepath.Join("/", "mnt", "reports")
	if err := os.MkdirAll(reportDir, 0600); err != nil {
		return fmt.Errorf("failed to create directory for reports to get copied to: %w", err)
	}

	cmd := exec.CommandContext(ctx, "/usr/local/bin/Rscript", "/app/create_report.R")
	cmd.Env = append(cmd.Env,
		"PORTFOLIO="+strings.Join(fileNames, ","),
		"HOME=/root", /* Required by pandoc */
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run pacta test CLI: %w", err)
	}

	// Download outputs from from /out and upload them to Azure
	dirEntries, err := os.ReadDir(reportDir)
	if err != nil {
		return fmt.Errorf("failed to read report directory: %w", err)
	}

	for _, dirEntry := range dirEntries {
		if !dirEntry.IsDir() {
			continue
		}
		dirPath := filepath.Join(reportDir, dirEntry.Name())
		if err := h.uploadDirectory(ctx, dirPath, h.reportContainer); err != nil {
			return fmt.Errorf("failed to upload report directory: %w", err)
		}
	}

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

type azureTokenCredential struct {
	accessToken string
}

func (a *azureTokenCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{
		Token: a.accessToken,
		// We just don't bother with expiration time
		ExpiresOn: time.Now().AddDate(1, 0, 0),
	}, nil
}

func asStrs[T ~string](in []T) []string {
	out := make([]string, len(in))
	for i, v := range in {
		out[i] = string(v)
	}
	return out
}

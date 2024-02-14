// Command parser is our shim for parsing + validating incoming PACTA porfolios,
// wrapping https://github.com/RMI-PACTA/workflow.portfolio.parsing
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
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

		azEventTopic    = fs.String("azure_event_topic", "", "The EventGrid topic to send notifications when tasks have finished")
		azTopicLocation = fs.String("azure_topic_location", "", "The location (like 'centralus-1') where our EventGrid topics are hosted")

		azStorageAccount         = fs.String("azure_storage_account", "", "The storage account to authenticate against for blob operations")
		azDestPortfolioContainer = fs.String("azure_dest_portfolio_container", "", "The container in the storage account where we write parsed portfolios")

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

	pubsubClient, err := publisher.NewClient(fmt.Sprintf("https://%s.%s.eventgrid.azure.net/api/events", *azEventTopic, *azTopicLocation), creds, nil)
	if err != nil {
		return fmt.Errorf("failed to init pub/sub client: %w", err)
	}

	blobClient, err := azblob.NewClient(creds, *azStorageAccount)
	if err != nil {
		return fmt.Errorf("failed to init blob client: %w", err)
	}

	h, err := async.New(&async.Config{
		Blob:   blobClient,
		PubSub: pubsubClient,
		Logger: logger,
	})
	if err != nil {
		return fmt.Errorf("failed to init async biz logic handler: %w", err)
	}

	taskID := task.ID(os.Getenv("TASK_ID"))
	if taskID == "" {
		return errors.New("no TASK_ID given")
	}

	req, err := async.LoadParsePortfolioRequestFromEnv()
	if err != nil {
		return fmt.Errorf("failed to parse portfolio request: %w", err)
	}

	logger.Info("running PACTA parsing task", zap.String("task_id", string(taskID)))

	if err := h.ParsePortfolio(ctx, taskID, req, *azDestPortfolioContainer); err != nil {
		return fmt.Errorf("error running task: %w", err)
	}

	logger.Info("ran PACTA parsing task successfully", zap.String("task_id", string(taskID)))

	return nil
}

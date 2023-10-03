package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/RMI/pacta/azure/azlog"
	"github.com/RMI/pacta/cmd/runner/taskrunner"
	"github.com/RMI/pacta/executors/local"
	"github.com/RMI/pacta/pacta"
	"github.com/RMI/pacta/task"
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

		minLogLevel zapcore.Level = zapcore.WarnLevel
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

	h := taskrunner.New(&local.Executor{}, logger)

	portfolioID := pacta.PortfolioID(os.Getenv("PORTFOLIO_ID"))
	logger.Info("running PACTA task", zap.String("portfolio_id", string(portfolioID)))
	if err := h.Execute(ctx, &task.StartRunRequest{
		PortfolioID: portfolioID,
	}); err != nil {
		return fmt.Errorf("failed to run task: %w", err)
	}

	return nil
}

// Package taskrunner implements the logic for preparing a portfolio for
// analysis, regardless of the underlying substrate we'll run the external
// processing logic on (e.g Docker or locally).
package taskrunner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/RMI/pacta/task"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Config struct {
	// ConfigPath should be a full path to a config file in the runner image,
	// like: /configs/{local,dev}.conf
	ConfigPath string

	// BaseImage is the runner image to execute, not specifying a tag.
	BaseImage *task.BaseImage

	Logger *zap.Logger

	Runner Runner
}

func (c *Config) validate() error {
	if c.ConfigPath == "" {
		return errors.New("no runner config path given")
	}

	if err := validateImage(c.BaseImage); err != nil {
		return fmt.Errorf("invalid base image: %w", err)
	}

	if c.Logger == nil {
		return errors.New("no logger given")
	}

	if c.Runner == nil {
		return errors.New("no runner given")
	}

	return nil
}

func validateImage(bi *task.BaseImage) error {
	if bi.Name == "" {
		return errors.New("no name given")
	}
	if bi.Registry == "" {
		return errors.New("no registry given")
	}
	return nil
}

type Runner interface {
	Run(ctx context.Context, cfg *task.Config) (task.RunnerID, error)
}

type TaskRunner struct {
	logger     *zap.Logger
	runner     Runner
	baseImage  *task.BaseImage
	configPath string
}

func New(cfg *Config) (*TaskRunner, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config given: %w", err)
	}

	return &TaskRunner{
		logger:     cfg.Logger,
		runner:     cfg.Runner,
		baseImage:  cfg.BaseImage,
		configPath: cfg.ConfigPath,
	}, nil
}

func (tr *TaskRunner) ParsePortfolio(ctx context.Context, req *task.ParsePortfolioRequest) (task.ID, task.RunnerID, error) {
	var taskBuffer bytes.Buffer
	if err := json.NewEncoder(&taskBuffer).Encode(req); err != nil {
		return "", "", fmt.Errorf("failed to encode ParsePortfolioRequest: %w", err)
	}
	tr.logger.Info("triggering parse portfolio task", zap.Any("req", req))
	return tr.run(ctx, []task.EnvVar{
		{
			Key:   "TASK_TYPE",
			Value: string(task.ParsePortfolio),
		},
		{
			Key:   "PARSE_PORTFOLIO_REQUEST",
			Value: taskBuffer.String(),
		},
	})
}

func (tr *TaskRunner) CreateReport(ctx context.Context, req *task.CreateReportRequest) (task.ID, task.RunnerID, error) {
	return tr.run(ctx, []task.EnvVar{
		{
			Key:   "TASK_TYPE",
			Value: string(task.CreateReport),
		},
		{
			Key:   "PORTFOLIO_ID",
			Value: string(req.PortfolioID),
		},
	})
}

func (tr *TaskRunner) run(ctx context.Context, env []task.EnvVar) (task.ID, task.RunnerID, error) {
	tr.logger.Info("triggering task run", zap.Any("env", env))
	taskID := uuid.NewString()
	runnerID, err := tr.runner.Run(ctx, &task.Config{
		Env: append(env, task.EnvVar{
			Key:   "TASK_ID",
			Value: taskID,
		}),
		Flags:   []string{"--config=" + tr.configPath},
		Command: []string{"/runner"},
		Image: &task.Image{
			Base: *tr.baseImage,
			// TODO: Take in the image digest as part of the task definition, as this can change per request.
			Tag: "latest",
		},
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to run task %q, %q: %w", taskID, runnerID, err)
	}
	return task.ID(taskID), runnerID, nil
}

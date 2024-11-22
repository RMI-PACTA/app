// Package taskrunner implements the logic for preparing a portfolio for
// analysis, regardless of the underlying substrate we'll run the external
// processing logic on (e.g Docker or locally).
//
// TODO: We use the tag "latest" throughout. For most cases, we'll want to
// version this and take in the image tag as part of the request.
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

	// RunnerImage is the runner image to execute, not specifying a tag.
	RunnerImage *task.BaseImage
	// RunnerImage is the parser image to execute, not specifying a tag.
	ParserImage *task.BaseImage

	Logger *zap.Logger

	Runner Runner
}

func (c *Config) validate() error {
	if c.ConfigPath == "" {
		return errors.New("no runner config path given")
	}

	if err := validateImage(c.RunnerImage); err != nil {
		return fmt.Errorf("invalid runner image: %w", err)
	}

	if err := validateImage(c.ParserImage); err != nil {
		return fmt.Errorf("invalid parser image: %w", err)
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
	logger      *zap.Logger
	runner      Runner
	runnerImage *task.BaseImage
	parserImage *task.BaseImage
	configPath  string
}

func New(cfg *Config) (*TaskRunner, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config given: %w", err)
	}

	return &TaskRunner{
		logger:      cfg.Logger,
		runner:      cfg.Runner,
		runnerImage: cfg.RunnerImage,
		parserImage: cfg.ParserImage,
		configPath:  cfg.ConfigPath,
	}, nil
}

type TaskRequest interface {
	task.ParsePortfolioRequest | task.CreateReportRequest | task.CreateAuditRequest
}

func encodeRequest[T TaskRequest](req *T) (string, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(req); err != nil {
		return "", fmt.Errorf("failed to encode request: %w", err)
	}
	value := buf.String()
	if len(value) > 128*1024 {
		return "", fmt.Errorf("request is too large: %d bytes > 128 kb", len(value))
	}
	return value, nil
}

func (tr *TaskRunner) ParsePortfolio(ctx context.Context, req *task.ParsePortfolioRequest) (task.ID, task.RunnerID, error) {
	value, err := encodeRequest(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to encode ParsePortfolioRequest: %w", err)
	}
	return tr.run(ctx, "/parser", withTag(tr.parserImage, "latest"), []task.EnvVar{
		{
			Key:   "TASK_TYPE",
			Value: string(task.ParsePortfolio),
		},
		{
			Key:   "PARSE_PORTFOLIO_REQUEST",
			Value: value,
		},
		// TODO(brandon): Unhardcode these
		{
			Key:   "BENCHMARK_DIR",
			Value: "/mnt/benchmark-data/65c1a416721b22a98c7925999ae03bc4",
		},
		{
			Key:   "PACTA_DATA_DIR",
			Value: "/mnt/pacta-data/2023Q4_20240718T150252Z",
		},
	})
}

func (tr *TaskRunner) CreateAudit(ctx context.Context, req *task.CreateAuditRequest) (task.ID, task.RunnerID, error) {
	value, err := encodeRequest(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to encode CreateAuditRequest: %w", err)
	}
	return tr.run(ctx, "/runner", withTag(tr.runnerImage, "latest"), []task.EnvVar{
		{
			Key:   "TASK_TYPE",
			Value: string(task.CreateAudit),
		},
		{
			Key:   "CREATE_AUDIT_REQUEST",
			Value: value,
		},
	})
}

func (tr *TaskRunner) CreateReport(ctx context.Context, req *task.CreateReportRequest) (task.ID, task.RunnerID, error) {
	value, err := encodeRequest(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to encode CreateReportRequest: %w", err)
	}
	return tr.run(ctx, "/runner", withTag(tr.runnerImage, "latest"), []task.EnvVar{
		{
			Key:   "TASK_TYPE",
			Value: string(task.CreateReport),
		},
		{
			Key:   "CREATE_REPORT_REQUEST",
			Value: value,
		},
		// TODO(brandon): Unhardcode these
		{
			Key:   "BENCHMARK_DIR",
			Value: "/mnt/benchmark-data/65c1a416721b22a98c7925999ae03bc4",
		},
		{
			Key:   "PACTA_DATA_DIR",
			Value: "/mnt/pacta-data/2023Q4_20240718T150252Z",
		},
	})
}

func withTag(img *task.BaseImage, tag string) *task.Image {
	return &task.Image{
		Base: *img,
		Tag:  tag,
	}
}

func (tr *TaskRunner) run(ctx context.Context, binary string, image *task.Image, env []task.EnvVar) (task.ID, task.RunnerID, error) {
	tr.logger.Info("triggering task run", zap.Any("env", env))
	taskID := uuid.NewString()
	runnerID, err := tr.runner.Run(ctx, &task.Config{
		Env: append(env, task.EnvVar{
			Key:   "TASK_ID",
			Value: taskID,
		}),
		Flags:   []string{"--config=" + tr.configPath},
		Command: []string{binary},
		Image:   image,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to run task %q, %q: %w", taskID, runnerID, err)
	}
	return task.ID(taskID), runnerID, nil
}

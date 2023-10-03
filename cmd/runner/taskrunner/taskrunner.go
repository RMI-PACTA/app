// Package taskrunner implements the logic for preparing a portfolio for
// analysis, regardless of the underlying substrate we'll run the external
// processing logic on (e.g Docker or locally).
package taskrunner

import (
	"context"
	"fmt"

	"github.com/RMI/pacta/pacta"
	"github.com/RMI/pacta/task"
	"go.uber.org/zap"
)

type Executor interface {
	// TODO: Update this interface to include relevant inputs and outputs
	ProcessPortfolio(ctx context.Context) (*pacta.PortfolioResult, error)
}

type Handler struct {
	logger   *zap.Logger
	executor Executor
}

func New(executor Executor, logger *zap.Logger) *Handler {
	return &Handler{
		logger:   logger,
		executor: executor,
	}
}

func (h *Handler) Execute(ctx context.Context, req *task.StartRunRequest) error {
	// TODO: Add logic for loading portfolio blobs and putting them in the right places.

	res, err := h.executor.ProcessPortfolio(ctx)
	if err != nil {
		return fmt.Errorf("failed to process portfolio: %w", err)
	}

	h.logger.Info("processed portfolio, result", zap.Any("result", res))

	return nil
}

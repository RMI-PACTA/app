// Package local implements an executor that runs PACTA processing directly
// on the current machine via os/exec. This implementation is used in deployed
// environments, where the runner binary will be baked into a Docker image
// with the portfolio processing code, and so has local access to a functioning
// portfolio-handling environment.
package local

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"github.com/RMI/pacta/pacta"
)

type Executor struct{}

func (e *Executor) ProcessPortfolio(ctx context.Context) (*pacta.PortfolioResult, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "echo", "TODO, the command to execute")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run pacta test CLI: %w", err)
	}

	return &pacta.PortfolioResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}, nil
}

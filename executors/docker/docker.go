// Package docker implements an executor that runs PACTA processing in a
// container against a local Docker daemon. This implementation is used for
// local testing, where one likely doesn't have a functional portfolio-handling
// environment on the host machine.
package docker

import (
	"bytes"
	"context"
	"fmt"

	"github.com/RMI/pacta/pacta"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"go.uber.org/zap"

	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

type Executor struct {
	client *client.Client
	logger *zap.Logger
}

func NewExecutor(dockerClient *client.Client, logger *zap.Logger) *Executor {
	return &Executor{client: dockerClient, logger: logger}
}

func (e *Executor) ProcessPortfolio(ctx context.Context) (*pacta.PortfolioResult, error) {
	cfg := &container.Config{
		// Use our runner image, which should contain PACTA processing infra.
		Image: "rmipacta.azurecr.io/runner:latest",

		// Run the script, tell it to output data to our mounted location.
		Cmd: []string{"echo", "TODO, the command to execute"},

		AttachStdin: false,
		Tty:         false,
	}

	hostCfg := &container.HostConfig{
		// TODO: Add relevant inputs + outputs, likely as mounts
	}
	platform := &specs.Platform{
		Architecture: "amd64",
		OS:           "linux",
	}

	resp, err := e.client.ContainerCreate(ctx, cfg, hostCfg, nil /* net config */, platform, "" /* random name */)
	if err != nil {
		return nil, fmt.Errorf("failed to create PACTA container: %w", err)
	}
	defer func() {
		err := e.client.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			RemoveLinks:   false,
			Force:         true,
		})
		if err != nil {
			e.logger.Error("failed to clean up container",
				zap.String("container_id", resp.ID),
				zap.Error(err))
		}
	}()

	if err := e.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, fmt.Errorf("failed to start PACTA container: %w", err)
	}

	waitC, errC := e.client.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	select {
	case resp := <-waitC:
		if resp.Error != nil {
			return nil, fmt.Errorf("error in container wait response: %+v", resp.Error)
		}
	case err := <-errC:
		return nil, fmt.Errorf("error while waiting for container to complete: %w", err)
	}

	// If we're here, container exited successfully, let's load the logs.
	logRC, err := e.client.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       "all",
		Details:    true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read container logs: %w", err)
	}
	defer func() {
		if err := logRC.Close(); err != nil {
			e.logger.Warn("failed to close contrainer log reader",
				zap.String("container_id", resp.ID),
				zap.Error(err))
		}
	}()

	var stdout, stderr bytes.Buffer
	if _, err := stdcopy.StdCopy(&stdout, &stderr, logRC); err != nil {
		return nil, fmt.Errorf("failed to read logs: %w", err)
	}

	return &pacta.PortfolioResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}, nil
}

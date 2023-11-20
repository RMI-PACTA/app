// Package dockertask implements PACTA processing in a container against a
// local Docker daemon. This implementation is used for local testing, where
// one likely doesn't have a functional portfolio-handling environment on the
// host machine.
package dockertask

import (
	"context"
	"fmt"

	"github.com/RMI/pacta/task"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"go.uber.org/zap"

	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

type Runner struct {
	client *client.Client
	logger *zap.Logger
	sp     *ServicePrincipal
}

type ServicePrincipal struct {
	TenantID     string
	ClientID     string
	ClientSecret string
}

func NewRunner(logger *zap.Logger, sp *ServicePrincipal) (*Runner, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Docker client: %w", err)
	}

	return &Runner{client: cli, logger: logger, sp: sp}, nil
}

func (r *Runner) Run(ctx context.Context, taskCfg *task.Config) (task.RunnerID, error) {
	r.logger.Info("Starting task run", zap.Any("task_config", taskCfg))

	env := []string{
		"AZURE_TENANT_ID=" + r.sp.TenantID,
		"AZURE_CLIENT_ID=" + r.sp.ClientID,
		"AZURE_CLIENT_SECRET=" + r.sp.ClientSecret,
	}
	for _, e := range taskCfg.Env {
		env = append(env, e.Key+"="+e.Value)
	}
	cfg := &container.Config{
		Image: taskCfg.Image.String(),

		// Run the script, tell it to output data to our mounted location.
		Cmd:        taskCfg.Flags,
		Entrypoint: taskCfg.Command,

		Env: env,

		AttachStdin: false,
		Tty:         false,
	}

	platform := &specs.Platform{
		Architecture: "amd64",
		OS:           "linux",
	}

	hostCfg := &container.HostConfig{
		// AutoRemove: true,
	}

	resp, err := r.client.ContainerCreate(ctx, cfg, hostCfg, nil /* net config */, platform, "" /* random name */)
	if err != nil {
		return "", fmt.Errorf("failed to create PACTA container: %w", err)
	}

	// ERROR IS HERE
	if err := r.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start PACTA container: %w", err)
	}

	// We don't wait for the container to exit, it's "fire and forget". We'll get an
	// Event Grid webhook when it's done.

	// waitC, errC := r.client.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

	// select {
	// case resp := <-waitC:
	// 	if resp.Error != nil {
	// 		return "", fmt.Errorf("error in container wait response: %+v", resp.Error)
	// 	}
	// case err := <-errC:
	// 	return "", fmt.Errorf("error while waiting for container to complete: %w", err)
	// }

	return task.RunnerID(resp.ID), nil
}

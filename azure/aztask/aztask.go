// Package aztask wraps Azure's Container Apps Jobs API to provide basic async
// task execution for the PACTA ecosystem.
package aztask

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	armappcontainers "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appcontainers/armappcontainers/v2"
	"github.com/RMI/pacta/task"
	"github.com/Silicon-Ally/idgen"
)

type Runner struct {
	client *armappcontainers.JobsClient

	cfg *Config
	gen *idgen.Generator
}

type Config struct {
	// The Azure Subscription to issue API calls against, like aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee
	SubscriptionID string
	// The resource group where our Container Apps Job is located, like rmi-pacta-{local,dev}
	ResourceGroup string
	// The Azure client ID of the managed identity that the job should run as, like ffffffff-0000-1111-2222-333333333333
	ManagedIdentityClientID string
	// The name of the Container Apps Job to start an execution of, like pacta-runner
	JobName string
}

func (c *Config) validate() error {
	if c.SubscriptionID == "" {
		return errors.New("no identity subscription ID given")
	}
	if c.ResourceGroup == "" {
		return errors.New("no identity resource group given")
	}
	if c.ManagedIdentityClientID == "" {
		return errors.New("no identity client ID given")
	}

	return nil
}

func NewRunner(creds azcore.TokenCredential, cfg *Config) (*Runner, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid task runner config: %w", err)
	}

	clientFactory, err := armappcontainers.NewClientFactory(cfg.SubscriptionID, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return &Runner{
		client: clientFactory.NewJobsClient(),
		cfg:    cfg,
	}, nil
}

func (r *Runner) Run(ctx context.Context, cfg *task.Config) (task.RunnerID, error) {
	envVars := []*armappcontainers.EnvironmentVar{
		{
			Name:  to.Ptr("AZURE_CLIENT_ID"),
			Value: to.Ptr(r.cfg.ManagedIdentityClientID),
		},
		{
			Name:  to.Ptr("MANAGED_IDENTITY_CLIENT_ID"),
			Value: to.Ptr(r.cfg.ManagedIdentityClientID),
		},
	}
	for _, v := range cfg.Env {
		envVars = append(envVars, &armappcontainers.EnvironmentVar{
			Name:  to.Ptr(v.Key),
			Value: to.Ptr(v.Value),
		})
	}

	poller, err := r.client.BeginStart(ctx, r.cfg.ResourceGroup, r.cfg.JobName, &armappcontainers.JobsClientBeginStartOptions{
		Template: &armappcontainers.JobExecutionTemplate{
			Containers: []*armappcontainers.JobExecutionContainer{{
				Args:    toPtrs(cfg.Flags),
				Command: toPtrs(cfg.Command),
				Env:     envVars,
				Image:   to.Ptr(cfg.Image.String()),
				Name:    to.Ptr("pacta-runner"),
				Resources: &armappcontainers.ContainerResources{
					CPU:    to.Ptr(1.0),
					Memory: to.Ptr("2Gi"),
				},
			}},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to start container app job: %w", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to poll for container app start: %w", err)
	}

	return task.RunnerID(*res.ID), nil
}

func toPtrs[T any](in []T) []*T {
	if in == nil {
		return nil
	}
	out := make([]*T, len(in))
	for i, v := range in {
		out[i] = &v
	}
	return out
}

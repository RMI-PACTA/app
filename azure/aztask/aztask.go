// Package aztask wraps Azure's Container Apps Jobs API to provide basic async
// task execution for the PACTA ecosystem.
package aztask

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

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
	// Location is the location to run the runner, like centralus
	Location string

	// Identity is the account the runner should act as.
	Identity *RunnerIdentity

	Rand *rand.Rand
}

func (c *Config) validate() error {
	if c.Location == "" {
		return errors.New("no container location given")
	}

	if err := c.Identity.validate(); err != nil {
		return fmt.Errorf("invalid identity config: %w", err)
	}

	if c.Rand == nil {
		return errors.New("no random number generator given")
	}

	return nil
}

type RunnerIdentity struct {
	// Like runner-local
	Name string
	// Like aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee
	SubscriptionID string
	// Like rmi-pacta-{local,dev}
	ResourceGroup string
	// Like ffffffff-0000-1111-2222-333333333333
	ClientID string
	// Like pacta-{local,dev}, the name of the Container Apps Environment
	ManagedEnvironment string
}

func (ri *RunnerIdentity) validate() error {
	if ri.Name == "" {
		return errors.New("no identity name given")
	}
	if ri.SubscriptionID == "" {
		return errors.New("no identity subscription ID given")
	}
	if ri.ResourceGroup == "" {
		return errors.New("no identity resource group given")
	}
	if ri.ClientID == "" {
		return errors.New("no identity client ID given")
	}
	return nil
}

func (r *RunnerIdentity) String() string {
	tmpl := "/subscriptions/%s/resourcegroups/%s/providers/Microsoft.ManagedIdentity/userAssignedIdentities/%s"
	return fmt.Sprintf(tmpl, r.SubscriptionID, r.ResourceGroup, r.Name)
}

func (r *RunnerIdentity) EnvironmentID() string {
	tmpl := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/managedEnvironments/%s"
	return fmt.Sprintf(tmpl, r.SubscriptionID, r.ResourceGroup, r.ManagedEnvironment)
}

func NewRunner(creds azcore.TokenCredential, cfg *Config) (*Runner, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid task runner config: %w", err)
	}

	clientFactory, err := armappcontainers.NewClientFactory(cfg.Identity.SubscriptionID, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	gen, err := idgen.New(cfg.Rand, idgen.WithDefaultLength(32), idgen.WithCharSet([]rune("abcdefghijklmnopqrstuvwxyz")))
	if err != nil {
		return nil, fmt.Errorf("failed to init ID generator: %w", err)
	}

	return &Runner{
		client: clientFactory.NewJobsClient(),
		cfg:    cfg,
		gen:    gen,
	}, nil
}

func (r *Runner) Run(ctx context.Context, cfg *task.Config) (task.ID, error) {

	name := r.gen.NewID()
	identity := r.cfg.Identity.String()
	envID := r.cfg.Identity.EnvironmentID()

	envVars := []*armappcontainers.EnvironmentVar{
		{
			Name:  to.Ptr("AZURE_CLIENT_ID"),
			Value: to.Ptr(r.cfg.Identity.ClientID),
		},
		{
			Name:  to.Ptr("MANAGED_IDENTITY_CLIENT_ID"),
			Value: to.Ptr(r.cfg.Identity.ClientID),
		},
	}
	for _, v := range cfg.Env {
		envVars = append(envVars, &armappcontainers.EnvironmentVar{
			Name:  to.Ptr(v.Key),
			Value: to.Ptr(v.Value),
		})
	}

	job := armappcontainers.Job{
		Location: &r.cfg.Location,
		Identity: &armappcontainers.ManagedServiceIdentity{
			Type: to.Ptr(armappcontainers.ManagedServiceIdentityTypeUserAssigned),
			UserAssignedIdentities: map[string]*armappcontainers.UserAssignedIdentity{
				identity: {},
			},
		},
		Properties: &armappcontainers.JobProperties{
			Configuration: &armappcontainers.JobConfiguration{
				ReplicaTimeout: to.Ptr(int32(60 * 60 * 2 /* two hours */)),
				TriggerType:    to.Ptr(armappcontainers.TriggerTypeManual),
				ManualTriggerConfig: &armappcontainers.JobConfigurationManualTriggerConfig{
					// Run one copy.
					Parallelism:            to.Ptr(int32(1)),
					ReplicaCompletionCount: to.Ptr(int32(1)),
				},
				// Don't retry, if it failed once, it'll probably fail again. We might relax
				// this in the future if we identify "transient" errors.
				ReplicaRetryLimit: to.Ptr(int32(0)),
				Registries: []*armappcontainers.RegistryCredentials{
					{
						Server:   to.Ptr(cfg.Image.Base.Registry),
						Identity: to.Ptr(identity),
					},
				},
				Secrets: []*armappcontainers.Secret{
					// TODO: Put any useful configuration here.
				},
			},
			EnvironmentID: to.Ptr(envID),
			Template: &armappcontainers.JobTemplate{
				Containers: []*armappcontainers.Container{
					{
						Args:    toPtrs(cfg.Flags),
						Command: toPtrs(cfg.Command),
						Env:     envVars,
						Image:   to.Ptr(cfg.Image.String()),
						Name:    to.Ptr(name),
						Probes:  []*armappcontainers.ContainerAppProbe{},
						Resources: &armappcontainers.ContainerResources{
							CPU:    to.Ptr(1.0),
							Memory: to.Ptr("2Gi"),
						},
						VolumeMounts: []*armappcontainers.VolumeMount{},
					},
				},
				Volumes: []*armappcontainers.Volume{
					// TODO: Mount any sources here.
				},
			},
		},
		Tags: map[string]*string{},
	}
	poller, err := r.client.BeginCreateOrUpdate(ctx, r.cfg.Identity.ResourceGroup, name, job, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create container app job: %w", err)
	}

	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to poll container group creation: %w", err)
	}

	poller2, err := r.client.BeginStart(ctx, r.cfg.Identity.ResourceGroup, name, nil)
	if err != nil {
		return "", fmt.Errorf("failed to start container app job: %w", err)
	}
	if _, err := poller2.PollUntilDone(ctx, nil); err != nil {
		return "", fmt.Errorf("failed to poll for container app start: %w", err)
	}

	return task.ID(*res.ID), nil
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

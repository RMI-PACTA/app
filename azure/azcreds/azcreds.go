// Package azcreds provides helpers for getting environment-appropriate credentials.
//
// The context is that Azure has 3-4 ways to authenticate as an identity to
// their APIs (KMS, storage, etc):
//
//   - When running locally, we use the "Environment" approach, which means we provide AZURE_* environment variables that authenticate against a local-only service account.
//   - When running in Azure Container Apps Jobs, we use the "ManagedIdentitiy" approach, meaning we pull ambiently from the infrastructure we're running on (via a metadata service).
//
// See https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#readme-credential-types for more info
package azcreds

import (
	"errors"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type Type string

const (
	Environment     = Type("ENVIRONMENT")
	ManagedIdentity = Type("MANAGED_IDENTITY")
)

// New returns appropriate credentials for the environment we're running in.
func New() (azcore.TokenCredential, Type, error) {
	if azClientSecret := os.Getenv("AZURE_CLIENT_SECRET"); azClientSecret != "" {
		return newEnvCreds()
	}

	// We use "ManagedIdentity" instead of just "Default" because the default
	// timeout is too low in azidentity.NewDefaultAzureCredentials, so it times out
	// and fails to run.
	azClientID := os.Getenv("AZURE_CLIENT_ID")
	if azClientID == "" {
		return nil, "", errors.New("no AZURE_CLIENT_SECERT or AZURE_CLIENT_ID found")
	}

	return newManagedIdentityCreds(azClientID)
}

func newEnvCreds() (*azidentity.EnvironmentCredential, Type, error) {
	creds, err := azidentity.NewEnvironmentCredential(nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to load Azure credentials from environment: %w", err)
	}
	return creds, Environment, nil

}

func newManagedIdentityCreds(azClientID string) (*azidentity.ManagedIdentityCredential, Type, error) {
	creds, err := azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{
		ID: azidentity.ClientID(azClientID),
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to load Azure credentials: %w", err)
	}
	return creds, ManagedIdentity, nil
}

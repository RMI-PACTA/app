// Package azblob wraps the existing Azure blob library to provide basic upload,
// download, and URL signing functionality against a standardized interface.
package azblob

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	azservice "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/RMI/pacta/blob"
)

const (
	Scheme = blob.Scheme("az")
)

type Client struct {
	storageAccount string
	now            func() time.Time

	client    *azblob.Client
	svcClient *azservice.Client

	cachedUDCMu     *sync.Mutex
	cachedUDC       *azservice.UserDelegationCredential
	cachedUDCExpiry time.Time
}

func NewClient(creds azcore.TokenCredential, storageAccount string) (*Client, error) {
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", storageAccount)

	client, err := azblob.NewClient(serviceURL, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to init Azure blob client: %w", err)
	}

	svcClient, err := azservice.NewClient(serviceURL, creds, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to init Azure blob service client: %w", err)
	}

	return &Client{
		storageAccount: storageAccount,
		now:            func() time.Time { return time.Now().UTC() },

		client:    client,
		svcClient: svcClient,

		cachedUDCMu: &sync.Mutex{},
	}, nil
}

func (c *Client) Scheme() blob.Scheme {
	return Scheme
}

func (c *Client) WriteBlob(ctx context.Context, uri string, r io.Reader) error {
	ctr, blb, ok := blob.SplitURI(Scheme, uri)
	if !ok {
		return fmt.Errorf("malformed URI %q is not for Azure", uri)
	}

	if _, err := c.client.UploadStream(ctx, ctr, blb, r, nil); err != nil {
		return fmt.Errorf("failed to upload blob: %w", err)
	}
	return nil
}

func (c *Client) ReadBlob(ctx context.Context, uri string) (io.ReadCloser, error) {
	ctr, blb, ok := blob.SplitURI(Scheme, uri)
	if !ok {
		return nil, fmt.Errorf("malformed URI %q is not for Azure", uri)
	}

	resp, err := c.client.DownloadStream(ctx, ctr, blb, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to read blob: %w", err)
	}

	return resp.Body, nil
}

func (c *Client) DeleteBlob(ctx context.Context, uri string) error {
	ctr, blb, ok := blob.SplitURI(Scheme, uri)
	if !ok {
		return fmt.Errorf("malformed URI %q is not for Azure", uri)
	}

	_, err := c.client.DeleteBlob(ctx, ctr, blb, nil)
	if err != nil {
		return fmt.Errorf("failed to delete blob: %w", err)
	}

	return nil
}

// SignedUploadURL returns a URL that is allowed to upload to the given URI.
// See https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/storage/azblob@v1.0.0/sas#example-package-UserDelegationSAS
func (c *Client) SignedUploadURL(ctx context.Context, uri string) (string, time.Time, error) {
	return c.signBlob(ctx, uri, &sas.BlobPermissions{Create: true, Write: true})
}

// SignedDownloadURL returns a URL that is allowed to download the file at the given URI.
// See https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/storage/azblob@v1.0.0/sas#example-package-UserDelegationSAS
func (c *Client) SignedDownloadURL(ctx context.Context, uri string) (string, time.Time, error) {
	return c.signBlob(ctx, uri, &sas.BlobPermissions{Read: true})
}

func (c *Client) signBlob(ctx context.Context, uri string, perms *sas.BlobPermissions) (string, time.Time, error) {
	ctr, blb, ok := blob.SplitURI(Scheme, uri)
	if !ok {
		return "", time.Time{}, fmt.Errorf("malformed URI %q is not for Azure", uri)
	}

	// The blob component is important, otherwise the signed URL is applicable to the whole container.
	if blb == "" {
		return "", time.Time{}, fmt.Errorf("uri %q did not contain a blob component", uri)
	}

	now := c.now().UTC().Add(-10 * time.Second)
	udc, err := c.getUserDelegationCredential(ctx, now)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get udc: %w", err)
	}

	expiry := now.Add(15 * time.Minute)

	// Create Blob Signature Values with desired permissions and sign with user delegation credential
	sasQueryParams, err := sas.BlobSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		StartTime:     now,
		ExpiryTime:    expiry,
		Permissions:   perms.String(),
		ContainerName: ctr,
		BlobName:      blb,
	}.SignWithUserDelegation(udc)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign blob: %w", err)
	}

	return fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s", c.storageAccount, ctr, blb, sasQueryParams.Encode()), expiry, nil
}

func (c *Client) ListBlobs(ctx context.Context, uriPrefix string) ([]string, error) {
	ctr, blobPrefix, ok := blob.SplitURI(Scheme, uriPrefix)
	if !ok {
		return nil, fmt.Errorf("malformed URI prefix %q is not for Azure", uriPrefix)
	}

	if blobPrefix == "" {
		return nil, fmt.Errorf("uri prefix %q did not contain a blob component", uriPrefix)
	}

	pager := c.client.NewListBlobsFlatPager(ctr, &azblob.ListBlobsFlatOptions{
		Prefix: &blobPrefix,
	})

	var blobs []string
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to load page of blobs: %w", err)
		}
		for _, bi := range resp.Segment.BlobItems {
			blobs = append(blobs, blob.Join(Scheme, ctr, *bi.Name))
		}
	}

	return blobs, nil
}

func (c *Client) getUserDelegationCredential(ctx context.Context, now time.Time) (*azservice.UserDelegationCredential, error) {
	c.cachedUDCMu.Lock()
	defer c.cachedUDCMu.Unlock()

	expiry := now.Add(48 * time.Hour)
	info := azservice.KeyInfo{
		Start:  to.Ptr(now.UTC().Format(sas.TimeFormat)),
		Expiry: to.Ptr(expiry.UTC().Format(sas.TimeFormat)),
	}

	if !c.cachedUDCExpiry.IsZero() && c.cachedUDCExpiry.Sub(now) > 1*time.Minute {
		return c.cachedUDC, nil
	}

	udc, err := c.svcClient.GetUserDelegationCredential(ctx, info, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get delegated credentials: %w", err)
	}
	c.cachedUDC = udc
	c.cachedUDCExpiry = expiry

	return udc, nil
}

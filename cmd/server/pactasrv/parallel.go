package pactasrv

import (
	"context"
	"fmt"

	"github.com/RMI/pacta/oapierr"
	"github.com/RMI/pacta/pacta"
	"go.uber.org/zap"
)

type blobDeleter interface {
	DeleteBlob(ctx context.Context, uri string) error
}

func deleteBlobs(ctx context.Context, bd blobDeleter, uris []string) error {
	// Implement parallel delete if slow - not prematurely optimizing.
	for i, uri := range uris {
		if err := bd.DeleteBlob(ctx, uri); err != nil {
			return oapierr.Internal("failed to delete blob",
				zap.String("uri", uri),
				zap.String("number_in_order", fmt.Sprintf("%d/%d", i+1, len(uris))),
				zap.Error(err))
		}
	}
	return nil
}

func (s *Server) deleteBlobs(ctx context.Context, uris ...pacta.BlobURI) error {
	return deleteBlobs(ctx, s.Blob, asStrs(uris))
}

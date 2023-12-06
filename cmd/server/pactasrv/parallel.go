package pactasrv

import (
	"context"
	"fmt"
)

type blobDeleter interface {
	DeleteBlob(ctx context.Context, uri string) error
}

func deleteBlobs(ctx context.Context, bd blobDeleter, uris []string) error {
	// Implement parallel delete if slow - not prematurely optimizing.
	for i, uri := range uris {
		if err := bd.DeleteBlob(ctx, uri); err != nil {
			return fmt.Errorf("deleting blob %d/%d: %w", i, len(uris), err)
		}
	}
	return nil
}

func (s *Server) deleteBlobs(ctx context.Context, uris []string) error {
	return deleteBlobs(ctx, s.Blob, uris)
}

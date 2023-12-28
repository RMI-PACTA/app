package pactasrv

import (
	"context"
	"time"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/oapierr"
	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
	"go.uber.org/zap"
)

func (s *Server) AccessBlobContent(ctx context.Context, request api.AccessBlobContentReq) (api.AccessBlobContentResp, error) {
	ownerID, err := s.getUserOwnerID(ctx)
	if err != nil {
		return nil, err
	}

	blobIDs := []pacta.BlobID{}
	for _, item := range request.Body.Items {
		blobIDs = append(blobIDs, pacta.BlobID(item.BlobID))
	}
	err404 := oapierr.NotFound("blob not found", zap.Strings("blob_ids", asStrs(blobIDs)))
	bois, err := s.DB.BlobOwners(s.DB.NoTxn(ctx), blobIDs)
	if err != nil {
		if db.IsNotFound(err) {
			return nil, err404
		}
		return nil, oapierr.Internal("error getting blob owners", zap.Error(err), zap.Strings("blob_ids", asStrs(request.BlobIDs)))
	}
	asMap := map[pacta.BlobID]*pacta.BlobOwnerInformation{}
	for _, boi := range bois {
		asMap[boi.BlobID] = boi
	}
	for _, blobID := range blobIDs {
		boi := asMap[blobID]
		if boi.OwnerID != ownerID {
			// TODO(grady) do AdminDebugEnabled & IsAdmin check here.
			return nil, err404
		}
	}

	blobs, err := s.DB.Blobs(s.DB.NoTxn(ctx), blobIDs)
	if err != nil {
		if db.IsNotFound(err) {
			return nil, err404
		}
		return nil, oapierr.Internal("error getting blobs", zap.Error(err), zap.Strings("blob_ids", asStrs(request.BlobIDs)))
	}

	// TODO(grady) - Add Audit Logs here

	response := api.AccessBlobContentResp{}
	for _, blob := range blobs {
		url, err := s.Blob.SignedDownloadURL(ctx, string(blob.BlobURI))
		if err != nil {
			return nil, oapierr.Internal("error getting signed download url", zap.Error(err), zap.String("blob_uri", string(blob.BlobURI)))
		}
		response.Items = append(response.Items, api.AccessBlobContentRespItems{
			BlobID:         blob.ID,
			URL:            url,
			ExpirationTime: s.Now().Add(15 * time.Minute), // See TODO in azblob.go
		})
	}
	return response, nil
}

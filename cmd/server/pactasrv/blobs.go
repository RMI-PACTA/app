package pactasrv

import (
	"context"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/oapierr"
	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
	"go.uber.org/zap"
)

func (s *Server) AccessBlobContent(ctx context.Context, request api.AccessBlobContentRequestObject) (api.AccessBlobContentResponseObject, error) {
	actorInfo, err := s.getactorInfoOrErrIfAnon(ctx)
	if err != nil {
		return nil, err
	}

	blobIDs := []pacta.BlobID{}
	for _, item := range request.Body.Items {
		blobIDs = append(blobIDs, pacta.BlobID(item.BlobId))
	}
	err404 := oapierr.NotFound("blob not found", zap.Strings("blob_ids", asStrs(blobIDs)))
	bos, err := s.DB.BlobContexts(s.DB.NoTxn(ctx), blobIDs)
	if err != nil {
		if db.IsNotFound(err) {
			return nil, err404
		}
		return nil, oapierr.Internal("error getting blob owners", zap.Error(err), zap.Strings("blob_ids", asStrs(blobIDs)))
	}
	asMap := map[pacta.BlobID]*pacta.BlobContext{}
	for _, boi := range bos {
		asMap[boi.BlobID] = boi
	}
	auditLogs := []*pacta.AuditLog{}
	for _, blobID := range blobIDs {
		boi := asMap[blobID]
		accessAsOwner := boi.PrimaryTargetOwnerID == actorInfo.OwnerID
		accessAsAdmin := boi.AdminDebugEnabled && actorInfo.IsAdmin
		accessAsSuperAdmin := boi.AdminDebugEnabled && actorInfo.IsSuperAdmin
		var actorType pacta.AuditLogActorType
		if accessAsOwner {
			actorType = pacta.AuditLogActorType_Owner
		} else if accessAsAdmin {
			actorType = pacta.AuditLogActorType_Admin
		} else if accessAsSuperAdmin {
			actorType = pacta.AuditLogActorType_SuperAdmin
		} else {
			// DENY CASE
			return nil, err404
		}
		auditLogs = append(auditLogs, &pacta.AuditLog{
			Action:             pacta.AuditLogAction_Download,
			ActorID:            string(actorInfo.UserID),
			ActorOwner:         &pacta.Owner{ID: actorInfo.OwnerID},
			ActorType:          actorType,
			PrimaryTargetType:  boi.PrimaryTargetType,
			PrimaryTargetID:    boi.PrimaryTargetID,
			PrimaryTargetOwner: &pacta.Owner{ID: boi.PrimaryTargetOwnerID},
		})
	}

	blobs, err := s.DB.Blobs(s.DB.NoTxn(ctx), blobIDs)
	if err != nil {
		if db.IsNotFound(err) {
			return nil, err404
		}
		return nil, oapierr.Internal("error getting blobs", zap.Error(err), zap.Strings("blob_ids", asStrs(blobIDs)))
	}

	if err = s.DB.CreateAuditLogs(s.DB.NoTxn(ctx), auditLogs); err != nil {
		return nil, oapierr.Internal("error creating audit logs - no download URLs generated", zap.Error(err), zap.Strings("blob_ids", asStrs(blobIDs)))
	}

	// Note, we're not parallelizing this because it is probably not nescessary.
	// The majority use case of this endpoint will be the user clicking a download
	// button, which will spin as it gets the URL, then turn into a dial as the
	// download starts. That allows us to only generate audit logs for true accesses,
	// and will typically happen on a single-file basis.
	response := api.AccessBlobContentResp{}
	for _, blob := range blobs {
		url, expiryTime, err := s.Blob.SignedDownloadURL(ctx, string(blob.BlobURI))
		if err != nil {
			return nil, oapierr.Internal("error getting signed download url", zap.Error(err), zap.String("blob_uri", string(blob.BlobURI)))
		}
		response.Items = append(response.Items, api.AccessBlobContentRespItem{
			BlobId:         string(blob.ID),
			DownloadUrl:    url,
			ExpirationTime: expiryTime,
		})
	}
	return api.AccessBlobContent200JSONResponse(response), nil
}

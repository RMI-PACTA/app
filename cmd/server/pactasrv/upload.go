package pactasrv

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/RMI/pacta/blob"
	"github.com/RMI/pacta/cmd/server/pactasrv/conv"
	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/oapierr"
	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
	"github.com/RMI/pacta/task"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Starts the process of uploading one or more portfolio files
// (POST /portfolio-upload)
func (s *Server) StartPortfolioUpload(ctx context.Context, request api.StartPortfolioUploadRequestObject) (api.StartPortfolioUploadResponseObject, error) {
	// TODO(#71) Implement basic limits
	actorInfo, err := s.getActorInfoOrFail(ctx)
	if err != nil {
		return nil, err
	}
	owner := &pacta.Owner{ID: actorInfo.OwnerID}
	holdingsDate, err := conv.HoldingsDateFromOAPI(&request.Body.HoldingsDate)
	if err != nil {
		return nil, err
	}
	n := len(request.Body.Items)
	if n > 25 {
		// TODO(#71) Implement basic limits
		return nil, oapierr.BadRequest("too many items")
	}
	blobs := make([]*pacta.Blob, n)
	respItems := make([]api.StartPortfolioUploadRespItem, n)
	for i, item := range request.Body.Items {
		fn := item.FileName
		extStr := filepath.Ext(fn)
		ft, err := pacta.ParseFileType(extStr)
		if err != nil {
			return nil, oapierr.BadRequest(
				fmt.Sprintf("invalid file type: %q", extStr),
				zap.String("file_type", extStr),
				zap.String("file_name", fn),
				zap.Error(err)).
				WithErrorID("INVALID_FILE_EXTENSION").
				WithMessage(fmt.Sprintf(`{"file": %q, "extension": %q}`, fn, extStr))
		}
		blobs[i] = &pacta.Blob{
			FileType: ft,
			FileName: item.FileName,
		}
		respItems[i].FileName = string(fn)
	}

	for i := range request.Body.Items {
		id := uuid.NewString()
		uri := blob.Join(s.Blob.Scheme(), s.PorfolioUploadURI, id)
		signed, _, err := s.Blob.SignedUploadURL(ctx, uri)
		if err != nil {
			return nil, oapierr.Internal("failed to sign blob URI", zap.String("uri", uri), zap.Error(err))
		}
		blobs[i].BlobURI = pacta.BlobURI(uri)
		respItems[i].UploadUrl = signed
	}

	err = s.DB.Transactional(ctx, func(tx db.Tx) error {
		for i := 0; i < n; i++ {
			blob := blobs[i]
			blobID, err := s.DB.CreateBlob(tx, blob)
			if err != nil {
				return fmt.Errorf("creating blob %d: %w", i, err)
			}
			blob.ID = blobID
			iuid, err := s.DB.CreateIncompleteUpload(tx, &pacta.IncompleteUpload{
				Blob:         blob,
				Name:         blob.FileName,
				HoldingsDate: holdingsDate,
				Owner:        owner,
			})
			if err != nil {
				return fmt.Errorf("creating incomplete upload %d: %w", i, err)
			}
			respItems[i].IncompleteUploadId = string(iuid)
			_, err = s.DB.CreateAuditLog(tx, &pacta.AuditLog{
				Action:             pacta.AuditLogAction_Create,
				ActorID:            string(actorInfo.UserID),
				ActorOwner:         owner,
				ActorType:          pacta.AuditLogActorType_Owner,
				PrimaryTargetType:  pacta.AuditLogTargetType_IncompleteUpload,
				PrimaryTargetID:    string(iuid),
				PrimaryTargetOwner: owner,
			})
			if err != nil {
				return fmt.Errorf("creating audit log %d: %w", i, err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, oapierr.Internal("failed to create uploads", zap.Error(err))
	}

	return api.StartPortfolioUpload200JSONResponse{Items: respItems}, nil
}

// Called after uploads of portfolios to cloud storage are complete.
// (POST /portfolio-upload:complete)
func (s *Server) CompletePortfolioUpload(ctx context.Context, request api.CompletePortfolioUploadRequestObject) (api.CompletePortfolioUploadResponseObject, error) {
	actorInfo, err := s.getActorInfoOrFail(ctx)
	if err != nil {
		return nil, err
	}
	ids := []pacta.IncompleteUploadID{}
	for _, item := range request.Body.Items {
		ids = append(ids, pacta.IncompleteUploadID(item.IncompleteUploadId))
	}
	if len(ids) == 0 {
		return nil, oapierr.BadRequest("no incomplete upload IDs provided")
	}
	// TODO(#71) Implement basic limits + validation
	blobURIs := []pacta.BlobURI{}
	err = s.DB.Transactional(ctx, func(tx db.Tx) error {
		ius, err := s.DB.IncompleteUploads(tx, ids)
		if err != nil {
			return oapierr.Internal("failed to query incomplete uploads", zap.Error(err))
		}
		blobIDs := []pacta.BlobID{}
		for _, id := range ids {
			iu := ius[id]
			if iu == nil || iu.Owner == nil || iu.Owner.ID != actorInfo.OwnerID {
				return oapierr.NotFound(
					fmt.Sprintf("incomplete upload %s does not belong to user", id),
					zap.String("incomplete_upload_id", string(id)),
					zap.String("owner_id", string(actorInfo.OwnerID)),
				)
			}
			blobIDs = append(blobIDs, iu.Blob.ID)
		}
		blobs, err := s.DB.Blobs(s.DB.NoTxn(ctx), blobIDs)
		if err != nil {
			return oapierr.Internal("failed to query blobs", zap.Error(err))
		}
		for _, blob := range blobs {
			blobURIs = append(blobURIs, blob.BlobURI)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	taskID, runnerID, err := s.TaskRunner.ParsePortfolio(ctx, &task.ParsePortfolioRequest{
		BlobURIs:            blobURIs,
		IncompleteUploadIDs: ids,
	})
	if err != nil {
		return nil, oapierr.Internal("failed to start task", zap.Error(err))
	}
	s.Logger.Info("triggered parse portfolio task",
		zap.String("task_id", string(taskID)),
		zap.String("task_runner_id", string(runnerID)))

	now := s.Now()
	err = s.DB.Transactional(ctx, func(tx db.Tx) error {
		for _, id := range ids {
			err := s.DB.UpdateIncompleteUpload(tx, id, db.SetIncompleteUploadRanAt(now))
			if err != nil {
				return oapierr.Internal("failed to update incomplete upload with ran_at", zap.Error(err))
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return api.CompletePortfolioUpload200JSONResponse{}, nil
}

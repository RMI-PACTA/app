package pactasrv

import (
	"context"
	"fmt"

	"github.com/RMI/pacta/cmd/server/pactasrv/conv"
	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/oapierr"
	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
	"go.uber.org/zap"
)

// (GET /incomplete-uploads)
func (s *Server) ListIncompleteUploads(ctx context.Context, request api.ListIncompleteUploadsRequestObject) (api.ListIncompleteUploadsResponseObject, error) {

	ownerID, err := s.getUserOwnerID(ctx)
	if err != nil {
		return nil, err
	}
	ius, err := s.DB.IncompleteUploadsByOwner(s.DB.NoTxn(ctx), ownerID)
	if err != nil {
		return nil, oapierr.Internal("failed to query incomplete uploads", zap.Error(err))
	}
	items, err := dereference(conv.IncompleteUploadsToOAPI(ius))
	if err != nil {
		return nil, err
	}
	return api.ListIncompleteUploads200JSONResponse{Items: items}, nil
}

// Deletes an incomplete upload by ID
// (DELETE /incomplete-upload/{id})
func (s *Server) DeleteIncompleteUpload(ctx context.Context, request api.DeleteIncompleteUploadRequestObject) (api.DeleteIncompleteUploadResponseObject, error) {
	id := pacta.IncompleteUploadID(request.Id)
	if err := s.incompleteUploadDoAuthzAndAuditLog(ctx, id, pacta.AuditLogAction_Delete); err != nil {
		return nil, err
	}
	blobURI, err := s.DB.DeleteIncompleteUpload(s.DB.NoTxn(ctx), id)
	if err != nil {
		return nil, oapierr.Internal("failed to delete incomplete upload", zap.Error(err))
	}
	if err := s.deleteBlobs(ctx, blobURI); err != nil {
		return nil, err
	}
	return api.DeleteIncompleteUpload204Response{}, nil
}

// Returns an incomplete upload by ID
// (GET /incomplete-upload/{id})
func (s *Server) FindIncompleteUploadById(ctx context.Context, request api.FindIncompleteUploadByIdRequestObject) (api.FindIncompleteUploadByIdResponseObject, error) {
	id := pacta.IncompleteUploadID(request.Id)
	if err := s.incompleteUploadDoAuthzAndAuditLog(ctx, id, pacta.AuditLogAction_ReadMetadata); err != nil {
		return nil, err
	}
	iu, err := s.DB.IncompleteUpload(s.DB.NoTxn(ctx), id)
	if err != nil {
		return nil, oapierr.Internal("failed to look up incomplete upload", zap.String("incomplete_upload_id", request.Id), zap.Error(err))
	}
	converted, err := conv.IncompleteUploadToOAPI(iu)
	if err != nil {
		return nil, err
	}
	return api.FindIncompleteUploadById200JSONResponse(*converted), nil
}

// Updates incomplete upload properties
// (PATCH /incomplete-upload/{id})
func (s *Server) UpdateIncompleteUpload(ctx context.Context, request api.UpdateIncompleteUploadRequestObject) (api.UpdateIncompleteUploadResponseObject, error) {
	if err := anyError(
		checkStringLimitSmallPtr("name", request.Body.Name),
		checkStringLimitMediumPtr("description", request.Body.Description),
	); err != nil {
		return nil, err
	}
	id := pacta.IncompleteUploadID(request.Id)
	mutations := []db.UpdateIncompleteUploadFn{}
	if request.Body.Name != nil {
		mutations = append(mutations, db.SetIncompleteUploadName(*request.Body.Name))
	}
	if request.Body.Description != nil {
		mutations = append(mutations, db.SetIncompleteUploadDescription(*request.Body.Description))
	}
	if len(mutations) > 0 {
		if err := s.incompleteUploadDoAuthzAndAuditLog(ctx, id, pacta.AuditLogAction_Update); err != nil {
			return nil, err
		}
	}
	if request.Body.AdminDebugEnabled != nil {
		if *request.Body.AdminDebugEnabled {
			if err := s.incompleteUploadDoAuthzAndAuditLog(ctx, id, pacta.AuditLogAction_EnableAdminDebug); err != nil {
				return nil, err
			}
		} else {
			if err := s.incompleteUploadDoAuthzAndAuditLog(ctx, id, pacta.AuditLogAction_DisableAdminDebug); err != nil {
				return nil, err
			}
		}
		mutations = append(mutations, db.SetIncompleteUploadAdminDebugEnabled(*request.Body.AdminDebugEnabled))
	}
	err := s.DB.UpdateIncompleteUpload(s.DB.NoTxn(ctx), id, mutations...)
	if err != nil {
		return nil, oapierr.Internal("failed to update incomplete upload", zap.Error(err))
	}
	return api.UpdateIncompleteUpload204Response{}, nil
}

func (s *Server) incompleteUploadDoAuthzAndAuditLog(ctx context.Context, iuID pacta.IncompleteUploadID, action pacta.AuditLogAction) error {
	actorInfo, err := s.getActorInfoOrErrIfAnon(ctx)
	if err != nil {
		return err
	}
	iu, err := s.DB.IncompleteUpload(s.DB.NoTxn(ctx), iuID)
	if err != nil {
		if db.IsNotFound(err) {
			return notFoundOrUnauthorized(actorInfo, action, pacta.AuditLogTargetType_IncompleteUpload, iuID)
		}
		return oapierr.Internal("failed to look up incomplete upload", zap.String("incomplete_upload_id", string(iuID)), zap.Error(err))
	}
	as := &authzStatus{
		primaryTargetID:      string(iuID),
		primaryTargetType:    pacta.AuditLogTargetType_IncompleteUpload,
		primaryTargetOwnerID: iu.Owner.ID,
		actorInfo:            actorInfo,
		action:               action,
	}
	switch action {
	case pacta.AuditLogAction_EnableAdminDebug, pacta.AuditLogAction_DisableAdminDebug:
		as.isAuthorized, as.authorizedAsActorType = allowIfOwner(actorInfo, iu.Owner.ID)
	case pacta.AuditLogAction_Update, pacta.AuditLogAction_Delete, pacta.AuditLogAction_ReadMetadata:
		as.isAuthorized, as.authorizedAsActorType = allowIfAdminOrOwner(actorInfo, iu.Owner.ID)
	default:
		return fmt.Errorf("unknown action %q for incomplete_upload authz", action)
	}
	return s.auditLogIfAuthorizedOrFail(ctx, as)
}

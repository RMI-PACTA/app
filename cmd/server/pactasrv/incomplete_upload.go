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
	// TODO(#12) Implement Authorization
	ownerID, err := s.getUserOwnerID(ctx)
	if err != nil {
		return nil, err
	}
	ius, err := s.DB.IncompleteUploadsByOwner(s.DB.NoTxn(ctx), ownerID)
	if err != nil {
		return nil, oapierr.Internal("failed to query incomplete uploads", zap.Error(err))
	}
	s.Logger.Info("queried incomplete uploads", zap.Int("count", len(ius)), zap.String("owner_id", string(ownerID)))
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
	_, err := s.checkIncompleteUploadAuthorization(ctx, id)
	if err != nil {
		return nil, err
	}
	blobURI, err := s.DB.DeleteIncompleteUpload(s.DB.NoTxn(ctx), id)
	if err != nil {
		return nil, oapierr.Internal("failed to delete incomplete upload", zap.Error(err))
	}
	if err := s.deleteBlobs(ctx, []string{string(blobURI)}); err != nil {
		return nil, oapierr.Internal("failed to delete blob", zap.Error(err))
	}
	return api.DeleteIncompleteUpload204Response{}, nil
}

// Returns an incomplete upload by ID
// (GET /incomplete-upload/{id})
func (s *Server) FindIncompleteUploadById(ctx context.Context, request api.FindIncompleteUploadByIdRequestObject) (api.FindIncompleteUploadByIdResponseObject, error) {
	iu, err := s.checkIncompleteUploadAuthorization(ctx, pacta.IncompleteUploadID(request.Id))
	if err != nil {
		return nil, err
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
	id := pacta.IncompleteUploadID(request.Id)
	_, err := s.checkIncompleteUploadAuthorization(ctx, id)
	if err != nil {
		return nil, err
	}
	// TODO(#12) Implement Authorization
	mutations := []db.UpdateIncompleteUploadFn{}
	if request.Body.Name != nil {
		mutations = append(mutations, db.SetIncompleteUploadName(*request.Body.Name))
	}
	if request.Body.Description != nil {
		mutations = append(mutations, db.SetIncompleteUploadDescription(*request.Body.Description))
	}
	if request.Body.AdminDebugEnabled != nil {
		mutations = append(mutations, db.SetIncompleteUploadAdminDebugEnabled(*request.Body.AdminDebugEnabled))
	}
	err = s.DB.UpdateIncompleteUpload(s.DB.NoTxn(ctx), id, mutations...)
	if err != nil {
		return nil, oapierr.Internal("failed to update incomplete upload", zap.Error(err))
	}
	return api.UpdateIncompleteUpload204Response{}, nil
}

func (s *Server) checkIncompleteUploadAuthorization(ctx context.Context, id pacta.IncompleteUploadID) (*pacta.IncompleteUpload, error) {
	// TODO(#12) Implement Authorization
	actorOwnerID, err := s.getUserOwnerID(ctx)
	if err != nil {
		return nil, err
	}
	// Extracted to a common variable so that we return the same response for not found and unauthorized.
	notFoundErr := func(err error) error {
		return oapierr.NotFound("incomplete upload not found", zap.String("incomplete_upload_id", string(id)), zap.Error(err))
	}
	iu, err := s.DB.IncompleteUpload(s.DB.NoTxn(ctx), id)
	if err != nil {
		if db.IsNotFound(err) {
			return nil, notFoundErr(err)
		}
		return nil, oapierr.Internal("failed to look up incomplete upload", zap.String("incomplete_upload_id", string(id)), zap.Error(err))
	}
	if iu.Owner.ID != actorOwnerID {
		return nil, notFoundErr(fmt.Errorf("incomplete upload does not belong to user: owner=%s actor=%s", iu.Owner.ID, actorOwnerID))
	}
	return iu, nil
}

package pactasrv

import (
	"context"

	"github.com/RMI/pacta/cmd/server/pactasrv/conv"
	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/oapierr"
	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
	"go.uber.org/zap"
)

// Returns a version of the PACTA model by ID
// (GET /pacta-version/{id})
func (s *Server) FindPactaVersionById(ctx context.Context, request api.FindPactaVersionByIdRequestObject) (api.FindPactaVersionByIdResponseObject, error) {
	// TODO(#12) Implement Authorization
	pv, err := s.DB.PACTAVersion(s.DB.NoTxn(ctx), pacta.PACTAVersionID(request.Id))
	if err != nil {
		if db.IsNotFound(err) {
			return nil, oapierr.NotFound("pacta version not found", zap.String("pacta_version_id", request.Id))
		}
		return nil, oapierr.Internal("failed to load pacta version", zap.String("pacta_version_id", request.Id), zap.Error(err))
	}
	resp, err := conv.PactaVersionToOAPI(pv)
	if err != nil {
		return nil, err
	}
	return api.FindPactaVersionById200JSONResponse(*resp), nil
}

// Returns all versions of the PACTA model
// (GET /pacta-versions)
func (s *Server) ListPactaVersions(ctx context.Context, request api.ListPactaVersionsRequestObject) (api.ListPactaVersionsResponseObject, error) {
	// TODO(#12) Implement Authorization
	pvs, err := s.DB.PACTAVersions(s.DB.NoTxn(ctx))
	if err != nil {
		return nil, oapierr.Internal("failed to list pacta versions", zap.Error(err))
	}
	resp, err := dereference(mapAll(pvs, conv.PactaVersionToOAPI))
	if err != nil {
		return nil, err
	}
	return api.ListPactaVersions200JSONResponse(resp), nil
}

// Creates a PACTA version
// (POST /pacta-versions)
func (s *Server) CreatePactaVersion(ctx context.Context, request api.CreatePactaVersionRequestObject) (api.CreatePactaVersionResponseObject, error) {
	// TODO(#12) Implement Authorization
	pv, err := conv.PactaVersionCreateFromOAPI(request.Body)
	if err != nil {
		return nil, err
	}
	if _, err := s.DB.CreatePACTAVersion(s.DB.NoTxn(ctx), pv); err != nil {
		return nil, oapierr.Internal("failed to create pacta version", zap.Error(err))
	}
	return api.CreatePactaVersion204Response{}, nil
}

// Updates a PACTA version
// (PATCH /pacta-version/{id})
func (s *Server) UpdatePactaVersion(ctx context.Context, request api.UpdatePactaVersionRequestObject) (api.UpdatePactaVersionResponseObject, error) {
	// TODO(#12) Implement Authorization
	id := pacta.PACTAVersionID(request.Id)
	mutations := []db.UpdatePACTAVersionFn{}
	b := request.Body
	if b.Description != nil {
		mutations = append(mutations, db.SetPACTAVersionDescription(*b.Description))
	}
	if b.Digest != nil {
		mutations = append(mutations, db.SetPACTAVersionDigest(*b.Digest))
	}
	if b.Name != nil {
		mutations = append(mutations, db.SetPACTAVersionName(*b.Name))
	}
	err := s.DB.UpdatePACTAVersion(s.DB.NoTxn(ctx), id, mutations...)
	if err != nil {
		return nil, oapierr.Internal("failed to update pacta version", zap.String("pacta_version_id", string(id)), zap.Error(err))
	}
	return api.UpdatePactaVersion204Response{}, nil

}

// Deletes a pacta version by ID
// (DELETE /pacta-version/{id})
func (s *Server) DeletePactaVersion(ctx context.Context, request api.DeletePactaVersionRequestObject) (api.DeletePactaVersionResponseObject, error) {
	// TODO(#12) Implement Authorization
	id := pacta.PACTAVersionID(request.Id)
	err := s.DB.DeletePACTAVersion(s.DB.NoTxn(ctx), id)
	if err != nil {
		return nil, oapierr.Internal("failed to delete pacta version", zap.String("pacta_version_id", string(id)), zap.Error(err))
	}
	return api.DeletePactaVersion204Response{}, nil
}

// Marks this version of the PACTA model as the default
// (POST /pacta-version/{id}/set-default)
func (s *Server) MarkPactaVersionAsDefault(ctx context.Context, request api.MarkPactaVersionAsDefaultRequestObject) (api.MarkPactaVersionAsDefaultResponseObject, error) {
	// TODO(#12) Implement Authorization
	id := pacta.PACTAVersionID(request.Id)
	err := s.DB.SetDefaultPACTAVersion(s.DB.NoTxn(ctx), id)
	if err != nil {
		return nil, oapierr.Internal("failed to set default pacta version", zap.String("pacta_version_id", string(id)), zap.Error(err))
	}
	return api.MarkPactaVersionAsDefault204Response{}, nil
}

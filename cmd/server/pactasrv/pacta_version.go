package pactasrv

import (
	"context"
	"fmt"

	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
	"github.com/Silicon-Ally/gqlerr"
	"go.uber.org/zap"
)

// Returns a version of the PACTA model by ID
// (GET /pacta-version/{id})
func (s *Server) FindPactaVersionById(ctx context.Context, request api.FindPactaVersionByIdRequestObject) (api.FindPactaVersionByIdResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// Returns all versions of the PACTA model
// (GET /pacta-versions)
func (s *Server) ListPactaVersions(ctx context.Context, request api.ListPactaVersionsRequestObject) (api.ListPactaVersionsResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// Creates a PACTA version
// (POST /pacta-versions)
func (s *Server) CreatePactaVersion(ctx context.Context, request api.CreatePactaVersionRequestObject) (api.CreatePactaVersionResponseObject, error) {
	// TODO(grady) Authz
	_, err := s.DB.CreatePACTAVersion(s.DB.NoTxn(ctx), &pacta.PACTAVersion{
		Name:        request.Body.Name,
		Description: request.Body.Description,
		Digest:      request.Body.Digest,
	})
	if err != nil {
		return nil, gqlerr.Internal(ctx, "failed to create PACTA version", zap.Error(err))
	}
	return nil, nil
}

// Updates a PACTA version
// (PATCH /pacta-version/{id})
func (s *Server) UpdatePactaVersion(ctx context.Context, request api.UpdatePactaVersionRequestObject) (api.UpdatePactaVersionResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// Deletes a pacta version by ID
// (DELETE /pacta-version/{id})
func (s *Server) DeletePactaVersion(ctx context.Context, request api.DeletePactaVersionRequestObject) (api.DeletePactaVersionResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

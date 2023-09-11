package pactasrv

import (
	"context"

	"github.com/RMI/pacta/db"
	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
)

// Returns a version of the PACTA model by ID
// (GET /pacta-version/{id})
func (s *Server) FindPactaVersionById(ctx context.Context, request api.FindPactaVersionByIdRequestObject) (api.FindPactaVersionByIdResponseObject, error) {
	pv, err := s.findPactaVersionById(ctx, request)
	if err != nil {
		e := errToAPIError(err)
		return api.FindPactaVersionByIddefaultJSONResponse{
			Body:       e,
			StatusCode: int(e.Code),
		}, nil
	}
	return api.FindPactaVersionById200JSONResponse(*pv), nil
}

func (s *Server) findPactaVersionById(ctx context.Context, request api.FindPactaVersionByIdRequestObject) (*api.PactaVersion, error) {
	// TODO(#12) Implement Authorization
	pv, err := s.DB.PACTAVersion(s.DB.NoTxn(ctx), pacta.PACTAVersionID(request.Id))
	if err != nil {
		return nil, errorInternal(err)
	}
	return pactaVersionToOAPI(pv)
}

// Returns all versions of the PACTA model
// (GET /pacta-versions)
func (s *Server) ListPactaVersions(ctx context.Context, request api.ListPactaVersionsRequestObject) (api.ListPactaVersionsResponseObject, error) {
	pvs, err := s.listPactaVersions(ctx, request)
	if err != nil {
		e := errToAPIError(err)
		return api.ListPactaVersionsdefaultJSONResponse{
			Body:       e,
			StatusCode: int(e.Code),
		}, nil
	}
	return api.ListPactaVersions200JSONResponse(pvs), nil
}

func (s *Server) listPactaVersions(ctx context.Context, _request api.ListPactaVersionsRequestObject) ([]api.PactaVersion, error) {
	// TODO(#12) Implement Authorization
	pvs, err := s.DB.PACTAVersions(s.DB.NoTxn(ctx))
	if err != nil {
		return nil, errorInternal(err)
	}
	return dereference(mapAll(pvs, pactaVersionToOAPI))
}

// Creates a PACTA version
// (POST /pacta-versions)
func (s *Server) CreatePactaVersion(ctx context.Context, request api.CreatePactaVersionRequestObject) (api.CreatePactaVersionResponseObject, error) {
	err := s.createPactaVersion(ctx, request)
	if err != nil {
		e := errToAPIError(err)
		return api.CreatePactaVersiondefaultJSONResponse{
			Body:       e,
			StatusCode: int(e.Code),
		}, nil
	}
	return api.CreatePactaVersion200JSONResponse(emptySuccess), nil
}

func (s *Server) createPactaVersion(ctx context.Context, request api.CreatePactaVersionRequestObject) error {
	// TODO(#12) Implement Authorization
	pv, err := pactaVersionCreateToPACTA(request.Body)
	if err != nil {
		return errorBadRequest("body", err.Error())
	}
	if _, err := s.DB.CreatePACTAVersion(s.DB.NoTxn(ctx), pv); err != nil {
		return errorInternal(err)
	}
	return nil
}

// Updates a PACTA version
// (PATCH /pacta-version/{id})
func (s *Server) UpdatePactaVersion(ctx context.Context, request api.UpdatePactaVersionRequestObject) (api.UpdatePactaVersionResponseObject, error) {
	err := s.updatePactaVersion(ctx, request)
	if err != nil {
		e := errToAPIError(err)
		return api.UpdatePactaVersiondefaultJSONResponse{
			Body:       e,
			StatusCode: int(e.Code),
		}, nil
	}
	return api.UpdatePactaVersion200JSONResponse(emptySuccess), nil
}

func (s *Server) updatePactaVersion(ctx context.Context, request api.UpdatePactaVersionRequestObject) error {
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
		return errorInternal(err)
	}
	return nil

}

// Deletes a pacta version by ID
// (DELETE /pacta-version/{id})
func (s *Server) DeletePactaVersion(ctx context.Context, request api.DeletePactaVersionRequestObject) (api.DeletePactaVersionResponseObject, error) {
	err := s.deletePactaVersion(ctx, request)
	if err != nil {
		e := errToAPIError(err)
		return api.DeletePactaVersiondefaultJSONResponse{
			Body:       e,
			StatusCode: int(e.Code),
		}, nil
	}
	return api.DeletePactaVersion200JSONResponse(emptySuccess), nil
}

func (s *Server) deletePactaVersion(ctx context.Context, request api.DeletePactaVersionRequestObject) error {
	// TODO(#12) Implement Authorization
	id := pacta.PACTAVersionID(request.Id)
	err := s.DB.DeletePACTAVersion(s.DB.NoTxn(ctx), id)
	if err != nil {
		return errorInternal(err)
	}
	return nil
}

// Marks this version of the PACTA model as the default
// (POST /pacta-version/{id}/set-default)
func (s *Server) MarkPactaVersionAsDefault(ctx context.Context, request api.MarkPactaVersionAsDefaultRequestObject) (api.MarkPactaVersionAsDefaultResponseObject, error) {
	err := s.markPactaVersionAsDefault(ctx, request)
	if err != nil {
		e := errToAPIError(err)
		return api.MarkPactaVersionAsDefaultdefaultJSONResponse{
			Body:       e,
			StatusCode: int(e.Code),
		}, nil
	}
	return api.MarkPactaVersionAsDefault200JSONResponse(emptySuccess), nil
}

func (s *Server) markPactaVersionAsDefault(ctx context.Context, request api.MarkPactaVersionAsDefaultRequestObject) error {
	// TODO(#12) Implement Authorization
	id := pacta.PACTAVersionID(request.Id)
	err := s.DB.SetDefaultPACTAVersion(s.DB.NoTxn(ctx), id)
	if err != nil {
		return errorInternal(err)
	}
	return nil
}

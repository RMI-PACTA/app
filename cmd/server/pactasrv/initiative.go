package pactasrv

import (
	"context"

	"github.com/RMI/pacta/db"
	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
)

// Creates a initiative
// (POST /initiatives)
func (s *Server) CreateInitiative(ctx context.Context, request api.CreateInitiativeRequestObject) (api.CreateInitiativeResponseObject, error) {
	err := s.createInitiative(ctx, request)
	if err != nil {
		return errToAPIError(err)
	}
	return api.CreateInitiative200JSONResponse{}, nil
}

func (s *Server) createInitiative(ctx context.Context, request api.CreateInitiativeRequestObject) error {
	// TODO(#12) Implement Authorization
	i, err := initiativeCreateFromOAPI(request.Body)
	if err != nil {
		return err
	}
	err = s.DB.CreateInitiative(s.DB.NoTxn(ctx), i)
	if err != nil {
		return errorInternal(err)
	}
	return nil
}

// Updates an initiative
// (PATCH /initiative/{id})
func (s *Server) UpdateInitiative(ctx context.Context, request api.UpdateInitiativeRequestObject) (api.UpdateInitiativeResponseObject, error) {
	err := s.updateInitiative(ctx, request)
	if err != nil {
		return errToAPIError(err)
	}
	return api.UpdateInitiative200JSONResponse{}, nil
}

func (s *Server) updateInitiative(ctx context.Context, request api.UpdateInitiativeRequestObject) error {
	// TODO(#12) Implement Authorization
	id := pacta.InitiativeID(request.Id)
	mutations := []db.UpdateInitiativeFn{}
	b := request.Params.Body
	if b.Affiliation != nil {
		mutations = append(mutations, db.SetInitiativeAffiliation(*b.Affiliation))
	}
	if b.InternalDescription != nil {
		mutations = append(mutations, db.SetInitiativeInternalDescription(*b.InternalDescription))
	}
	if b.IsAcceptingNewMembers != nil {
		mutations = append(mutations, db.SetInitiativeIsAcceptingNewMembers(*b.IsAcceptingNewMembers))
	}
	if b.IsAcceptingNewPortfolios != nil {
		mutations = append(mutations, db.SetInitiativeIsAcceptingNewPortfolios(*b.IsAcceptingNewPortfolios))
	}
	if b.Language != nil {
		lang, err := pacta.ParseLanguage(string(*b.Language))
		if err != nil {
			return errorBadRequest("language", err.Error())
		}
		mutations = append(mutations, db.SetInitiativeLanguage(lang))
	}
	if b.Name != nil {
		mutations = append(mutations, db.SetInitiativeName(*b.Name))
	}
	if b.PactaVersion != nil {
		mutations = append(mutations, db.SetInitiativePACTAVersion(pacta.PACTAVersionID(*b.PactaVersion)))
	}
	if b.PublicDescription != nil {
		mutations = append(mutations, db.SetInitiativePublicDescription(*b.PublicDescription))
	}
	if b.RequiresInvitationToJoin != nil {
		mutations = append(mutations, db.SetInitiativeRequiresInvitationToJoin(*b.RequiresInvitationToJoin))
	}
	err := s.DB.UpdateInitiative(s.DB.NoTxn(ctx), id, mutations...)
	if err != nil {
		return errorInternal(err)
	}
	return nil
}

// Deletes an initiative by id
// (DELETE /initiative/{id})
func (s *Server) DeleteInitiative(ctx context.Context, request api.DeleteInitiativeRequestObject) (api.DeleteInitiativeResponseObject, error) {
	err := s.deleteInitiative(ctx, request)
	if err != nil {
		return errToAPIError(err)
	}
	return api.DeleteInitiative200JSONResponse{}, nil
}

func (s *Server) deleteInitiative(ctx context.Context, request api.DeleteInitiativeRequestObject) error {
	// TODO(#12) Implement Authorization
	err := s.DB.DeleteInitiative(s.DB.NoTxn(ctx), pacta.InitiativeID(request.Id))
	if err != nil {
		return errorInternal(err)
	}
	return nil
}

// Returns an initiative by ID
// (GET /initiative/{id})
func (s *Server) FindInitiativeById(ctx context.Context, request api.FindInitiativeByIdRequestObject) (api.FindInitiativeByIdResponseObject, error) {
	result, err := s.findInitiativeById(ctx, request)
	if err != nil {
		return errToAPIError(err)
	}
	return api.FindInitiativeById200JSONResponse(*result), nil
}

func (s *Server) findInitiativeById(ctx context.Context, request api.FindInitiativeByIdRequestObject) (*api.Initiative, error) {
	// TODO(#12) Implement Authorization
	i, err := s.DB.Initiative(s.DB.NoTxn(ctx), pacta.InitiativeID(request.Id))
	if err != nil {
		if db.IsNotFound(err) {
			return nil, errorNotFound("initiative", request.Id)
		}
		return nil, errorInternal(err)
	}
	return initiativeToOAPI(i)
}

// Returns all initiatives
// (GET /initiatives)
func (s *Server) ListInitiatives(ctx context.Context, request api.ListInitiativesRequestObject) (api.ListInitiativesResponseObject, error) {
	result, err := s.listInitiatives(ctx, request)
	if err != nil {
		return errToAPIError(err)
	}
	return api.ListInitiatives200JSONResponse(result), nil
}

func (s *Server) listInitiatives(ctx context.Context, request api.ListInitiativesRequestObject) ([]api.Initiative, error) {
	is, err := s.DB.AllInitiatives(s.DB.NoTxn(ctx))
	if err != nil {
		return nil, errorInternal(err)
	}
	return dereference(mapAll(is, initiativeToOAPI))
}

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

// Creates a initiative
// (POST /initiatives)
func (s *Server) CreateInitiative(ctx context.Context, request api.CreateInitiativeRequestObject) (api.CreateInitiativeResponseObject, error) {
	if err := anyError(
		checkStringLimitSmall("id", request.Body.Id),
		checkStringLimitSmallPtr("affiliation", request.Body.Affiliation),
		checkStringLimitSmall("name", request.Body.Name),
		checkStringLimitMediumPtr("internal_description", request.Body.InternalDescription),
		checkStringLimitMediumPtr("public_description", request.Body.PublicDescription),
	); err != nil {
		return nil, err
	}
	// TODO(#12) Implement Authorization
	i, err := conv.InitiativeCreateFromOAPI(request.Body)
	if err != nil {
		return nil, err
	}
	err = s.DB.CreateInitiative(s.DB.NoTxn(ctx), i)
	if err != nil {
		return nil, oapierr.Internal("failed to create initiative", zap.Error(err))
	}
	return api.CreateInitiative204Response{}, nil
}

// Updates an initiative
// (PATCH /initiative/{id})
func (s *Server) UpdateInitiative(ctx context.Context, request api.UpdateInitiativeRequestObject) (api.UpdateInitiativeResponseObject, error) {
	if err := anyError(
		checkStringLimitSmallPtr("affiliation", request.Body.Affiliation),
		checkStringLimitSmallPtr("name", request.Body.Name),
		checkStringLimitMediumPtr("internal_description", request.Body.InternalDescription),
		checkStringLimitMediumPtr("public_description", request.Body.PublicDescription),
	); err != nil {
		return nil, err
	}
	// TODO(#12) Implement Authorization
	id := pacta.InitiativeID(request.Id)
	mutations := []db.UpdateInitiativeFn{}
	b := request.Body
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
			return nil, oapierr.BadRequest("failed to parse language", zap.Error(err))
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
		return nil, oapierr.Internal("failed to update initiative", zap.String("initiative_id", string(id)), zap.Error(err))
	}
	return api.UpdateInitiative204Response{}, nil
}

// Deletes an initiative by id
// (DELETE /initiative/{id})
func (s *Server) DeleteInitiative(ctx context.Context, request api.DeleteInitiativeRequestObject) (api.DeleteInitiativeResponseObject, error) {
	// TODO(#12) Implement Authorization
	err := s.DB.DeleteInitiative(s.DB.NoTxn(ctx), pacta.InitiativeID(request.Id))
	if err != nil {
		return nil, oapierr.Internal("failed to delete initiative", zap.Error(err))
	}
	return api.DeleteInitiative204Response{}, nil
}

// Returns an initiative by ID
// (GET /initiative/{id})
func (s *Server) FindInitiativeById(ctx context.Context, request api.FindInitiativeByIdRequestObject) (api.FindInitiativeByIdResponseObject, error) {
	// TODO(#12) Implement Authorization
	i, err := s.DB.Initiative(s.DB.NoTxn(ctx), pacta.InitiativeID(request.Id))
	if err != nil {
		if db.IsNotFound(err) {
			return nil, oapierr.NotFound("initiative not found", zap.String("initiative_id", request.Id))
		}
		return nil, oapierr.Internal("failed to load initiative", zap.String("initiative_id", request.Id), zap.Error(err))
	}
	resp, err := conv.InitiativeToOAPI(i)
	if err != nil {
		return nil, err
	}
	return api.FindInitiativeById200JSONResponse(*resp), nil
}

// Returns all initiatives
// (GET /initiatives)
func (s *Server) ListInitiatives(ctx context.Context, request api.ListInitiativesRequestObject) (api.ListInitiativesResponseObject, error) {
	is, err := s.DB.AllInitiatives(s.DB.NoTxn(ctx))
	if err != nil {
		return nil, oapierr.Internal("failed to load all initiatives", zap.Error(err))
	}
	result, err := dereference(mapAll(is, conv.InitiativeToOAPI))
	if err != nil {
		return nil, err
	}
	return api.ListInitiatives200JSONResponse(result), nil
}

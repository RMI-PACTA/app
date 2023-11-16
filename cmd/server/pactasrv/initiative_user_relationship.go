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

// Returns all initiative user relationships for the user that the user has access to view
// (GET /initiative/{id}/user-relationships)
func (s *Server) ListInitiativeUserRelationshipsByUser(ctx context.Context, request api.ListInitiativeUserRelationshipsByUserRequestObject) (api.ListInitiativeUserRelationshipsByUserResponseObject, error) {
	// TODO(#12) Implement Authorization
	iurs, err := s.DB.InitiativeUserRelationshipsByUser(s.DB.NoTxn(ctx), pacta.UserID(request.UserId))
	if err != nil {
		return nil, oapierr.Internal("failed to retrieve initiative user relationships by user", zap.Error(err))
	}
	result, err := dereference(mapAll(iurs, conv.InitiativeUserRelationshipToOAPI))
	if err != nil {
		return nil, err
	}
	return api.ListInitiativeUserRelationshipsByUser200JSONResponse(result), nil
}

// Returns all initiative user relationships for the initiative that the user has access to view
// (GET /initiative/user-relationships/{id})
func (s *Server) ListInitiativeUserRelationshipsByInitiative(ctx context.Context, request api.ListInitiativeUserRelationshipsByInitiativeRequestObject) (api.ListInitiativeUserRelationshipsByInitiativeResponseObject, error) {
	// TODO(#12) Implement Authorization
	iurs, err := s.DB.InitiativeUserRelationshipsByInitiatives(s.DB.NoTxn(ctx), pacta.InitiativeID(request.InitiativeId))
	if err != nil {
		return nil, oapierr.Internal("failed to retrieve initiative user relationships by user", zap.Error(err))
	}
	result, err := dereference(mapAll(iurs, conv.InitiativeUserRelationshipToOAPI))
	if err != nil {
		return nil, err
	}
	return api.ListInitiativeUserRelationshipsByInitiative200JSONResponse(result), nil
}

// Returns the initiative user relationship from this id, if it exists
// (GET /initiative/{initiativeId}/user-relationship/{userId})
func (s *Server) GetInitiativeUserRelationship(ctx context.Context, request api.GetInitiativeUserRelationshipRequestObject) (api.GetInitiativeUserRelationshipResponseObject, error) {
	// TODO(#12) Implement Authorization
	iur, err := s.DB.InitiativeUserRelationship(s.DB.NoTxn(ctx), pacta.InitiativeID(request.InitiativeId), pacta.UserID(request.UserId))
	if err != nil {
		return nil, oapierr.Internal("failed to retrieve initiative user relationship", zap.Error(err))
	}
	result, err := conv.InitiativeUserRelationshipToOAPI(iur)
	if err != nil {
		return nil, err
	}
	return api.GetInitiativeUserRelationship200JSONResponse(*result), nil
}

// Updates initiative user relationship properties
// (PATCH /initiative/{initiativeId}/user-relationship/{userId})
func (s *Server) UpdateInitiativeUserRelationship(ctx context.Context, request api.UpdateInitiativeUserRelationshipRequestObject) (api.UpdateInitiativeUserRelationshipResponseObject, error) {
	// TODO(#12) Implement Authorization
	mutations := []db.UpdateInitiativeUserRelationshipFn{}
	if request.Body.Manager != nil {
		mutations = append(mutations, db.SetInitiativeUserRelationshipManager(*request.Body.Manager))
	}
	if request.Body.Member != nil {
		mutations = append(mutations, db.SetInitiativeUserRelationshipMember(*request.Body.Member))
	}
	err := s.DB.UpdateInitiativeUserRelationship(s.DB.NoTxn(ctx), pacta.InitiativeID(request.InitiativeId), pacta.UserID(request.UserId), mutations...)
	if err != nil {
		return nil, oapierr.Internal("failed to update initiative user relationship", zap.Error(err))
	}
	return api.UpdateInitiativeUserRelationship204Response{}, nil
}

package pactasrv

import (
	"context"
	"fmt"

	"github.com/RMI/pacta/cmd/server/pactasrv/conv"
	"github.com/RMI/pacta/oapierr"
	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
	"go.uber.org/zap"
)

// Creates an initiative invitation
// (POST /initiative-invitation)
func (s *Server) CreateInitiativeInvitation(ctx context.Context, request api.CreateInitiativeInvitationRequestObject) (api.CreateInitiativeInvitationResponseObject, error) {
	// TODO(#12) Implement Authorization
	ii, err := conv.InitiativeInvitationFromOAPI(request.Body)
	if err != nil {
		return nil, err
	}
	id, err = s.DB.CreateInitiativeInvitation(s.DB.NoTxn(ctx), ii)
	if err != nil {
		return nil, oapierr.Internal("failed to create initiative invitation", zap.Error(err))
	}
	if id != ii.ID {
		return nil, oapierr.Internal(
			"failed to create initiative invitation: ID mismatch",
			zap.String("requested_id", string(ii.ID)),
			zap.String("actual_id", string(id)),
		)
	}
	return api.CreateInitiativeInvitation204Response{}, nil
}

// Deletes an initiative invitation by id
// (DELETE /initiative-invitation/{id})
func (s *Server) DeleteInitiativeInvitation(ctx context.Context, request api.DeleteInitiativeInvitationRequestObject) (api.DeleteInitiativeInvitationResponseObject, error) {
	// TODO(#12) Implement Authorization
	err := s.DB.DeleteInitiativeInvitation(s.DB.NoTxn(ctx), pacta.InitiativeInvitationID(request.Id))
	if err != nil {
		return nil, oapierr.Internal("failed to delete initiative invitation", zap.Error(err))
	}
	return api.DeleteInitiativeInvitation204Response{}, nil
}

// Returns the initiative invitation from this id, if it exists
// (GET /initiative-invitation/{id})
func (s *Server) GetInitiativeInvitation(ctx context.Context, request api.GetInitiativeInvitationRequestObject) (api.GetInitiativeInvitationResponseObject, error) {
	// TODO(#12) Implement Authorization
	ii, err := s.DB.InitiativeInvitation(s.DB.NoTxn(ctx), pacta.InitiativeInvitationID(request.Id))
	if err != nil {
		return nil, oapierr.Internal("failed to retrieve initiative invitation", zap.Error(err))
	}
	result, err := conv.InitiativeInvitationToOAPI(ii)
	if err != nil {
		return nil, err
	}
	return api.GetInitiativeInvitation200JSONResponse(*result), nil
}

// Claims this initiative invitation, if it exists
// (POST /initiative-invitation/{id})
func (s *Server) ClaimInitiativeInvitation(ctx context.Context, request api.ClaimInitiativeInvitationRequestObject) (api.ClaimInitiativeInvitationResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}

// Returns all initiative invitations associated with the initiative
// (GET /initiative/{id}/invitations)
func (s *Server) ListInitiativeInvitations(ctx context.Context, request api.ListInitiativeInvitationsRequestObject) (api.ListInitiativeInvitationsResponseObject, error) {
	// TODO(#12) Implement Authorization
	iis, err := s.DB.InitiativeInvitationsByInitiative(s.DB.NoTxn(ctx), pacta.InitiativeID(request.Id))
	if err != nil {
		return nil, oapierr.Internal("failed to list initiative invitations", zap.Error(err))
	}
	result, err := dereference(mapAll(iis, conv.InitiativeInvitationToOAPI))
	if err != nil {
		return nil, err
	}
	return api.ListInitiativeInvitations200JSONResponse(result), nil
}

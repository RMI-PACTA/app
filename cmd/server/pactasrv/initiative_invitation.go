package pactasrv

import (
	"context"
	"fmt"
	"time"

	"github.com/RMI/pacta/cmd/server/pactasrv/conv"
	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/oapierr"
	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
	"go.uber.org/zap"
)

// Creates an initiative invitation
// (POST /initiative-invitation)
func (s *Server) CreateInitiativeInvitation(ctx context.Context, request api.CreateInitiativeInvitationRequestObject) (api.CreateInitiativeInvitationResponseObject, error) {
	if err := checkStringLimitSmall("id", request.Body.Id); err != nil {
		return nil, err
	}
	// TODO(#12) Implement Authorization
	ii, err := conv.InitiativeInvitationFromOAPI(request.Body)
	if err != nil {
		return nil, err
	}
	id, err := s.DB.CreateInitiativeInvitation(s.DB.NoTxn(ctx), ii)
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
// (POST /initiative-invitation/{id}:claim)
func (s *Server) ClaimInitiativeInvitation(ctx context.Context, request api.ClaimInitiativeInvitationRequestObject) (api.ClaimInitiativeInvitationResponseObject, error) {
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}
	var customErr api.ClaimInitiativeInvitationResponseObject
	err = s.DB.Transactional(ctx, func(tx db.Tx) error {
		ii, err := s.DB.InitiativeInvitation(tx, pacta.InitiativeInvitationID(request.Id))
		if err != nil {
			return fmt.Errorf("looking up initiative invite: %w", err)
		}
		if !ii.UsedAt.IsZero() || ii.UsedBy != nil {
			if ii.UsedBy != nil && ii.UsedBy.ID == userID {
				// We don't return an error if the same user tries to claim the same invitation twice,
				// which might happen by accident, but wouldn't impact the state of the initiative memberships.
				// We may want to log this, though.
				return nil
			} else {
				customErr = api.ClaimInitiativeInvitation409Response{}
				return fmt.Errorf("initiative is already used: %+v", ii)
			}
		}
		err = s.DB.UpdateInitiativeInvitation(tx, ii.ID,
			db.SetInitiativeInvitationUsedAt(time.Now()),
			db.SetInitiativeInvitationUsedBy(userID))
		if err != nil {
			return fmt.Errorf("updating initiative invite: %w", err)
		}
		err = s.DB.UpdateInitiativeUserRelationship(tx, ii.Initiative.ID, userID,
			db.SetInitiativeUserRelationshipMember(true))
		if err != nil {
			return fmt.Errorf("creating initiative membership: %w", err)
		}
		return nil
	})
	if err != nil {
		if customErr != nil {
			return customErr, nil
		}
		return nil, oapierr.Internal("failed to claim initiative invitation", zap.Error(err))
	}
	return api.ClaimInitiativeInvitation204Response{}, nil
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

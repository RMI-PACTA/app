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
	ii, err := conv.InitiativeInvitationFromOAPI(request.Body)
	if err != nil {
		return nil, err
	}
	if err := s.initiativeInvitationDoAuthzAndAuditLog(ctx, ii.Initiative.ID, ii.ID, pacta.AuditLogAction_Create); err != nil {
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
	iiID := pacta.InitiativeInvitationID(request.Id)
	ii, err := s.DB.InitiativeInvitation(s.DB.NoTxn(ctx), iiID)
	if err != nil {
		return nil, oapierr.Internal("failed to retrieve initiative invitation", zap.Error(err))
	}
	if err := s.initiativeInvitationDoAuthzAndAuditLog(ctx, ii.Initiative.ID, iiID, pacta.AuditLogAction_Delete); err != nil {
		return nil, err
	}
	err = s.DB.DeleteInitiativeInvitation(s.DB.NoTxn(ctx), pacta.InitiativeInvitationID(request.Id))
	if err != nil {
		return nil, oapierr.Internal("failed to delete initiative invitation", zap.Error(err))
	}
	return api.DeleteInitiativeInvitation204Response{}, nil
}

// Returns the initiative invitation from this id, if it exists
// (GET /initiative-invitation/{id})
func (s *Server) GetInitiativeInvitation(ctx context.Context, request api.GetInitiativeInvitationRequestObject) (api.GetInitiativeInvitationResponseObject, error) {
	iiID := pacta.InitiativeInvitationID(request.Id)
	ii, err := s.DB.InitiativeInvitation(s.DB.NoTxn(ctx), iiID)
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
	id := pacta.InitiativeID(request.Id)
	if err := s.initiativeInvitationDoAuthzAndAuditLog(ctx, id, "", pacta.AuditLogAction_ReadMetadata); err != nil {
		return nil, err
	}
	iis, err := s.DB.InitiativeInvitationsByInitiative(s.DB.NoTxn(ctx), id)
	if err != nil {
		return nil, oapierr.Internal("failed to list initiative invitations", zap.Error(err))
	}
	result, err := dereference(mapAll(iis, conv.InitiativeInvitationToOAPI))
	if err != nil {
		return nil, err
	}
	return api.ListInitiativeInvitations200JSONResponse(result), nil
}

func (s *Server) userIsInitiativeManagerOrAdmin(ctx context.Context, iID pacta.InitiativeID) (bool, error) {
	actorInfo, err := s.getActorInfoOrErrIfAnon(ctx)
	if err != nil {
		return false, err
	}
	if actorInfo.IsAdmin || actorInfo.IsSuperAdmin {
		return true, nil
	}
	iurs, err := s.DB.InitiativeUserRelationshipsByInitiative(s.DB.NoTxn(ctx), iID)
	if err != nil {
		return false, oapierr.Internal("failed to list initiative user relationships", zap.Error(err))
	}
	for _, iur := range iurs {
		if iur.User.ID == actorInfo.UserID && iur.Manager {
			return true, nil
		}
	}
	return false, nil
}

func (s *Server) initiativeInvitationDoAuthzAndAuditLog(ctx context.Context, iID pacta.InitiativeID, iiID pacta.InitiativeInvitationID, action pacta.AuditLogAction) error {
	actorInfo, err := s.getActorInfoOrErrIfAnon(ctx)
	if err != nil {
		return err
	}
	iurs, err := s.DB.InitiativeUserRelationshipsByInitiative(s.DB.NoTxn(ctx), iID)
	if err != nil {
		return oapierr.Internal("failed to list initiative user relationships", zap.Error(err))
	}
	actorIsInitiativeManager := false
	for _, iur := range iurs {
		if iur.User.ID == actorInfo.UserID {
			actorIsInitiativeManager = iur.Manager
			break
		}
	}
	as := &authzStatus{
		primaryTargetID:      string(iID),
		primaryTargetType:    pacta.AuditLogTargetType_Initiative,
		primaryTargetOwnerID: systemOwnedEntityOwner,
		actorInfo:            actorInfo,
		action:               action,
	}
	if iiID != "" {
		as.secondaryTargetID = string(iiID)
		as.secondaryTargetType = pacta.AuditLogTargetType_InitiativeInvitation
		as.secondaryTargetOwnerID = pacta.OwnerID(systemOwnedEntityOwner)
	}
	switch action {
	case pacta.AuditLogAction_Delete, pacta.AuditLogAction_ReadMetadata, pacta.AuditLogAction_Create:
		if actorIsInitiativeManager {
			as.authorizedAsActorType = ptr(pacta.AuditLogActorType_Owner)
			as.isAuthorized = true
		} else {
			as.isAuthorized, as.authorizedAsActorType = allowIfAdmin(actorInfo)
		}
	default:
		return fmt.Errorf("unknown action %q for initiative_invitation authz", action)
	}
	return s.auditLogIfAuthorizedOrFail(ctx, as)
}

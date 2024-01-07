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

// Returns all initiative user relationships for the user that the user has access to view
// (GET /initiative/{id}/user-relationships)
func (s *Server) ListInitiativeUserRelationshipsByUser(ctx context.Context, request api.ListInitiativeUserRelationshipsByUserRequestObject) (api.ListInitiativeUserRelationshipsByUserResponseObject, error) {
	id := pacta.UserID(request.UserId)
	if err := s.userDoAuthzAndAuditLog(ctx, id, pacta.AuditLogAction_ReadMetadata); err != nil {
		return nil, err
	}
	iurs, err := s.DB.InitiativeUserRelationshipsByUser(s.DB.NoTxn(ctx), id)
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
	id := pacta.InitiativeID(request.InitiativeId)
	if err := s.initiativeDoAuthzAndAuditLog(ctx, id, pacta.AuditLogAction_ReadMetadata); err != nil {
		return nil, err
	}
	iurs, err := s.DB.InitiativeUserRelationshipsByInitiative(s.DB.NoTxn(ctx), id)
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
	iID := pacta.InitiativeID(request.InitiativeId)
	uID := pacta.UserID(request.UserId)
	mutations := []db.UpdateInitiativeUserRelationshipFn{}
	if request.Body.Manager != nil {
		if err := s.initiativeUserRelationshipDoAuthzAndAuditLog(ctx, iID, uID, pacta.AuditLogAction_Update); err != nil {
			return nil, err
		}
		mutations = append(mutations, db.SetInitiativeUserRelationshipManager(*request.Body.Manager))
	}
	if request.Body.Member != nil {
		if *request.Body.Member {
			if err := s.initiativeUserRelationshipDoAuthzAndAuditLog(ctx, iID, uID, pacta.AuditLogAction_AddTo); err != nil {
				return nil, err
			}
		} else {
			if err := s.initiativeUserRelationshipDoAuthzAndAuditLog(ctx, iID, uID, pacta.AuditLogAction_RemoveFrom); err != nil {
				return nil, err
			}
		}
		mutations = append(mutations, db.SetInitiativeUserRelationshipMember(*request.Body.Member))
	}
	err := s.DB.UpdateInitiativeUserRelationship(s.DB.NoTxn(ctx), pacta.InitiativeID(request.InitiativeId), pacta.UserID(request.UserId), mutations...)
	if err != nil {
		return nil, oapierr.Internal("failed to update initiative user relationship", zap.Error(err))
	}
	return api.UpdateInitiativeUserRelationship204Response{}, nil
}

func (s *Server) initiativeUserRelationshipDoAuthzAndAuditLog(ctx context.Context, iID pacta.InitiativeID, targetUserID pacta.UserID, action pacta.AuditLogAction) error {
	actorInfo, err := s.getActorInfoOrFail(ctx)
	if err != nil {
		return err
	}
	i, err := s.DB.Initiative(s.DB.NoTxn(ctx), iID)
	if err != nil {
		if db.IsNotFound(err) {
			return notFoundOrUnauthorized(actorInfo, action, pacta.AuditLogTargetType_Initiative, iID)
		}
		return oapierr.Internal("failed to retrieve initiative", zap.Error(err))
	}
	iurs, err := s.DB.InitiativeUserRelationshipsByInitiative(s.DB.NoTxn(ctx), iID)
	if err != nil {
		return oapierr.Internal("failed to list initiative user relationships", zap.Error(err))
	}
	actorIsInitiativeManager := false
	for _, iur := range iurs {
		if iur.User.ID == actorInfo.UserID && iur.Manager {
			actorIsInitiativeManager = true
			break
		}
	}
	targetOwnerID, err := s.DB.GetOwnerForUser(s.DB.NoTxn(ctx), targetUserID)
	if err != nil {
		return err
	}
	targetIsActor := actorInfo.UserID == targetUserID
	as := &authzStatus{
		primaryTargetID:        string(iID),
		primaryTargetType:      pacta.AuditLogTargetType_Initiative,
		primaryTargetOwnerID:   systemOwnedEntityOwner,
		secondaryTargetID:      ptr(string(targetUserID)),
		secondaryTargetType:    ptr(pacta.AuditLogTargetType_User),
		secondaryTargetOwnerID: ptr(targetOwnerID),
		actorInfo:              actorInfo,
		action:                 action,
	}
	switch action {
	case pacta.AuditLogAction_AddTo:
		if i.IsAcceptingNewMembers {
			if actorIsInitiativeManager {
				as.authorizedAsActorType = ptr(pacta.AuditLogActorType_Owner)
				as.isAuthorized = true
			} else if targetIsActor && !i.RequiresInvitationToJoin {
				as.authorizedAsActorType = ptr(pacta.AuditLogActorType_Public)
				as.isAuthorized = true
			} else {
				as.isAuthorized, as.authorizedAsActorType = allowIfAdmin(actorInfo)
			}
		}
	case pacta.AuditLogAction_Update:
		if actorIsInitiativeManager {
			as.authorizedAsActorType = ptr(pacta.AuditLogActorType_Owner)
			as.isAuthorized = true
		} else {
			as.isAuthorized, as.authorizedAsActorType = allowIfAdmin(actorInfo)
		}
	case pacta.AuditLogAction_RemoveFrom:
		if actorIsInitiativeManager {
			as.authorizedAsActorType = ptr(pacta.AuditLogActorType_Owner)
			as.isAuthorized = true
		} else if targetIsActor {
			as.authorizedAsActorType = ptr(pacta.AuditLogActorType_Public)
			as.isAuthorized = true
		} else {
			as.isAuthorized, as.authorizedAsActorType = allowIfAdmin(actorInfo)
		}
	default:
		return fmt.Errorf("unknown action %q for initiative_invitation authz", action)
	}
	return s.auditLogIfAuthorizedOrFail(ctx, as)
}

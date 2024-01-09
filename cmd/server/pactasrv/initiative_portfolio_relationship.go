package pactasrv

import (
	"context"
	"fmt"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/oapierr"
	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
	"go.uber.org/zap"
)

// creates an initiative portfolio relationship
// (POST /initiative/{initiativeId}/portfolio-relationship/{portfolioId})
func (s *Server) CreateInitiativePortfolioRelationship(ctx context.Context, request api.CreateInitiativePortfolioRelationshipRequestObject) (api.CreateInitiativePortfolioRelationshipResponseObject, error) {
	pID := pacta.PortfolioID(request.PortfolioId)
	iID := pacta.InitiativeID(request.InitiativeId)
	if err := s.initiativePortfolioRelationshipDoAuthzAndAuditLog(ctx, iID, pID, pacta.AuditLogAction_AddTo); err != nil {
		return nil, err
	}
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}
	err = s.DB.CreatePortfolioInitiativeMembership(s.DB.NoTxn(ctx), &pacta.PortfolioInitiativeMembership{
		Portfolio:  &pacta.Portfolio{ID: pID},
		Initiative: &pacta.Initiative{ID: iID},
		AddedBy:    &pacta.User{ID: userID},
	})
	if err != nil {
		return nil, oapierr.Internal("failed to create initiative portfolio membership", zap.Error(err))
	}
	return api.CreateInitiativePortfolioRelationship204Response{}, nil
}

// Deletes an initiative:portfolio relationship
// (DELETE /initiative/{initiativeId}/portfolio-relationship/{portfolioId})
func (s *Server) DeleteInitiativePortfolioRelationship(ctx context.Context, request api.DeleteInitiativePortfolioRelationshipRequestObject) (api.DeleteInitiativePortfolioRelationshipResponseObject, error) {
	pID := pacta.PortfolioID(request.PortfolioId)
	iID := pacta.InitiativeID(request.InitiativeId)
	if err := s.initiativePortfolioRelationshipDoAuthzAndAuditLog(ctx, iID, pID, pacta.AuditLogAction_RemoveFrom); err != nil {
		return nil, err
	}
	err := s.DB.DeletePortfolioInitiativeMembership(s.DB.NoTxn(ctx), pID, iID)
	if err != nil {
		return nil, oapierr.Internal("failed to delete initiative-portfolio relationship",
			zap.String("initiative_id", request.InitiativeId),
			zap.String("portfolio_id", request.PortfolioId),
			zap.Error(err))
	}
	return api.DeleteInitiativePortfolioRelationship204Response{}, nil
}

func (s *Server) initiativePortfolioRelationshipDoAuthzAndAuditLog(ctx context.Context, iID pacta.InitiativeID, pID pacta.PortfolioID, action pacta.AuditLogAction) error {
	actorInfo, err := s.getActorInfoOrErrIfAnon(ctx)
	if err != nil {
		return err
	}
	i, err := s.DB.Initiative(s.DB.NoTxn(ctx), iID)
	if err != nil {
		if db.IsNotFound(err) {
			return notFoundOrUnauthorized(actorInfo, action, pacta.AuditLogTargetType_Initiative, iID)
		}
		return oapierr.Internal("failed to look up initiative", zap.String("initiative_id", string(iID)), zap.Error(err))
	}
	iurs, err := s.DB.InitiativeUserRelationshipsByInitiative(s.DB.NoTxn(ctx), iID)
	if err != nil {
		return oapierr.Internal("failed to list initiative user relationships", zap.Error(err))
	}
	actorIsInitiativeManager := false
	actorIsInitiativeMember := false
	for _, iur := range iurs {
		if iur.User.ID == actorInfo.UserID {
			if iur.Manager {
				actorIsInitiativeManager = true
			}
			if iur.Member {
				actorIsInitiativeMember = true
			}
			break
		}
	}
	p, err := s.DB.Portfolio(s.DB.NoTxn(ctx), pID)
	if err != nil {
		if db.IsNotFound(err) {
			return notFoundOrUnauthorized(actorInfo, action, pacta.AuditLogTargetType_Portfolio, pID)
		}
		return oapierr.Internal("failed to look up portfolio", zap.String("portfolio_id", string(pID)), zap.Error(err))
	}
	targetIsOwnedByActor := actorInfo.OwnerID == p.Owner.ID
	as := &authzStatus{
		primaryTargetID:        string(iID),
		primaryTargetType:      pacta.AuditLogTargetType_Initiative,
		primaryTargetOwnerID:   systemOwnedEntityOwner,
		secondaryTargetID:      string(pID),
		secondaryTargetType:    pacta.AuditLogTargetType_Portfolio,
		secondaryTargetOwnerID: p.Owner.ID,
		actorInfo:              actorInfo,
		action:                 action,
	}
	switch action {
	case pacta.AuditLogAction_AddTo:
		if i.IsAcceptingNewPortfolios {
			if actorIsInitiativeManager {
				as.authorizedAsActorType = ptr(pacta.AuditLogActorType_Owner)
				as.isAuthorized = true
			} else if actorIsInitiativeMember && targetIsOwnedByActor {
				as.authorizedAsActorType = ptr(pacta.AuditLogActorType_Public)
				as.isAuthorized = true
			} else {
				as.isAuthorized, as.authorizedAsActorType = allowIfAdmin(actorInfo)
			}
		}
	case pacta.AuditLogAction_RemoveFrom:
		if actorIsInitiativeManager {
			as.authorizedAsActorType = ptr(pacta.AuditLogActorType_Owner)
			as.isAuthorized = true
		} else if actorIsInitiativeMember && targetIsOwnedByActor {
			as.authorizedAsActorType = ptr(pacta.AuditLogActorType_Public)
			as.isAuthorized = true
		} else {
			as.isAuthorized, as.authorizedAsActorType = allowIfAdmin(actorInfo)
		}
	default:
		return fmt.Errorf("unknown action %q for initiative_portfolio_relationship authz", action)
	}
	return s.auditLogIfAuthorizedOrFail(ctx, as)
}

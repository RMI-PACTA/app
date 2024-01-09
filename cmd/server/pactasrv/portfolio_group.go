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

// Returns a portfolio group by ID
// (GET /portfolio-group/{id})
func (s *Server) FindPortfolioGroupById(ctx context.Context, request api.FindPortfolioGroupByIdRequestObject) (api.FindPortfolioGroupByIdResponseObject, error) {
	id := pacta.PortfolioGroupID(request.Id)
	if err := s.portfolioGroupAuthz(ctx, id, pacta.AuditLogAction_ReadMetadata); err != nil {
		return nil, err
	}
	pg, err := s.DB.PortfolioGroup(s.DB.NoTxn(ctx), id)
	if err != nil {
		if db.IsNotFound(err) {
			return nil, oapierr.NotFound("portfolio group not found", zap.String("portfolio_group_id", request.Id))
		}
		return nil, oapierr.Internal("failed to load portfolio_group", zap.String("portfolio_group_id", request.Id), zap.Error(err))
	}
	if err := s.populatePortfoliosInPortfolioGroups(ctx, []*pacta.PortfolioGroup{pg}); err != nil {
		return nil, err
	}
	resp, err := conv.PortfolioGroupToOAPI(pg)
	if err != nil {
		return nil, err
	}
	return api.FindPortfolioGroupById200JSONResponse(*resp), nil
}

// Returns the portfolio groups that the user has access to
// (GET /portfolio-groups)
func (s *Server) ListPortfolioGroups(ctx context.Context, request api.ListPortfolioGroupsRequestObject) (api.ListPortfolioGroupsResponseObject, error) {
	ownerID, err := s.getUserOwnerID(ctx)
	if err != nil {
		return nil, err
	}
	pgs, err := s.DB.PortfolioGroupsByOwner(s.DB.NoTxn(ctx), ownerID)
	if err != nil {
		return nil, oapierr.Internal("failed to query portfolio groups", zap.Error(err))
	}
	if err := s.populatePortfoliosInPortfolioGroups(ctx, pgs); err != nil {
		return nil, err
	}
	items, err := dereference(conv.PortfolioGroupsToOAPI(pgs))
	if err != nil {
		return nil, err
	}
	return api.ListPortfolioGroups200JSONResponse{Items: items}, nil
}

// Creates a portfolio group
// (POST /portfolio-groups)
func (s *Server) CreatePortfolioGroup(ctx context.Context, request api.CreatePortfolioGroupRequestObject) (api.CreatePortfolioGroupResponseObject, error) {
	actorInfo, err := s.getActorInfoOrErrIfAnon(ctx)
	if err != nil {
		return nil, err
	}
	pg, err := conv.PortfolioGroupCreateFromOAPI(request.Body, actorInfo.OwnerID)
	if err != nil {
		return nil, err
	}
	id, err := s.DB.CreatePortfolioGroup(s.DB.NoTxn(ctx), pg)
	if err != nil {
		return nil, oapierr.Internal("failed to create portfolio group", zap.Error(err))
	}
	pg.ID = id
	pg.CreatedAt = s.Now()
	resp, err := conv.PortfolioGroupToOAPI(pg)
	if err != nil {
		return nil, err
	}
	if err := s.auditLogForCreateEvent(ctx, actorInfo, pacta.AuditLogActorType_Owner, pacta.AuditLogTargetType_PortfolioGroup, string(id)); err != nil {
		return nil, err
	}
	return api.CreatePortfolioGroup200JSONResponse(*resp), nil
}

// Updates portfolio group properties
// (PATCH /portfolio-group/{id})
func (s *Server) UpdatePortfolioGroup(ctx context.Context, request api.UpdatePortfolioGroupRequestObject) (api.UpdatePortfolioGroupResponseObject, error) {
	id := pacta.PortfolioGroupID(request.Id)
	if err := s.portfolioGroupAuthz(ctx, id, pacta.AuditLogAction_Update); err != nil {
		return nil, err
	}
	mutations := []db.UpdatePortfolioGroupFn{}
	b := request.Body
	if b.Name != nil {
		mutations = append(mutations, db.SetPortfolioGroupName(*b.Name))
	}
	if b.Description != nil {
		mutations = append(mutations, db.SetPortfolioGroupDescription(*b.Description))
	}
	err := s.DB.UpdatePortfolioGroup(s.DB.NoTxn(ctx), id, mutations...)
	if err != nil {
		return nil, oapierr.Internal("failed to update portfolio group", zap.String("portfolio_group_id", string(id)), zap.Error(err))
	}
	return api.UpdatePortfolioGroup204Response{}, nil
}

// Deletes a portfolio group by ID
// (DELETE /portfolio-group/{id})
func (s *Server) DeletePortfolioGroup(ctx context.Context, request api.DeletePortfolioGroupRequestObject) (api.DeletePortfolioGroupResponseObject, error) {
	id := pacta.PortfolioGroupID(request.Id)
	if err := s.portfolioGroupAuthz(ctx, id, pacta.AuditLogAction_Delete); err != nil {
		return nil, err
	}
	err := s.DB.DeletePortfolioGroup(s.DB.NoTxn(ctx), id)
	if err != nil {
		return nil, oapierr.Internal("failed to delete portfolio group", zap.String("portfolio_group_id", request.Id), zap.Error(err))
	}
	return api.DeletePortfolioGroup204Response{}, nil
}

// Deletes a portfolio group membership
// (DELETE /portfolio-group-membership)
func (s *Server) DeletePortfolioGroupMembership(ctx context.Context, request api.DeletePortfolioGroupMembershipRequestObject) (api.DeletePortfolioGroupMembershipResponseObject, error) {
	pgID := pacta.PortfolioGroupID(request.Body.PortfolioGroupId)
	pID := pacta.PortfolioID(request.Body.PortfolioId)
	if err := s.portfolioGroupMembershipAuthz(ctx, pgID, pacta.AuditLogAction_RemoveFrom, pID); err != nil {
		return nil, err
	}
	err := s.DB.DeletePortfolioGroupMembership(s.DB.NoTxn(ctx), pgID, pID)
	if err != nil {
		return nil, oapierr.Internal("failed to delete portfolio group membership", zap.String("portfolio_group_id", string(pgID)), zap.String("portfolio_id", string(pID)), zap.Error(err))
	}
	return api.DeletePortfolioGroupMembership204Response{}, nil
}

// creates a portfolio group membership
// (PUT /portfolio-group-membership)
func (s *Server) CreatePortfolioGroupMembership(ctx context.Context, request api.CreatePortfolioGroupMembershipRequestObject) (api.CreatePortfolioGroupMembershipResponseObject, error) {
	pgID := pacta.PortfolioGroupID(request.Body.PortfolioGroupId)
	pID := pacta.PortfolioID(request.Body.PortfolioId)
	if err := s.portfolioGroupMembershipAuthz(ctx, pgID, pacta.AuditLogAction_AddTo, pID); err != nil {
		return nil, err
	}
	err := s.DB.CreatePortfolioGroupMembership(s.DB.NoTxn(ctx), pgID, pID)
	if err != nil {
		return nil, oapierr.Internal("failed to create portfolio group membership", zap.String("portfolio_group_id", string(pgID)), zap.String("portfolio_id", string(pID)), zap.Error(err))
	}
	return api.CreatePortfolioGroupMembership204Response{}, nil
}

func (s *Server) portfolioGroupAuthz(ctx context.Context, pgID pacta.PortfolioGroupID, action pacta.AuditLogAction) error {
	actorInfo, err := s.getActorInfoOrErrIfAnon(ctx)
	if err != nil {
		return err
	}
	pg, err := s.DB.PortfolioGroup(s.DB.NoTxn(ctx), pgID)
	if err != nil {
		if db.IsNotFound(err) {
			return notFoundOrUnauthorized(actorInfo, action, pacta.AuditLogTargetType_PortfolioGroup, pgID)
		}
		return oapierr.Internal("failed to look up portfolio_group", zap.String("portfolio_group_id", string(pgID)), zap.Error(err))
	}
	as := &authzStatus{
		primaryTargetID:      string(pgID),
		primaryTargetType:    pacta.AuditLogTargetType_PortfolioGroup,
		primaryTargetOwnerID: pg.Owner.ID,
		actorInfo:            actorInfo,
		action:               action,
	}
	switch action {
	case pacta.AuditLogAction_Update, pacta.AuditLogAction_Delete, pacta.AuditLogAction_ReadMetadata:
		as.isAuthorized, as.authorizedAsActorType = allowIfAdminOrOwner(actorInfo, pg.Owner.ID)
	default:
		return fmt.Errorf("unknown action %q for portfolio_group authz", action)
	}
	return s.auditLogIfAuthorizedOrFail(ctx, as)
}

func (s *Server) portfolioGroupMembershipAuthz(ctx context.Context, pgID pacta.PortfolioGroupID, action pacta.AuditLogAction, pID pacta.PortfolioID) error {
	actorInfo, err := s.getActorInfoOrErrIfAnon(ctx)
	if err != nil {
		return err
	}
	pg, err := s.DB.PortfolioGroup(s.DB.NoTxn(ctx), pgID)
	if err != nil {
		if db.IsNotFound(err) {
			return notFoundOrUnauthorized(actorInfo, action, pacta.AuditLogTargetType_PortfolioGroup, pgID)
		}
		return oapierr.Internal("failed to look up portfolio_group", zap.String("portfolio_group_id", string(pgID)), zap.Error(err))
	}
	p, err := s.DB.Portfolio(s.DB.NoTxn(ctx), pID)
	if err != nil {
		if db.IsNotFound(err) {
			return notFoundOrUnauthorized(actorInfo, action, pacta.AuditLogTargetType_Portfolio, pID)
		}
		return oapierr.Internal("failed to look up portfolio for pgm", zap.String("portfolio_id", string(pID)), zap.Error(err))
	}
	as := &authzStatus{
		primaryTargetID:        string(pgID),
		primaryTargetType:      pacta.AuditLogTargetType_PortfolioGroup,
		primaryTargetOwnerID:   pg.Owner.ID,
		secondaryTargetID:      ptr(string(pID)),
		secondaryTargetType:    ptr(pacta.AuditLogTargetType_Portfolio),
		secondaryTargetOwnerID: ptr(p.Owner.ID),
		actorInfo:              actorInfo,
		action:                 action,
	}
	switch action {
	case pacta.AuditLogAction_AddTo, pacta.AuditLogAction_RemoveFrom:
		// NOTE! The actor must be the owner of BOTH the portfolio group and the portfolio in order to add/remove it.
		if actorInfo.OwnerID == pg.Owner.ID && actorInfo.OwnerID == p.Owner.ID {
			as.isAuthorized = true
			as.authorizedAsActorType = ptr(pacta.AuditLogActorType_Owner)
		} else {
			as.isAuthorized, as.authorizedAsActorType = allowIfAdminOrOwner(actorInfo, pg.Owner.ID)
		}
	default:
		return fmt.Errorf("unknown action %q for portfolio_group_membership authz", action)
	}
	return s.auditLogIfAuthorizedOrFail(ctx, as)
}

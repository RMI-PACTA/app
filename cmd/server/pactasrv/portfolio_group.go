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

// Returns a portfolio group by ID
// (GET /portfolio-group/{id})
func (s *Server) FindPortfolioGroupById(ctx context.Context, request api.FindPortfolioGroupByIdRequestObject) (api.FindPortfolioGroupByIdResponseObject, error) {
	// TODO(#12) Implement Authorization
	pg, err := s.DB.PortfolioGroup(s.DB.NoTxn(ctx), pacta.PortfolioGroupID(request.Id))
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
	// TODO(#12) Implement Authorization
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
	ownerID, err := s.getUserOwnerID(ctx)
	if err != nil {
		return nil, err
	}
	pg, err := conv.PortfolioGroupCreateFromOAPI(request.Body, ownerID)
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
	return api.CreatePortfolioGroup200JSONResponse(*resp), nil
}

// Updates portfolio group properties
// (PATCH /portfolio-group/{id})
func (s *Server) UpdatePortfolioGroup(ctx context.Context, request api.UpdatePortfolioGroupRequestObject) (api.UpdatePortfolioGroupResponseObject, error) {
	// TODO(#12) Implement Authorization
	id := pacta.PortfolioGroupID(request.Id)
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
	// TODO(#12) Implement Authorization
	err := s.DB.DeletePortfolioGroup(s.DB.NoTxn(ctx), pacta.PortfolioGroupID(request.Id))
	if err != nil {
		return nil, oapierr.Internal("failed to delete portfolio group", zap.String("portfolio_group_id", request.Id), zap.Error(err))
	}
	return api.DeletePortfolioGroup204Response{}, nil
}

// Deletes a portfolio group membership
// (DELETE /portfolio-group-membership)
func (s *Server) DeletePortfolioGroupMembership(ctx context.Context, request api.DeletePortfolioGroupMembershipRequestObject) (api.DeletePortfolioGroupMembershipResponseObject, error) {
	// TODO(#12) Implement Authorization
	pgID := pacta.PortfolioGroupID(request.Body.PortfolioGroupId)
	pID := pacta.PortfolioID(request.Body.PortfolioId)
	err := s.DB.DeletePortfolioGroupMembership(s.DB.NoTxn(ctx), pgID, pID)
	if err != nil {
		return nil, oapierr.Internal("failed to delete portfolio group membership", zap.String("portfolio_group_id", string(pgID)), zap.String("portfolio_id", string(pID)), zap.Error(err))
	}
	return api.DeletePortfolioGroupMembership204Response{}, nil
}

// creates a portfolio group membership
// (PUT /portfolio-group-membership)
func (s *Server) CreatePortfolioGroupMembership(ctx context.Context, request api.CreatePortfolioGroupMembershipRequestObject) (api.CreatePortfolioGroupMembershipResponseObject, error) {
	// TODO(#12) Implement Authorization
	pgID := pacta.PortfolioGroupID(request.Body.PortfolioGroupId)
	pID := pacta.PortfolioID(request.Body.PortfolioId)
	err := s.DB.CreatePortfolioGroupMembership(s.DB.NoTxn(ctx), pgID, pID)
	if err != nil {
		return nil, oapierr.Internal("failed to create portfolio group membership", zap.String("portfolio_group_id", string(pgID)), zap.String("portfolio_id", string(pID)), zap.Error(err))
	}
	return api.CreatePortfolioGroupMembership204Response{}, nil
}

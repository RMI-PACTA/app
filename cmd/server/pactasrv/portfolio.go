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

// (GET /portfolios)
func (s *Server) ListPortfolios(ctx context.Context, request api.ListPortfoliosRequestObject) (api.ListPortfoliosResponseObject, error) {
	ownerID, err := s.getUserOwnerID(ctx)
	if err != nil {
		return nil, err
	}
	ps, err := s.DB.PortfoliosByOwner(s.DB.NoTxn(ctx), ownerID)
	if err != nil {
		return nil, oapierr.Internal("failed to query portfolios", zap.Error(err))
	}
	if err := s.populatePortfolioGroupsInPortfolios(ctx, ps); err != nil {
		return nil, err
	}
	if err := s.populateInitiativesInPortfolios(ctx, ps); err != nil {
		return nil, err
	}
	items, err := dereference(conv.PortfoliosToOAPI(ps))
	if err != nil {
		return nil, err
	}
	return api.ListPortfolios200JSONResponse{Items: items}, nil
}

// Deletes an portfolio by ID
// (DELETE /portfolio/{id})
func (s *Server) DeletePortfolio(ctx context.Context, request api.DeletePortfolioRequestObject) (api.DeletePortfolioResponseObject, error) {
	id := pacta.PortfolioID(request.Id)
	if err := s.portfolioDoAuthzAndAuditLog(ctx, id, pacta.AuditLogAction_Delete); err != nil {
		return nil, err
	}
	blobURIs, err := s.DB.DeletePortfolio(s.DB.NoTxn(ctx), id)
	if err != nil {
		return nil, oapierr.Internal("failed to delete portfolio", zap.Error(err))
	}
	if err := s.deleteBlobs(ctx, blobURIs...); err != nil {
		return nil, err
	}
	return api.DeletePortfolio204Response{}, nil
}

// Returns an portfolio by ID
// (GET /portfolio/{id})
func (s *Server) FindPortfolioById(ctx context.Context, request api.FindPortfolioByIdRequestObject) (api.FindPortfolioByIdResponseObject, error) {
	id := pacta.PortfolioID(request.Id)
	if err := s.portfolioDoAuthzAndAuditLog(ctx, id, pacta.AuditLogAction_ReadMetadata); err != nil {
		return nil, err
	}
	p, err := s.DB.Portfolio(s.DB.NoTxn(ctx), id)
	if err != nil {
		return nil, oapierr.Internal("failed to look up portfolio", zap.String("portfolio_id", string(id)), zap.Error(err))
	}
	if err := s.populatePortfolioGroupsInPortfolios(ctx, []*pacta.Portfolio{p}); err != nil {
		return nil, err
	}
	if err := s.populateInitiativesInPortfolios(ctx, []*pacta.Portfolio{p}); err != nil {
		return nil, err
	}
	converted, err := conv.PortfolioToOAPI(p)
	if err != nil {
		return nil, err
	}
	return api.FindPortfolioById200JSONResponse(*converted), nil
}

// Updates portfolio properties
// (PATCH /portfolio/{id})
func (s *Server) UpdatePortfolio(ctx context.Context, request api.UpdatePortfolioRequestObject) (api.UpdatePortfolioResponseObject, error) {
	id := pacta.PortfolioID(request.Id)
	if err := s.portfolioDoAuthzAndAuditLog(ctx, id, pacta.AuditLogAction_Update); err != nil {
		return nil, err
	}
	mutations := []db.UpdatePortfolioFn{}
	if request.Body.Name != nil {
		mutations = append(mutations, db.SetPortfolioName(*request.Body.Name))
	}
	if request.Body.Description != nil {
		mutations = append(mutations, db.SetPortfolioDescription(*request.Body.Description))
	}
	if request.Body.AdminDebugEnabled != nil {
		if *request.Body.AdminDebugEnabled {
			if err := s.portfolioDoAuthzAndAuditLog(ctx, id, pacta.AuditLogAction_EnableAdminDebug); err != nil {
				return nil, err
			}
		} else {
			if err := s.portfolioDoAuthzAndAuditLog(ctx, id, pacta.AuditLogAction_DisableAdminDebug); err != nil {
				return nil, err
			}
		}
		mutations = append(mutations, db.SetPortfolioAdminDebugEnabled(*request.Body.AdminDebugEnabled))
	}
	err := s.DB.UpdatePortfolio(s.DB.NoTxn(ctx), id, mutations...)
	if err != nil {
		return nil, oapierr.Internal("failed to update portfolio", zap.Error(err))
	}
	return api.UpdatePortfolio204Response{}, nil
}

func (s *Server) portfolioDoAuthzAndAuditLog(ctx context.Context, pID pacta.PortfolioID, action pacta.AuditLogAction) error {
	actorInfo, err := s.getActorInfoOrFail(ctx)
	if err != nil {
		return err
	}
	p, err := s.DB.Portfolio(s.DB.NoTxn(ctx), pID)
	if err != nil {
		if db.IsNotFound(err) {
			return notFoundOrUnauthorized(actorInfo, action, pacta.AuditLogTargetType_Portfolio, pID)
		}
		return oapierr.Internal("failed to look up portfolio", zap.String("portfolio_id", string(pID)), zap.Error(err))
	}
	as := &authzStatus{
		primaryTargetID:      string(pID),
		primaryTargetType:    pacta.AuditLogTargetType_Portfolio,
		primaryTargetOwnerID: p.Owner.ID,
		actorInfo:            actorInfo,
		action:               action,
	}
	switch action {
	case pacta.AuditLogAction_EnableAdminDebug, pacta.AuditLogAction_DisableAdminDebug:
		as.isAuthorized, as.authorizedAsActorType = allowIfOwner(actorInfo, p.Owner.ID)
	case pacta.AuditLogAction_Update, pacta.AuditLogAction_Delete, pacta.AuditLogAction_ReadMetadata:
		as.isAuthorized, as.authorizedAsActorType = allowIfAdminOrOwner(actorInfo, p.Owner.ID)
	default:
		return fmt.Errorf("unknown action %q for portfolio authz", action)
	}
	return s.auditLogIfAuthorizedOrFail(ctx, as)
}

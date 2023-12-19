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
	"go.uber.org/zap/zapcore"
)

// (GET /portfolios)
func (s *Server) ListPortfolios(ctx context.Context, request api.ListPortfoliosRequestObject) (api.ListPortfoliosResponseObject, error) {
	// TODO(#12) Implement Authorization
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
	_, err := s.checkPortfolioAuthorization(ctx, id)
	if err != nil {
		return nil, err
	}
	blobURIs, err := s.DB.DeletePortfolio(s.DB.NoTxn(ctx), id)
	if err != nil {
		return nil, oapierr.Internal("failed to delete portfolio", zap.Error(err))
	}
	if err := s.deleteBlobs(ctx, asStrs(blobURIs)); err != nil {
		return nil, oapierr.Internal("failed to delete blob", zap.Error(err))
	}
	return api.DeletePortfolio204Response{}, nil
}

// Returns an portfolio by ID
// (GET /portfolio/{id})
func (s *Server) FindPortfolioById(ctx context.Context, request api.FindPortfolioByIdRequestObject) (api.FindPortfolioByIdResponseObject, error) {
	p, err := s.checkPortfolioAuthorization(ctx, pacta.PortfolioID(request.Id))
	if err != nil {
		return nil, err
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
	_, err := s.checkPortfolioAuthorization(ctx, id)
	if err != nil {
		return nil, err
	}
	// TODO(#12) Implement Authorization
	mutations := []db.UpdatePortfolioFn{}
	if request.Body.Name != nil {
		mutations = append(mutations, db.SetPortfolioName(*request.Body.Name))
	}
	if request.Body.Description != nil {
		mutations = append(mutations, db.SetPortfolioDescription(*request.Body.Description))
	}
	if request.Body.AdminDebugEnabled != nil {
		mutations = append(mutations, db.SetPortfolioAdminDebugEnabled(*request.Body.AdminDebugEnabled))
	}
	err = s.DB.UpdatePortfolio(s.DB.NoTxn(ctx), id, mutations...)
	if err != nil {
		return nil, oapierr.Internal("failed to update portfolio", zap.Error(err))
	}
	return api.UpdatePortfolio204Response{}, nil
}

func (s *Server) checkPortfolioAuthorization(ctx context.Context, id pacta.PortfolioID) (*pacta.Portfolio, error) {
	// TODO(#12) Implement Authorization
	actorOwnerID, err := s.getUserOwnerID(ctx)
	if err != nil {
		return nil, err
	}
	// Extracted to a common variable so that we return the same response for not found and unauthorized.
	notFoundErr := func(fields ...zapcore.Field) error {
		fs := append(fields, zap.String("portfolio_id", string(id)))
		return oapierr.NotFound("portfolio not found", fs...)
	}
	iu, err := s.DB.Portfolio(s.DB.NoTxn(ctx), id)
	if err != nil {
		if db.IsNotFound(err) {
			return nil, notFoundErr(zap.Error(err))
		}
		return nil, oapierr.Internal("failed to look up portfolio", zap.String("portfolio_id", string(id)), zap.Error(err))
	}
	if iu.Owner.ID != actorOwnerID {
		return nil, notFoundErr(zap.Error(fmt.Errorf("portfolio does not belong to user")), zap.String("owner", string(iu.Owner.ID)), zap.String("actor", string(actorOwnerID)))
	}
	return iu, nil
}

package pactasrv

import (
	"context"

	"github.com/RMI/pacta/oapierr"
	api "github.com/RMI/pacta/openapi/pacta"
)

func (s *Server) CreatePortfolioAsset(ctx context.Context, req api.CreatePortfolioAssetRequestObject) (api.CreatePortfolioAssetResponseObject, error) {
	return nil, oapierr.NotImplemented("no longer implemented")
}

func (s *Server) ParsePortfolio(ctx context.Context, req api.ParsePortfolioRequestObject) (api.ParsePortfolioResponseObject, error) {
	return nil, oapierr.NotImplemented("no longer implemented")
}

// (GET /portfolios)

func (s *Server) ListPortfolios(ctx context.Context, request api.ListPortfoliosRequestObject) (api.ListPortfoliosResponseObject, error) {
	return nil, oapierr.NotImplemented("not implemented")
}

package pactasrv

import (
	"context"

	"github.com/RMI/pacta/oapierr"
	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
	"go.uber.org/zap"
)

// creates an initiative portfolio relationship
// (POST /initiative/{initiativeId}/portfolio-relationship/{portfolioId})
func (s *Server) CreateInitiativePortfolioRelationship(ctx context.Context, request api.CreateInitiativePortfolioRelationshipRequestObject) (api.CreateInitiativePortfolioRelationshipResponseObject, error) {
	// TODO(#12) Implement Authorization
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}
	err = s.DB.CreatePortfolioInitiativeMembership(s.DB.NoTxn(ctx), &pacta.PortfolioInitiativeMembership{
		Portfolio:  &pacta.Portfolio{ID: pacta.PortfolioID(request.PortfolioId)},
		Initiative: &pacta.Initiative{ID: pacta.InitiativeID(request.InitiativeId)},
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
	// TODO(#12) Implement Authorization
	err := s.DB.DeletePortfolioInitiativeMembership(s.DB.NoTxn(ctx), pacta.PortfolioID(request.PortfolioId), pacta.InitiativeID(request.InitiativeId))
	if err != nil {
		return nil, oapierr.Internal("failed to delete initiative-portfolio relationship",
			zap.String("initiative_id", request.InitiativeId),
			zap.String("portfolio_id", request.PortfolioId),
			zap.Error(err))
	}
	return api.DeleteInitiativePortfolioRelationship204Response{}, nil
}

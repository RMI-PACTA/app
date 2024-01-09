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

// Creates a initiative
// (POST /initiatives)
func (s *Server) CreateInitiative(ctx context.Context, request api.CreateInitiativeRequestObject) (api.CreateInitiativeResponseObject, error) {
	// TODO(#12) Implement Authorization
	i, err := conv.InitiativeCreateFromOAPI(request.Body)
	if err != nil {
		return nil, err
	}
	err = s.DB.CreateInitiative(s.DB.NoTxn(ctx), i)
	if err != nil {
		return nil, oapierr.Internal("failed to create initiative", zap.Error(err))
	}
	return api.CreateInitiative204Response{}, nil
}

// Updates an initiative
// (PATCH /initiative/{id})
func (s *Server) UpdateInitiative(ctx context.Context, request api.UpdateInitiativeRequestObject) (api.UpdateInitiativeResponseObject, error) {
	// TODO(#12) Implement Authorization
	id := pacta.InitiativeID(request.Id)
	mutations := []db.UpdateInitiativeFn{}
	b := request.Body
	if b.Affiliation != nil {
		mutations = append(mutations, db.SetInitiativeAffiliation(*b.Affiliation))
	}
	if b.InternalDescription != nil {
		mutations = append(mutations, db.SetInitiativeInternalDescription(*b.InternalDescription))
	}
	if b.IsAcceptingNewMembers != nil {
		mutations = append(mutations, db.SetInitiativeIsAcceptingNewMembers(*b.IsAcceptingNewMembers))
	}
	if b.IsAcceptingNewPortfolios != nil {
		mutations = append(mutations, db.SetInitiativeIsAcceptingNewPortfolios(*b.IsAcceptingNewPortfolios))
	}
	if b.Language != nil {
		lang, err := conv.LanguageFromOAPI(*b.Language)
		if err != nil {
			return nil, oapierr.BadRequest("failed to parse language", zap.Error(err))
		}
		mutations = append(mutations, db.SetInitiativeLanguage(lang))
	}
	if b.Name != nil {
		mutations = append(mutations, db.SetInitiativeName(*b.Name))
	}
	if b.PactaVersion != nil {
		mutations = append(mutations, db.SetInitiativePACTAVersion(pacta.PACTAVersionID(*b.PactaVersion)))
	}
	if b.PublicDescription != nil {
		mutations = append(mutations, db.SetInitiativePublicDescription(*b.PublicDescription))
	}
	if b.RequiresInvitationToJoin != nil {
		mutations = append(mutations, db.SetInitiativeRequiresInvitationToJoin(*b.RequiresInvitationToJoin))
	}
	err := s.DB.UpdateInitiative(s.DB.NoTxn(ctx), id, mutations...)
	if err != nil {
		return nil, oapierr.Internal("failed to update initiative", zap.String("initiative_id", string(id)), zap.Error(err))
	}
	return api.UpdateInitiative204Response{}, nil
}

// Deletes an initiative by id
// (DELETE /initiative/{id})
func (s *Server) DeleteInitiative(ctx context.Context, request api.DeleteInitiativeRequestObject) (api.DeleteInitiativeResponseObject, error) {
	// TODO(#12) Implement Authorization
	err := s.DB.DeleteInitiative(s.DB.NoTxn(ctx), pacta.InitiativeID(request.Id))
	if err != nil {
		return nil, oapierr.Internal("failed to delete initiative", zap.Error(err))
	}
	return api.DeleteInitiative204Response{}, nil
}

// Returns an initiative by ID
// (GET /initiative/{id})
func (s *Server) FindInitiativeById(ctx context.Context, request api.FindInitiativeByIdRequestObject) (api.FindInitiativeByIdResponseObject, error) {
	// TODO(#12) Implement Authorization
	i, err := s.DB.Initiative(s.DB.NoTxn(ctx), pacta.InitiativeID(request.Id))
	if err != nil {
		if db.IsNotFound(err) {
			return nil, oapierr.NotFound("initiative not found", zap.String("initiative_id", request.Id))
		}
		return nil, oapierr.Internal("failed to load initiative", zap.String("initiative_id", request.Id), zap.Error(err))
	}
	portfolios, err := s.DB.PortfolioInitiativeMembershipsByInitiative(s.DB.NoTxn(ctx), i.ID)
	if err != nil {
		return nil, oapierr.Internal("failed to load portfolios for initiative", zap.String("initiative_id", string(i.ID)), zap.Error(err))
	}
	i.PortfolioInitiativeMemberships = portfolios
	resp, err := conv.InitiativeToOAPI(i)
	if err != nil {
		return nil, err
	}
	return api.FindInitiativeById200JSONResponse(*resp), nil
}

// Returns all initiatives
// (GET /initiatives)
func (s *Server) ListInitiatives(ctx context.Context, request api.ListInitiativesRequestObject) (api.ListInitiativesResponseObject, error) {
	is, err := s.DB.AllInitiatives(s.DB.NoTxn(ctx))
	if err != nil {
		return nil, oapierr.Internal("failed to load all initiatives", zap.Error(err))
	}
	result, err := dereference(mapAll(is, conv.InitiativeToOAPI))
	if err != nil {
		return nil, err
	}
	return api.ListInitiatives200JSONResponse(result), nil
}

// Returns all of the portfolios that are participating in the initiative
// (GET /initiative/{id}/all-data)
func (s *Server) AllInitiativeData(ctx context.Context, request api.AllInitiativeDataRequestObject) (api.AllInitiativeDataResponseObject, error) {
	actorInfo, err := s.getactorInfoOrErrIfAnon(ctx)
	if err != nil {
		return nil, err
	}
	// TODO(#12) Implement Authorization, along the lines of #121
	i, err := s.DB.Initiative(s.DB.NoTxn(ctx), pacta.InitiativeID(request.Id))
	if err != nil {
		if db.IsNotFound(err) {
			return nil, oapierr.NotFound("initiative not found", zap.String("initiative_id", request.Id))
		}
		return nil, oapierr.Internal("failed to load initiative", zap.String("initiative_id", request.Id), zap.Error(err))
	}
	portfolioMembers, err := s.DB.PortfolioInitiativeMembershipsByInitiative(s.DB.NoTxn(ctx), i.ID)
	if err != nil {
		return nil, oapierr.Internal("failed to load portfolio memberships for initiative", zap.String("initiative_id", string(i.ID)), zap.Error(err))
	}
	portfolioIDs := []pacta.PortfolioID{}
	for _, pm := range portfolioMembers {
		portfolioIDs = append(portfolioIDs, pm.Portfolio.ID)
	}
	portfolios, err := s.DB.Portfolios(s.DB.NoTxn(ctx), portfolioIDs)
	if err != nil {
		return nil, oapierr.Internal("failed to load portfolios for initiative", zap.String("initiative_id", string(i.ID)), zap.Error(err))
	}
	if err := s.populateBlobsInPortfolios(ctx, values(portfolios)...); err != nil {
		return nil, err
	}
	auditLogs := []*pacta.AuditLog{}
	for _, p := range portfolios {
		auditLogs = append(auditLogs, &pacta.AuditLog{
			Action:               pacta.AuditLogAction_Download,
			ActorType:            pacta.AuditLogActorType_Admin, // TODO(#12) When merging with #121, use the actor type from authorization
			ActorID:              string(actorInfo.UserID),
			ActorOwner:           &pacta.Owner{ID: actorInfo.OwnerID},
			PrimaryTargetType:    pacta.AuditLogTargetType_Portfolio,
			PrimaryTargetID:      string(p.ID),
			PrimaryTargetOwner:   p.Owner,
			SecondaryTargetType:  pacta.AuditLogTargetType_Initiative,
			SecondaryTargetID:    string(i.ID),
			SecondaryTargetOwner: &pacta.Owner{ID: "SYSTEM"}, // TODO(#12) When merging with #121, use the const type.
		})
	}
	if err := s.DB.CreateAuditLogs(s.DB.NoTxn(ctx), auditLogs); err != nil {
		return nil, oapierr.Internal("failed to create audit logs nescessary to return download urls", zap.Error(err))
	}

	// Note, it is likely this code will need to be parallelized in the future - initiatives may eventually become large.
	// However, since this action is unlikely to be taken frequently, and will only be taken by admins, getting the experience
	// perfect here is not a priority.
	response := api.InitiativeAllData{}
	for _, portfolio := range portfolios {
		url, expiryTime, err := s.Blob.SignedDownloadURL(ctx, string(portfolio.Blob.BlobURI))
		if err != nil {
			return nil, oapierr.Internal("error getting signed download url", zap.Error(err), zap.String("blob_uri", string(portfolio.Blob.BlobURI)))
		}
		response.Items = append(response.Items, api.InitiativeAllDataPortfolioItem{
			Name:           portfolio.Name,
			BlobId:         string(portfolio.Blob.BlobURI),
			DownloadUrl:    url,
			ExpirationTime: expiryTime,
		})
	}

	return api.AllInitiativeData200JSONResponse(response), nil
}

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

// Creates a initiative
// (POST /initiatives)
func (s *Server) CreateInitiative(ctx context.Context, request api.CreateInitiativeRequestObject) (api.CreateInitiativeResponseObject, error) {
	if err := anyError(
		checkStringLimitSmall("id", request.Body.Id),
		checkStringLimitSmallPtr("affiliation", request.Body.Affiliation),
		checkStringLimitSmall("name", request.Body.Name),
		checkStringLimitMediumPtr("internal_description", request.Body.InternalDescription),
		checkStringLimitMediumPtr("public_description", request.Body.PublicDescription),
	); err != nil {
		return nil, err
	}
	actorInfo, err := s.getActorInfoOrErrIfAnon(ctx)
	if err != nil {
		return nil, err
	}
	var auditLogActorType pacta.AuditLogActorType
	if actorInfo.IsAdmin {
		auditLogActorType = pacta.AuditLogActorType_Admin
	} else if actorInfo.IsSuperAdmin {
		auditLogActorType = pacta.AuditLogActorType_SuperAdmin
	} else {
		return nil, oapierr.Forbidden("only admins and super admins can create initiatives")
	}
	i, err := conv.InitiativeCreateFromOAPI(request.Body)
	if err != nil {
		return nil, err
	}
	err = s.DB.CreateInitiative(s.DB.NoTxn(ctx), i)
	if err != nil {
		return nil, oapierr.Internal("failed to create initiative", zap.Error(err))
	}
	if err := s.auditLogForCreateEvent(
		ctx,
		actorInfo,
		auditLogActorType,
		pacta.AuditLogTargetType_Initiative,
		string(i.ID)); err != nil {
		return nil, err
	}
	return api.CreateInitiative204Response{}, nil
}

// Updates an initiative
// (PATCH /initiative/{id})
func (s *Server) UpdateInitiative(ctx context.Context, request api.UpdateInitiativeRequestObject) (api.UpdateInitiativeResponseObject, error) {
	if err := anyError(
		checkStringLimitSmallPtr("affiliation", request.Body.Affiliation),
		checkStringLimitSmallPtr("name", request.Body.Name),
		checkStringLimitMediumPtr("internal_description", request.Body.InternalDescription),
		checkStringLimitMediumPtr("public_description", request.Body.PublicDescription),
	); err != nil {
		return nil, err
	}
	id := pacta.InitiativeID(request.Id)
	if _, err := s.initiativeDoAuthzAndAuditLog(ctx, id, pacta.AuditLogAction_Update); err != nil {
		return nil, err
	}
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
	id := pacta.InitiativeID(request.Id)
	if _, err := s.initiativeDoAuthzAndAuditLog(ctx, id, pacta.AuditLogAction_Delete); err != nil {
		return nil, err
	}
	buris, err := s.DB.DeleteInitiative(s.DB.NoTxn(ctx), id)
	if err != nil {
		return nil, oapierr.Internal("failed to delete initiative", zap.Error(err))
	}
	if err := s.deleteBlobs(ctx, buris...); err != nil {
		return nil, err
	}
	return api.DeleteInitiative204Response{}, nil
}

// Returns an initiative by ID
// (GET /initiative/{id})
func (s *Server) FindInitiativeById(ctx context.Context, request api.FindInitiativeByIdRequestObject) (api.FindInitiativeByIdResponseObject, error) {
	actorInfo, err := s.getActorInfoOrAnon(ctx)
	if err != nil {
		return nil, err
	}
	id := pacta.InitiativeID(request.Id)
	info, err := s.initiativeDoAuthzAndAuditLog(ctx, id, pacta.AuditLogAction_ReadMetadata)
	if err != nil {
		return nil, err
	}
	i, err := s.DB.Initiative(s.DB.NoTxn(ctx), id)
	if err != nil {
		if db.IsNotFound(err) {
			return nil, oapierr.NotFound("initiative not found", zap.String("initiative_id", request.Id))
		}
		return nil, oapierr.Internal("failed to load initiative", zap.String("initiative_id", request.Id), zap.Error(err))
	}
	if info.CanManageUsersAndPortfolios {
		portfolios, err := s.DB.PortfolioInitiativeMembershipsByInitiative(s.DB.NoTxn(ctx), i.ID)
		if err != nil {
			return nil, oapierr.Internal("failed to load portfolios for initiative", zap.String("initiative_id", string(i.ID)), zap.Error(err))
		}
		i.PortfolioInitiativeMemberships = portfolios
		if err := s.populatePortfoliosInInitiatives(ctx, []*pacta.Initiative{i}); err != nil {
			return nil, err
		}
		relationships, err := s.DB.InitiativeUserRelationshipsByInitiative(s.DB.NoTxn(ctx), i.ID)
		if err != nil {
			return nil, oapierr.Internal("failed to load initiative user relationships for initiative", zap.String("initiative_id", string(i.ID)), zap.Error(err))
		}
		i.InitiativeUserRelationships = relationships
	} else if actorInfo.UserID != "" {
		iur, err := s.DB.InitiativeUserRelationship(s.DB.NoTxn(ctx), i.ID, actorInfo.UserID)
		if err != nil {
			return nil, oapierr.Internal("failed to load singular initiative user relationship for initiative", zap.String("initiative_id", string(i.ID)), zap.Error(err))
		}
		i.InitiativeUserRelationships = []*pacta.InitiativeUserRelationship{iur}
	}
	if !info.CanSeeInternalInfo {
		i.InternalDescription = ""
	}
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
	for _, i := range is {
		i.InternalDescription = ""
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
	actorInfo, err := s.getActorInfoOrErrIfAnon(ctx)
	if err != nil {
		return nil, err
	}
	id := pacta.InitiativeID(request.Id)
	_, err = s.initiativeDoAuthzAndAuditLog(ctx, id, pacta.AuditLogAction_Download)
	if err != nil {
		return nil, err
	}
	portfolioMembers, err := s.DB.PortfolioInitiativeMembershipsByInitiative(s.DB.NoTxn(ctx), id)
	if err != nil {
		return nil, oapierr.Internal("failed to load portfolio memberships for initiative", zap.String("initiative_id", string(id)), zap.Error(err))
	}
	portfolioIDs := []pacta.PortfolioID{}
	for _, pm := range portfolioMembers {
		portfolioIDs = append(portfolioIDs, pm.Portfolio.ID)
	}
	portfolios, err := s.DB.Portfolios(s.DB.NoTxn(ctx), portfolioIDs)
	if err != nil {
		return nil, oapierr.Internal("failed to load portfolios for initiative", zap.String("initiative_id", string(id)), zap.Error(err))
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
			SecondaryTargetID:    string(id),
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

type initiativeAuthzVisibilityInfo struct {
	CanManageUsersAndPortfolios bool
	CanSeeInternalInfo          bool
}

func (s *Server) initiativeDoAuthzAndAuditLog(ctx context.Context, iID pacta.InitiativeID, action pacta.AuditLogAction) (*initiativeAuthzVisibilityInfo, error) {
	actorInfo, err := s.getActorInfoOrAnon(ctx)
	if err != nil {
		return nil, err
	}
	userIsInitiativeManager := false
	userIsInitiativeMember := false
	if actorInfo.UserID != "" {
		iur, err := s.DB.InitiativeUserRelationship(s.DB.NoTxn(ctx), iID, actorInfo.UserID)
		if err != nil {
			if !db.IsNotFound(err) {
				return nil, oapierr.Internal("failed to find initiative user relationships", zap.Error(err))
			}
		} else {
			userIsInitiativeMember = iur.Member
			userIsInitiativeManager = iur.Manager
		}
	}
	as := &authzStatus{
		primaryTargetID:      string(iID),
		primaryTargetType:    pacta.AuditLogTargetType_Initiative,
		primaryTargetOwnerID: systemOwnedEntityOwner,
		actorInfo:            actorInfo,
		action:               action,
	}
	switch action {
	case pacta.AuditLogAction_ReadMetadata:
		as.isAuthorized = true
		if userIsInitiativeManager {
			as.authorizedAsActorType = ptr(pacta.AuditLogActorType_Owner)
		} else if actorInfo.IsAdmin || actorInfo.IsSuperAdmin {
			as.authorizedAsActorType = ptr(pacta.AuditLogActorType_Admin)
		} else {
			as.authorizedAsActorType = ptr(pacta.AuditLogActorType_Public)
		}
	case pacta.AuditLogAction_Delete, pacta.AuditLogAction_Create, pacta.AuditLogAction_Update, pacta.AuditLogAction_Download:
		if userIsInitiativeManager {
			as.authorizedAsActorType = ptr(pacta.AuditLogActorType_Owner)
			as.isAuthorized = true
		} else {
			as.isAuthorized, as.authorizedAsActorType = allowIfAdmin(actorInfo)
		}
	default:
		return nil, fmt.Errorf("unknown action %q for initiative authz", action)
	}
	if err := s.auditLogIfAuthorizedOrFail(ctx, as); err != nil {
		return nil, err
	}
	return &initiativeAuthzVisibilityInfo{
		CanSeeInternalInfo:          userIsInitiativeMember || userIsInitiativeManager || actorInfo.IsAdmin || actorInfo.IsSuperAdmin,
		CanManageUsersAndPortfolios: userIsInitiativeManager || actorInfo.IsAdmin || actorInfo.IsSuperAdmin,
	}, nil
}

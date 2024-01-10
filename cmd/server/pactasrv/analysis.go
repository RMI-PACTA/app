package pactasrv

import (
	"context"
	"errors"
	"fmt"

	"github.com/RMI/pacta/cmd/server/pactasrv/conv"
	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/oapierr"
	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
	"github.com/RMI/pacta/task"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// (GET /analyses)
func (s *Server) ListAnalyses(ctx context.Context, request api.ListAnalysesRequestObject) (api.ListAnalysesResponseObject, error) {
	ownerID, err := s.getUserOwnerID(ctx)
	if err != nil {
		return nil, err
	}
	as, err := s.DB.AnalysesByOwner(s.DB.NoTxn(ctx), ownerID)
	if err != nil {
		return nil, oapierr.Internal("failed to query analyses", zap.Error(err))
	}
	items, err := dereference(conv.AnalysesToOAPI(as))
	if err != nil {
		return nil, err
	}
	return api.ListAnalyses200JSONResponse{Items: items}, nil
}

// Deletes an analysis (and its artifacts) by ID
// (DELETE /analysis/{id})
func (s *Server) DeleteAnalysis(ctx context.Context, request api.DeleteAnalysisRequestObject) (api.DeleteAnalysisResponseObject, error) {
	id := pacta.AnalysisID(request.Id)
	if err := s.analysisDoAuthzAndAuditLog(ctx, id, pacta.AuditLogAction_Delete); err != nil {
		return nil, err
	}
	blobURIs, err := s.DB.DeleteAnalysis(s.DB.NoTxn(ctx), id)
	if err != nil {
		return nil, oapierr.Internal("failed to delete analysis", zap.Error(err))
	}
	if err := s.deleteBlobs(ctx, blobURIs...); err != nil {
		return nil, err
	}
	return api.DeleteAnalysis204Response{}, nil
}

// Returns an analysis by ID
// (GET /analysis/{id})
func (s *Server) FindAnalysisById(ctx context.Context, request api.FindAnalysisByIdRequestObject) (api.FindAnalysisByIdResponseObject, error) {
	id := pacta.AnalysisID(request.Id)
	if err := s.analysisDoAuthzAndAuditLog(ctx, id, pacta.AuditLogAction_ReadMetadata); err != nil {
		return nil, err
	}
	a, err := s.DB.Analysis(s.DB.NoTxn(ctx), id)
	if err != nil {
		return nil, oapierr.Internal("failed to query analysis", zap.Error(err))
	}
	if err := s.populateArtifactsInAnalyses(ctx, a); err != nil {
		return nil, err
	}
	if err := s.populateBlobsInAnalysisArtifacts(ctx, a.Artifacts...); err != nil {
		return nil, err
	}
	converted, err := conv.AnalysisToOAPI(a)
	if err != nil {
		return nil, err
	}
	return api.FindAnalysisById200JSONResponse(*converted), nil
}

// Updates writable analysis properties
// (PATCH /analysis/{id})
func (s *Server) UpdateAnalysis(ctx context.Context, request api.UpdateAnalysisRequestObject) (api.UpdateAnalysisResponseObject, error) {
	id := pacta.AnalysisID(request.Id)
	if err := s.analysisDoAuthzAndAuditLog(ctx, id, pacta.AuditLogAction_Update); err != nil {
		return nil, err
	}
	mutations := []db.UpdateAnalysisFn{}
	if request.Body.Name != nil {
		mutations = append(mutations, db.SetAnalysisName(*request.Body.Name))
	}
	if request.Body.Description != nil {
		mutations = append(mutations, db.SetAnalysisDescription(*request.Body.Description))
	}
	if err := s.DB.UpdateAnalysis(s.DB.NoTxn(ctx), id, mutations...); err != nil {
		return nil, oapierr.Internal("failed to update analysis", zap.Error(err))
	}
	return api.UpdateAnalysis204Response{}, nil
}

// Deletes an analysis artifact by ID
// (DELETE /analysis-artifact/{id})
func (s *Server) DeleteAnalysisArtifact(ctx context.Context, request api.DeleteAnalysisArtifactRequestObject) (api.DeleteAnalysisArtifactResponseObject, error) {
	id := pacta.AnalysisArtifactID(request.Id)
	if err := s.analysisArtifactDoAuthzAndAuditLog(ctx, id, pacta.AuditLogAction_Delete); err != nil {
		return nil, err
	}
	blobURI, err := s.DB.DeleteAnalysisArtifact(s.DB.NoTxn(ctx), id)
	if err != nil {
		return nil, oapierr.Internal("failed to delete analysis artifact", zap.Error(err))
	}
	if err := s.deleteBlobs(ctx, blobURI); err != nil {
		return nil, err
	}
	return api.DeleteAnalysisArtifact204Response{}, nil
}

// Updates writable analysis artifact properties
// (PATCH /analysis-artifact/{id})
func (s *Server) UpdateAnalysisArtifact(ctx context.Context, request api.UpdateAnalysisArtifactRequestObject) (api.UpdateAnalysisArtifactResponseObject, error) {
	id := pacta.AnalysisArtifactID(request.Id)
	if err := s.analysisArtifactDoAuthzAndAuditLog(ctx, id, pacta.AuditLogAction_Update); err != nil {
		return nil, err
	}
	mutations := []db.UpdateAnalysisArtifactFn{}
	if request.Body.AdminDebugEnabled != nil {
		mutations = append(mutations, db.SetAnalysisArtifactAdminDebugEnabled(*request.Body.AdminDebugEnabled))
	}
	if request.Body.SharedToPublic != nil {
		mutations = append(mutations, db.SetAnalysisArtifactSharedToPublic(*request.Body.SharedToPublic))
	}
	err := s.DB.UpdateAnalysisArtifact(s.DB.NoTxn(ctx), id, mutations...)
	if err != nil {
		return nil, oapierr.Internal("failed to update analysis artifact", zap.Error(err))
	}
	return api.UpdateAnalysisArtifact204Response{}, nil
}

// Requests an analysis be run
// (POST /run-analysis)
func (s *Server) RunAnalysis(ctx context.Context, request api.RunAnalysisRequestObject) (api.RunAnalysisResponseObject, error) {
	actorInfo, err := s.getActorInfoOrErrIfAnon(ctx)
	if err != nil {
		return nil, err
	}

	analysisType, err := conv.AnalysisTypeFromOAPI(&request.Body.AnalysisType)
	if err != nil {
		return nil, err
	}

	var ais []entityForAnalysis
	if request.Body.InitiativeId != nil {
		ais = append(ais, initiativeAnalysis{iID: pacta.InitiativeID(*request.Body.InitiativeId), s: s})
	}
	if request.Body.PortfolioGroupId != nil {
		ais = append(ais, portfolioGroupAnalysis{pgID: pacta.PortfolioGroupID(*request.Body.PortfolioGroupId), s: s})
	}
	if request.Body.PortfolioId != nil {
		ais = append(ais, portfolioAnalysis{pID: pacta.PortfolioID(*request.Body.PortfolioId), s: s})
	}
	if len(ais) == 0 {
		return nil, oapierr.BadRequest("one of initiative_id, portfolio_group_id, or portfolio_id is required")
	}
	if len(ais) > 1 {
		return nil, oapierr.BadRequest("only one of initiative_id, portfolio_group_id, or portfolio_id may be set")
	}
	ai := ais[0]

	var analysisID pacta.AnalysisID
	var blobURIs []pacta.BlobURI
	err = s.DB.Transactional(ctx, func(tx db.Tx) error {
		var pvID pacta.PACTAVersionID
		if request.Body.PactaVersionId == nil {
			pv, err := s.DB.DefaultPACTAVersion(tx)
			if err != nil {
				return fmt.Errorf("looking up default pacta version: %w", err)
			}
			pvID = pv.ID
		} else {
			pvID = pacta.PACTAVersionID(*request.Body.PactaVersionId)
			_, err := s.DB.PACTAVersion(tx, pvID)
			if err != nil {
				return oapierr.BadRequest("pacta_version_id is invalid", zap.Error(err), zap.String("pacta_version_id", string(pvID)))
			}
		}

		var snapshotID pacta.PortfolioSnapshotID
		var blobIDs []pacta.BlobID
		if err := ai.checkAuth(ctx, tx); err != nil {
			return err
		}

		snapshotID, blobIDs, err := ai.createSnapshot(tx)
		if err != nil {
			return err
		}

		blobs, err := s.DB.Blobs(tx, blobIDs)
		if err != nil {
			return fmt.Errorf("looking up blobs: %w", err)
		}
		for _, blob := range blobs {
			blobURIs = append(blobURIs, blob.BlobURI)
		}

		aID, err := s.DB.CreateAnalysis(tx, &pacta.Analysis{
			AnalysisType:      *analysisType,
			PortfolioSnapshot: &pacta.PortfolioSnapshot{ID: snapshotID},
			PACTAVersion:      &pacta.PACTAVersion{ID: pvID},
			Owner:             &pacta.Owner{ID: actorInfo.OwnerID},
			Name:              request.Body.Name,
			Description:       request.Body.Description,
		})
		if err != nil {
			return fmt.Errorf("creating analysis: %w", err)
		}
		if _, err := s.DB.CreateAuditLog(tx, &pacta.AuditLog{
			ActorType:          pacta.AuditLogActorType_Owner,
			ActorID:            string(actorInfo.UserID),
			ActorOwner:         &pacta.Owner{ID: actorInfo.OwnerID},
			Action:             pacta.AuditLogAction_Create,
			PrimaryTargetType:  pacta.AuditLogTargetType_Analysis,
			PrimaryTargetID:    string(analysisID),
			PrimaryTargetOwner: &pacta.Owner{ID: actorInfo.OwnerID},
		}); err != nil {
			return fmt.Errorf("creating audit log: %w", err)
		}
		analysisID = aID
		return nil
	})
	if err != nil {
		e := &oapierr.Error{}
		if errors.As(err, &e) {
			return nil, e
		}
		return nil, oapierr.Internal("failed to create analysis", zap.Error(err))
	}

	switch *analysisType {
	case pacta.AnalysisType_Audit:
		taskID, runnerID, err := s.TaskRunner.CreateAudit(ctx, &task.CreateAuditRequest{
			AnalysisID: analysisID,
			BlobURIs:   blobURIs,
		})
		if err != nil {
			return nil, oapierr.Internal("failed to create audit task", zap.Error(err))
		}
		s.Logger.Info("created audit task", zap.String("task_id", string(taskID)), zap.String("runner_id", string(runnerID)), zap.String("analysis_id", string(analysisID)))
	case pacta.AnalysisType_Report:
		taskID, runnerID, err := s.TaskRunner.CreateReport(ctx, &task.CreateReportRequest{
			AnalysisID: analysisID,
			BlobURIs:   blobURIs,
		})
		if err != nil {
			return nil, oapierr.Internal("failed to create report task", zap.Error(err))
		}
		s.Logger.Info("created report task", zap.String("task_id", string(taskID)), zap.String("runner_id", string(runnerID)), zap.String("analysis_id", string(analysisID)))
	default:
		return nil, oapierr.Internal("unknown analysis type", zap.String("analysis_type", string(*analysisType)))
	}

	return api.RunAnalysis200JSONResponse{AnalysisId: string(analysisID)}, nil
}

// This name is awkward, but it just encapsulates things we can run an analysis
// on that represent one or more underlying portfolios.
type entityForAnalysis interface {
	checkAuth(context.Context, db.Tx) error
	createSnapshot(db.Tx) (pacta.PortfolioSnapshotID, []pacta.BlobID, error)
}

type portfolioAnalysis struct {
	pID pacta.PortfolioID
	s   *Server

	p *pacta.Portfolio
}

// Allows consistent handling between NOT FOUND and UNAUTHORIZED.
func notFoundErr[T ~string](typeName string, id T, fields ...zapcore.Field) error {
	fs := append(fields, zap.String(fmt.Sprintf("%s_id", typeName), string(id)))
	return oapierr.NotFound(fmt.Sprintf("%s not found", typeName), fs...)
}

func (pa portfolioAnalysis) checkAuth(ctx context.Context, tx db.Tx) error {
	actorInfo, err := pa.s.getActorInfoOrErrIfAnon(ctx)
	if err != nil {
		return err
	}
	p, err := pa.s.DB.Portfolio(tx, pa.pID)
	if err != nil {
		if db.IsNotFound(err) {
			return notFoundErr("portfolio", pa.pID, zap.Error(err))
		}
		return fmt.Errorf("looking up portfolio: %w", err)
	}
	if p.Owner.ID != actorInfo.OwnerID {
		return notFoundErr("portfolio", pa.pID,
			zap.Error(fmt.Errorf("portfolio does not belong to user")),
			zap.String("portfolio_owner_id", string(p.Owner.ID)),
			zap.String("actor_owner_id", string(actorInfo.OwnerID)))
	}
	pa.p = p
	return nil
}

func (pa portfolioAnalysis) createSnapshot(tx db.Tx) (pacta.PortfolioSnapshotID, []pacta.BlobID, error) {
	sID, err := pa.s.DB.CreateSnapshotOfPortfolio(tx, pa.pID)
	if err != nil {
		return "", nil, fmt.Errorf("creating snapshot of portfolio: %w", err)
	}
	return sID, []pacta.BlobID{pa.p.Blob.ID}, nil
}

type portfolioGroupAnalysis struct {
	pgID pacta.PortfolioGroupID
	s    *Server

	pg *pacta.PortfolioGroup
}

func (pga portfolioGroupAnalysis) checkAuth(ctx context.Context, tx db.Tx) error {
	actorInfo, err := pga.s.getActorInfoOrErrIfAnon(ctx)
	if err != nil {
		return err
	}
	pg, err := pga.s.DB.PortfolioGroup(tx, pga.pgID)
	if err != nil {
		if db.IsNotFound(err) {
			return notFoundErr("portfolio_group", pga.pgID, zap.Error(err))
		}
		return fmt.Errorf("looking up portfolio_group: %w", err)
	}
	if pg.Owner.ID != actorInfo.OwnerID {
		return notFoundErr("portfolio_group", pga.pgID,
			zap.Error(fmt.Errorf("portfolio group does not belong to user")),
			zap.String("pg_owner_id", string(pg.Owner.ID)),
			zap.String("actor_owner_id", string(actorInfo.OwnerID)))
	}
	pga.pg = pg
	return nil
}

func (pga portfolioGroupAnalysis) createSnapshot(tx db.Tx) (pacta.PortfolioSnapshotID, []pacta.BlobID, error) {
	sID, err := pga.s.DB.CreateSnapshotOfPortfolioGroup(tx, pga.pgID)
	if err != nil {
		return "", nil, fmt.Errorf("creating snapshot of portfolio group: %w", err)
	}
	pids := []pacta.PortfolioID{}
	for _, pm := range pga.pg.PortfolioGroupMemberships {
		pids = append(pids, pm.Portfolio.ID)
	}
	portfolios, err := pga.s.DB.Portfolios(tx, pids)
	if err != nil {
		return "", nil, fmt.Errorf("looking up portfolios: %w", err)
	}
	var blobIDs []pacta.BlobID
	for _, p := range portfolios {
		blobIDs = append(blobIDs, p.Blob.ID)
	}
	return sID, blobIDs, nil
}

type initiativeAnalysis struct {
	iID pacta.InitiativeID
	s   *Server

	i *pacta.Initiative
}

func (ia initiativeAnalysis) checkAuth(ctx context.Context, tx db.Tx) error {
	// This crudely tests whether or not a user is a manager of the initiative.
	if err := ia.s.initiativeDoAuthzAndAuditLog(ctx, ia.iID, pacta.AuditLogAction_Update); err != nil {
		return err
	}
	i, err := ia.s.DB.Initiative(tx, ia.iID)
	if err != nil {
		if db.IsNotFound(err) {
			return notFoundErr("initiative", ia.iID, zap.Error(err))
		}
		return fmt.Errorf("looking up initiative: %w", err)
	}
	ia.i = i
	return nil
}

func (ia initiativeAnalysis) createSnapshot(tx db.Tx) (pacta.PortfolioSnapshotID, []pacta.BlobID, error) {
	sID, err := ia.s.DB.CreateSnapshotOfInitiative(tx, ia.iID)
	if err != nil {
		return "", nil, fmt.Errorf("creating snapshot of initiative: %w", err)
	}
	pids := []pacta.PortfolioID{}
	for _, pm := range ia.i.PortfolioInitiativeMemberships {
		pids = append(pids, pm.Portfolio.ID)
	}
	portfolios, err := ia.s.DB.Portfolios(tx, pids)
	if err != nil {
		return "", nil, fmt.Errorf("looking up portfolios: %w", err)
	}
	var blobIDs []pacta.BlobID
	for _, p := range portfolios {
		blobIDs = append(blobIDs, p.Blob.ID)
	}
	return sID, blobIDs, nil
}

func (s *Server) analysisDoAuthzAndAuditLog(ctx context.Context, analysisID pacta.AnalysisID, action pacta.AuditLogAction) error {
	actorInfo, err := s.getActorInfoOrErrIfAnon(ctx)
	if err != nil {
		return err
	}
	analysis, err := s.DB.Analysis(s.DB.NoTxn(ctx), analysisID)
	if err != nil {
		if db.IsNotFound(err) {
			return notFoundOrUnauthorized(actorInfo, action, pacta.AuditLogTargetType_Analysis, analysisID)
		}
		return oapierr.Internal("querying analysis for authz failed", zap.Error(err))
	}
	as := &authzStatus{
		primaryTargetID:      string(analysisID),
		primaryTargetType:    pacta.AuditLogTargetType_Analysis,
		primaryTargetOwnerID: analysis.Owner.ID,
		actorInfo:            actorInfo,
		action:               action,
	}
	switch action {
	case pacta.AuditLogAction_Update, pacta.AuditLogAction_Delete, pacta.AuditLogAction_ReadMetadata:
		as.isAuthorized, as.authorizedAsActorType = allowIfAdminOrOwner(actorInfo, analysis.Owner.ID)
	default:
		return fmt.Errorf("unknown action %q for analysis authz", action)

	}
	return s.auditLogIfAuthorizedOrFail(ctx, as)
}

func (s *Server) analysisArtifactDoAuthzAndAuditLog(ctx context.Context, aaID pacta.AnalysisArtifactID, action pacta.AuditLogAction) error {
	actorInfo, err := s.getActorInfoOrErrIfAnon(ctx)
	if err != nil {
		return err
	}
	artifact, err := s.DB.AnalysisArtifact(s.DB.NoTxn(ctx), aaID)
	if err != nil {
		if db.IsNotFound(err) {
			return notFoundOrUnauthorized(actorInfo, action, pacta.AuditLogTargetType_AnalysisArtifact, aaID)
		}
		return oapierr.Internal("failed to look up analysis artifact", zap.String("analysis_artifact_id", string(aaID)), zap.Error(err))
	}
	aID := artifact.AnalysisID
	analysis, err := s.DB.Analysis(s.DB.NoTxn(ctx), aID)
	if err != nil {
		if db.IsNotFound(err) {
			return notFoundOrUnauthorized(actorInfo, action, pacta.AuditLogTargetType_AnalysisArtifact, aaID)
		}
		return oapierr.Internal("failed to look up analysis for analysis artifact", zap.String("analysis_id", string(aID)), zap.Error(err))
	}

	as := &authzStatus{
		primaryTargetID:        string(aID),
		primaryTargetType:      pacta.AuditLogTargetType_Analysis,
		primaryTargetOwnerID:   analysis.Owner.ID,
		secondaryTargetID:      string(aaID),
		secondaryTargetType:    pacta.AuditLogTargetType_AnalysisArtifact,
		secondaryTargetOwnerID: analysis.Owner.ID,
		actorInfo:              actorInfo,
		action:                 action,
	}
	switch action {
	case pacta.AuditLogAction_Download:
		if actorInfo.OwnerID == analysis.Owner.ID {
			as.isAuthorized, as.authorizedAsActorType = true, ptr(pacta.AuditLogActorType_Owner)
		} else if artifact.SharedToPublic {
			as.isAuthorized, as.authorizedAsActorType = true, ptr(pacta.AuditLogActorType_Public)
		} else if artifact.AdminDebugEnabled {
			as.isAuthorized, as.authorizedAsActorType = allowIfAdmin(actorInfo)
		} else {
			as.isAuthorized, as.authorizedAsActorType = false, nil
		}
	case pacta.AuditLogAction_EnableAdminDebug,
		pacta.AuditLogAction_DisableAdminDebug,
		pacta.AuditLogAction_EnableSharing,
		pacta.AuditLogAction_DisableSharing:
		as.isAuthorized, as.authorizedAsActorType = allowIfOwner(actorInfo, analysis.Owner.ID)
	case pacta.AuditLogAction_Update,
		pacta.AuditLogAction_Delete:
		as.isAuthorized, as.authorizedAsActorType = allowIfAdminOrOwner(actorInfo, analysis.Owner.ID)
	default:
		return fmt.Errorf("unknown action %q for analysis_artifact authz", action)
	}
	return s.auditLogIfAuthorizedOrFail(ctx, as)
}

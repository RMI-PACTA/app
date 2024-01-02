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
	_, err := s.checkAnalysisAuthorization(ctx, id)
	if err != nil {
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
	a, err := s.checkAnalysisAuthorization(ctx, pacta.AnalysisID(request.Id))
	if err != nil {
		return nil, err
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
	_, err := s.checkAnalysisAuthorization(ctx, id)
	if err != nil {
		return nil, err
	}
	mutations := []db.UpdateAnalysisFn{}
	if request.Body.Name != nil {
		mutations = append(mutations, db.SetAnalysisName(*request.Body.Name))
	}
	if request.Body.Description != nil {
		mutations = append(mutations, db.SetAnalysisDescription(*request.Body.Description))
	}
	err = s.DB.UpdateAnalysis(s.DB.NoTxn(ctx), id, mutations...)
	if err != nil {
		return nil, oapierr.Internal("failed to update analysis", zap.Error(err))
	}
	return api.UpdateAnalysis204Response{}, nil
}

func (s *Server) checkAnalysisAuthorization(ctx context.Context, id pacta.AnalysisID) (*pacta.Analysis, error) {
	actorOwnerID, err := s.getUserOwnerID(ctx)
	if err != nil {
		return nil, err
	}
	// Extracted to a common variable so that we return the same response for not found and unauthorized.
	notFoundErr := func(fields ...zapcore.Field) error {
		fs := append(fields, zap.String("analysis_id", string(id)))
		return oapierr.NotFound("analysis not found", fs...)
	}
	a, err := s.DB.Analysis(s.DB.NoTxn(ctx), id)
	if err != nil {
		if db.IsNotFound(err) {
			return nil, notFoundErr(zap.Error(err))
		}
		return nil, oapierr.Internal(
			"failed to look up analysis",
			zap.Error(err))
	}
	if a.Owner.ID != actorOwnerID {
		return nil, notFoundErr(
			zap.Error(fmt.Errorf("analysis does not belong to user")),
			zap.String("owner_id", string(a.Owner.ID)),
			zap.String("actor_id", string(actorOwnerID)))
	}
	return a, nil
}

// Deletes an analysis artifact by ID
// (DELETE /analysis-artifact/{id})
func (s *Server) DeleteAnalysisArtifact(ctx context.Context, request api.DeleteAnalysisArtifactRequestObject) (api.DeleteAnalysisArtifactResponseObject, error) {
	id := pacta.AnalysisArtifactID(request.Id)
	err := s.checkAnalysisArtifactAuthorization(ctx, id)
	if err != nil {
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
	err := s.checkAnalysisArtifactAuthorization(ctx, id)
	if err != nil {
		return nil, err
	}
	mutations := []db.UpdateAnalysisArtifactFn{}
	if request.Body.AdminDebugEnabled != nil {
		mutations = append(mutations, db.SetAnalysisArtifactAdminDebugEnabled(*request.Body.AdminDebugEnabled))
	}
	if request.Body.SharedToPublic != nil {
		mutations = append(mutations, db.SetAnalysisArtifactSharedToPublic(*request.Body.SharedToPublic))
	}
	err = s.DB.UpdateAnalysisArtifact(s.DB.NoTxn(ctx), id, mutations...)
	if err != nil {
		return nil, oapierr.Internal("failed to update analysis artifact", zap.Error(err))
	}
	return api.UpdateAnalysisArtifact204Response{}, nil
}

func (s *Server) checkAnalysisArtifactAuthorization(ctx context.Context, id pacta.AnalysisArtifactID) error {
	actorOwnerID, err := s.getUserOwnerID(ctx)
	if err != nil {
		return err
	}
	// Extracted to a common variable so that we return the same response for not found and unauthorized.
	notFoundErr := func(fields ...zapcore.Field) error {
		fs := append(fields, zap.String("analysis_artifact_id", string(id)))
		return oapierr.NotFound("analysis artifact not found", fs...)
	}
	aa, err := s.DB.AnalysisArtifact(s.DB.NoTxn(ctx), id)
	if err != nil {
		if db.IsNotFound(err) {
			return notFoundErr(zap.Error(err))
		}
		return oapierr.Internal("failed to look up analysis artifact", zap.String("analysis_artifact_id", string(id)), zap.Error(err))
	}
	a, err := s.DB.Analysis(s.DB.NoTxn(ctx), aa.AnalysisID)
	if err != nil {
		if db.IsNotFound(err) {
			return notFoundErr(zap.Error(err))
		}
		return oapierr.Internal("failed to look up analysis for analysis artifact",
			zap.String("analysis_id", string(aa.AnalysisID)),
			zap.Error(err))
	}
	if a.Owner.ID != actorOwnerID {
		return notFoundErr(
			zap.Error(fmt.Errorf("analysis artifact does not belong to user")),
			zap.String("owner_id", string(a.Owner.ID)),
			zap.String("actor_id", string(actorOwnerID)),
			zap.String("analysis_id", string(aa.AnalysisID)))
	}
	return nil
}

// Requests an anslysis be run
// (POST /run-analysis)
func (s *Server) RunAnalysis(ctx context.Context, request api.RunAnalysisRequestObject) (api.RunAnalysisResponseObject, error) {
	ownerID, err := s.getUserOwnerID(ctx)
	if err != nil {
		return nil, err
	}
	analysisType, err := conv.AnalysisTypeFromOAPI(&request.Body.AnalysisType)
	if err != nil {
		return nil, err
	}

	found := 0
	var iID pacta.InitiativeID
	var pgID pacta.PortfolioGroupID
	var pID pacta.PortfolioID
	if request.Body.InitiativeId != nil {
		iID = pacta.InitiativeID(*request.Body.InitiativeId)
		found++
	}
	if request.Body.PortfolioGroupId != nil {
		pgID = pacta.PortfolioGroupID(*request.Body.PortfolioGroupId)
		found++
	}
	if request.Body.PortfolioId != nil {
		pID = pacta.PortfolioID(*request.Body.PortfolioId)
		found++
	}
	if found == 0 {
		return nil, oapierr.BadRequest("one of initiative_id, portfolio_group_id, or portfolio_id is required")
	}
	if found > 1 {
		return nil, oapierr.BadRequest("only one of initiative_id, portfolio_group_id, or portfolio_id may be set")
	}

	// Allows consistent handling between NOT FOUND and UNAUTHORIZED.
	notFoundErr := func(typeName string, id string, fields ...zapcore.Field) error {
		fs := append(fields, zap.String(fmt.Sprintf("%s_id", typeName), string(id)))
		return oapierr.NotFound(fmt.Sprintf("%s not found", typeName), fs...)
	}
	var result pacta.AnalysisID
	var endUserErr error
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
				endUserErr = oapierr.BadRequest("pacta_version_id is invalid", zap.Error(err), zap.String("pacta_version_id", string(pvID)))
				return nil
			}
		}

		var snapshotID pacta.PortfolioSnapshotID
		if pID != "" {
			p, err := s.DB.Portfolio(tx, pID)
			if err != nil {
				if db.IsNotFound(err) {
					endUserErr = notFoundErr("portfolio", string(pID), zap.Error(err))
					return nil
				}
				return fmt.Errorf("looking up portfolio: %w", err)
			}
			if p.Owner.ID != ownerID {
				endUserErr = notFoundErr("portfolio", string(pID),
					zap.Error(fmt.Errorf("portfolio does not belong to user")),
					zap.String("portfolio_owner_id", string(p.Owner.ID)),
					zap.String("actor_owner_id", string(ownerID)))
				return nil
			}
			sID, err := s.DB.CreateSnapshotOfPortfolio(tx, pID)
			if err != nil {
				return fmt.Errorf("creating snapshot of portfolio: %w", err)
			}
			snapshotID = sID
		} else if pgID != "" {
			pg, err := s.DB.PortfolioGroup(tx, pgID)
			if err != nil {
				if db.IsNotFound(err) {
					endUserErr = notFoundErr("portfolio_group", string(pgID), zap.Error(err))
					return nil
				}
				return fmt.Errorf("looking up portfolio_group: %w", err)
			}
			if pg.Owner.ID != ownerID {
				endUserErr = notFoundErr("portfolio_group", string(pgID),
					zap.Error(fmt.Errorf("portfolio group does not belong to user")),
					zap.String("pg_owner_id", string(pg.Owner.ID)),
					zap.String("actor_owner_id", string(ownerID)))
				return nil
			}
			sID, err := s.DB.CreateSnapshotOfPortfolioGroup(tx, pgID)
			if err != nil {
				return fmt.Errorf("creating snapshot of portfolio group: %w", err)
			}
			snapshotID = sID
		} else if iID != "" {
			_, err := s.DB.Initiative(tx, iID)
			if err != nil {
				if db.IsNotFound(err) {
					endUserErr = notFoundErr("initiative", string(iID), zap.Error(err))
					return nil
				}
				return fmt.Errorf("looking up initiative: %w", err)
			}
			// TODO(#12) Implement Authorization Here
			sID, err := s.DB.CreateSnapshotOfInitiative(tx, iID)
			if err != nil {
				return fmt.Errorf("creating snapshot of initiative: %w", err)
			}
			snapshotID = sID
		}
		if snapshotID == "" {
			return fmt.Errorf("snapshot id is empty, something is wrong in the bizlogic")
		}

		analysisID, err := s.DB.CreateAnalysis(tx, &pacta.Analysis{
			AnalysisType:      *analysisType,
			PortfolioSnapshot: &pacta.PortfolioSnapshot{ID: snapshotID},
			PACTAVersion:      &pacta.PACTAVersion{ID: pvID},
			Owner:             &pacta.Owner{ID: ownerID},
			Name:              request.Body.Name,
			Description:       request.Body.Description,
		})
		if err != nil {
			return fmt.Errorf("creating analysis: %w", err)
		}
		result = analysisID
		return nil
	})
	if endUserErr != nil {
		return nil, endUserErr
	}
	if err != nil {
		return nil, oapierr.Internal("failed to create analysis", zap.Error(err))
	}

	// TODO - here this is where we'd kick off the analysis run.

	return api.RunAnalysis200JSONResponse{AnalysisId: string(result)}, nil
}

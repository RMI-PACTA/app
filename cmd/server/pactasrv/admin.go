package pactasrv

import (
	"context"
	"fmt"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/oapierr"
	api "github.com/RMI/pacta/openapi/pacta"
	"github.com/RMI/pacta/pacta"
	"go.uber.org/zap"
)

// Merges two users together
// (POST /admin/merge-users)
func (s *Server) MergeUsers(ctx context.Context, request api.MergeUsersRequestObject) (api.MergeUsersResponseObject, error) {
	req := request.Body
	actorUserInfo, err := s.getActorInfoOrErrIfAnon(ctx)
	if err != nil {
		return nil, err
	}
	fieldsIfErr := []zap.Field{
		zap.String("actor_user_id", string(actorUserInfo.UserID)),
		zap.String("from_user_id", req.FromUserId),
		zap.String("to_user_id", req.ToUserId),
	}
	if !actorUserInfo.IsAdmin && !actorUserInfo.IsSuperAdmin {
		return nil, oapierr.Forbidden("only admins can merge users", fieldsIfErr...)
	}

	sourceUID := pacta.UserID(req.FromUserId)
	destUID := pacta.UserID(req.ToUserId)
	if sourceUID == destUID {
		return nil, oapierr.BadRequest("cannot merge user into themselves", fieldsIfErr...)
	}

	var (
		numIncompleteUploads, numAnalyses, numPortfolios, numPortfolioGroups, numAuditLogsCreated int
		buris                                                                                     []pacta.BlobURI
	)

	err = s.DB.Transactional(ctx, func(tx db.Tx) error {
		sourceOwner, err := s.DB.GetOwnerForUser(tx, sourceUID)
		if err != nil {
			return fmt.Errorf("failed to get owner for source user: %w", err)
		}
		destOwner, err := s.DB.GetOwnerForUser(tx, destUID)
		if err != nil {
			return fmt.Errorf("failed to get owner for destination user: %w", err)
		}

		if err = s.DB.RecordUserMerge(tx, sourceUID, destUID, actorUserInfo.UserID); err != nil {
			return fmt.Errorf("failed to record user merge: %w", err)
		}
		if err = s.DB.RecordOwnerMerge(tx, sourceOwner, destOwner, actorUserInfo.UserID); err != nil {
			return fmt.Errorf("failed to record owner merge: %w", err)
		}

		auditLogsToCreate := []*pacta.AuditLog{}
		addAuditLog := func(t pacta.AuditLogTargetType, id string) {
			auditLogsToCreate = append(auditLogsToCreate, &pacta.AuditLog{
				Action:               pacta.AuditLogAction_TransferOwnership,
				ActorType:            pacta.AuditLogActorType_Admin,
				ActorID:              string(actorUserInfo.UserID),
				ActorOwner:           &pacta.Owner{ID: actorUserInfo.OwnerID},
				PrimaryTargetType:    t,
				PrimaryTargetID:      id,
				PrimaryTargetOwner:   &pacta.Owner{ID: destOwner},
				SecondaryTargetType:  pacta.AuditLogTargetType_User,
				SecondaryTargetID:    string(sourceUID),
				SecondaryTargetOwner: &pacta.Owner{ID: sourceOwner},
			})
		}

		incompleteUploads, err := s.DB.IncompleteUploadsByOwner(tx, sourceOwner)
		if err != nil {
			return fmt.Errorf("failed to get incomplete uploads for source owner: %w", err)
		}
		for i, upload := range incompleteUploads {
			err := s.DB.UpdateIncompleteUpload(tx, upload.ID, db.SetIncompleteUploadOwner(destOwner))
			if err != nil {
				return fmt.Errorf("failed to update upload owner %d/%d: %w", i+1, len(incompleteUploads), err)
			}
			addAuditLog(pacta.AuditLogTargetType_IncompleteUpload, string(upload.ID))
		}
		numIncompleteUploads = len(incompleteUploads)

		analyses, err := s.DB.AnalysesByOwner(tx, sourceOwner)
		if err != nil {
			return fmt.Errorf("failed to get analyses for source owner: %w", err)
		}
		for i, analysis := range analyses {
			err := s.DB.UpdateAnalysis(tx, analysis.ID, db.SetAnalysisOwner(destOwner))
			if err != nil {
				return fmt.Errorf("failed to update analysis owner %d/%d: %w", i+1, len(analyses), err)
			}
			addAuditLog(pacta.AuditLogTargetType_Analysis, string(analysis.ID))
		}
		numAnalyses = len(analyses)

		portfolios, err := s.DB.PortfoliosByOwner(tx, sourceOwner)
		if err != nil {
			return fmt.Errorf("failed to get portfolios for source owner: %w", err)
		}
		for i, portfolio := range portfolios {
			err := s.DB.UpdatePortfolio(tx, portfolio.ID, db.SetPortfolioOwner(destOwner))
			if err != nil {
				return fmt.Errorf("failed to update portfolio owner %d/%d: %w", i+1, len(portfolios), err)
			}
			addAuditLog(pacta.AuditLogTargetType_Portfolio, string(portfolio.ID))
		}
		numPortfolios = len(portfolios)

		portfolioGroups, err := s.DB.PortfolioGroupsByOwner(tx, sourceOwner)
		if err != nil {
			return fmt.Errorf("failed to get portfolio groups for source owner: %w", err)
		}
		for i, portfolioGroup := range portfolioGroups {
			err := s.DB.UpdatePortfolioGroup(tx, portfolioGroup.ID, db.SetPortfolioGroupOwner(destOwner))
			if err != nil {
				return fmt.Errorf("failed to update portfolio group owner %d/%d: %w", i+1, len(portfolioGroups), err)
			}
			addAuditLog(pacta.AuditLogTargetType_PortfolioGroup, string(portfolioGroup.ID))
		}
		numPortfolioGroups = len(portfolioGroups)

		if err := s.DB.CreateAuditLogs(tx, auditLogsToCreate); err != nil {
			return fmt.Errorf("failed to create audit logs: %w", err)
		}
		numAuditLogsCreated = len(auditLogsToCreate)

		// Now that we've transferred all the entities, we can delete the user.
		deletedUserBuris, err := s.DB.DeleteUser(tx, sourceUID)
		if err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}
		if len(deletedUserBuris) > 0 {
			// Note in this case we won't commit the transaction, so this data won't be orphaned.
			return fmt.Errorf("failed to delete user: user still has blobs: %v", deletedUserBuris)
		}

		return nil
	})
	if err != nil {
		fieldsIfErr := append(fieldsIfErr, zap.Error(err))
		return nil, oapierr.Internal("failed to merge users", fieldsIfErr...)
	}

	if err := s.deleteBlobs(ctx, buris...); err != nil {
		return nil, err
	}

	s.Logger.Info("user merge completed successfully",
		zap.String("actor_user_id", string(actorUserInfo.UserID)),
		zap.String("from_user_id", req.FromUserId),
		zap.String("to_user_id", req.ToUserId),
		zap.Int("num_incomplete_uploads", numIncompleteUploads),
		zap.Int("num_analyses", numAnalyses),
		zap.Int("num_portfolios", numPortfolios),
		zap.Int("num_portfolio_groups", numPortfolioGroups),
		zap.Int("num_audit_logs_created", numAuditLogsCreated),
	)
	return api.MergeUsers200JSONResponse{
		AuditLogsCreated:      numAuditLogsCreated,
		IncompleteUploadCount: numIncompleteUploads,
		PortfolioCount:        numPortfolios,
		PortfolioGroupCount:   numPortfolioGroups,
		AnalysisCount:         numAnalyses,
	}, nil
}

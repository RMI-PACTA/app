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
	actorUserInfo, err := s.getActorInfoOrFail(ctx)
	if err != nil {
		return nil, err
	}
	if !actorUserInfo.IsAdmin && !actorUserInfo.IsSuperAdmin {
		return nil, oapierr.Forbidden("only admins can merge users",
			zap.String("actor_user_id", string(actorUserInfo.UserID)),
			zap.String("from_user_id", req.FromUserId),
			zap.String("to_user_id", req.ToUserId))
	}

	sourceUID := pacta.UserID(req.FromUserId)
	destUID := pacta.UserID(req.ToUserId)

	var (
		numIncompleteUploads, numAnalyses, numPortfolios, numPortfolioGroups, numAuditLogs, numAuditLogsCreated int
		buris                                                                                                   []pacta.BlobURI
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

		// Note we do an audit log transfer FIRST so that we don't transfer the audit logs generated from the transfer itself.
		nal, err := s.DB.TransferAuditLogOwnership(tx, sourceUID, destUID, sourceOwner, destOwner)
		if err != nil {
			return fmt.Errorf("failed to transfer audit log ownership: %w", err)
		}
		numAuditLogs = nal

		auditLogsToCreate := []pacta.AuditLog{}
		addAuditLog := func(t pacta.AuditLogTargetType, id string) {
			auditLogsToCreate = append(auditLogsToCreate, pacta.AuditLog{
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
				return fmt.Errorf("failed to update upload owner %d/%d: %w", i, len(incompleteUploads), err)
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
				return fmt.Errorf("failed to update analysis owner %d/%d: %w", i, len(analyses), err)
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
				return fmt.Errorf("failed to update portfolio owner %d/%d: %w", i, len(portfolios), err)
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
				return fmt.Errorf("failed to update portfolio group owner %d/%d: %w", i, len(portfolioGroups), err)
			}
			addAuditLog(pacta.AuditLogTargetType_PortfolioGroup, string(portfolioGroup.ID))
		}
		numPortfolioGroups = len(portfolioGroups)

		for _, auditLog := range auditLogsToCreate {
			_, err := s.DB.CreateAuditLog(tx, &auditLog)
			if err != nil {
				return fmt.Errorf("failed to create audit log: %w", err)
			}
		}
		numAuditLogsCreated = len(auditLogsToCreate)

		// Now that we've transferred all the audit logs, we can delete the user.
		newBuris, err := s.DB.DeleteUser(tx, sourceUID)
		if err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}
		buris = append(buris, newBuris...)

		return nil
	})
	if err != nil {
		return nil, oapierr.Internal("failed to merge users", zap.Error(err), zap.String("actor_user_id", string(actorUserInfo.UserID)), zap.String("from_user_id", req.FromUserId), zap.String("to_user_id", req.ToUserId))
	}

	if err := s.deleteBlobs(ctx, buris...); err != nil {
		return nil, err
	}

	s.Logger.Info("user merge completed successfully",
		zap.String("actor_user_id", string(actorUserInfo.UserID)),
		zap.String("from_user_id", req.FromUserId),
		zap.String("to_user_id", req.ToUserId),
		zap.Int("num_audit_logs_transferred", numAuditLogs),
		zap.Int("num_incomplete_uploads", numIncompleteUploads),
		zap.Int("num_analyses", numAnalyses),
		zap.Int("num_portfolios", numPortfolios),
		zap.Int("num_portfolio_groups", numPortfolioGroups),
		zap.Int("num_audit_logs_created", numAuditLogsCreated),
	)
	return api.MergeUsers204Response{}, nil
}

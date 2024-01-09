package pactasrv

import (
	"context"

	"github.com/RMI/pacta/oapierr"
	"github.com/RMI/pacta/pacta"
	"go.uber.org/zap"
)

type authzStatus struct {
	primaryTargetID      string
	primaryTargetType    pacta.AuditLogTargetType
	primaryTargetOwnerID pacta.OwnerID

	secondaryTargetID      string
	secondaryTargetType    pacta.AuditLogTargetType
	secondaryTargetOwnerID pacta.OwnerID

	actorInfo actorInfo

	action                pacta.AuditLogAction
	isAuthorized          bool
	authorizedAsActorType *pacta.AuditLogActorType
}

func (as *authzStatus) actorUserID() string {
	if as.actorInfo.UserID == "" {
		return "ANONYMOUS"
	}
	return string(as.actorInfo.UserID)
}

func (as *authzStatus) actorOwner() *pacta.Owner {
	if as.actorInfo.OwnerID == "" {
		return &pacta.Owner{ID: "ANONYMOUS"}
	}
	return &pacta.Owner{ID: as.actorInfo.OwnerID}
}

func (as *authzStatus) ToAuditLog() (*pacta.AuditLog, error) {
	fieldsIfErr := []zap.Field{
		zap.String("action", string(as.action)),
		zap.String("target_type", string(as.primaryTargetType)),
		zap.String("target_id", as.primaryTargetID),
	}
	if !as.isAuthorized {
		return nil, oapierr.Internal("cannot create audit log for unauthorized action", fieldsIfErr...)
	}
	if as.authorizedAsActorType == nil {
		return nil, oapierr.Internal("cannot create audit log for an unknown actor role", fieldsIfErr...)
	}
	result := &pacta.AuditLog{
		ActorType:  *as.authorizedAsActorType,
		ActorID:    as.actorUserID(),
		ActorOwner: as.actorOwner(),

		Action: as.action,

		PrimaryTargetType:  as.primaryTargetType,
		PrimaryTargetID:    as.primaryTargetID,
		PrimaryTargetOwner: &pacta.Owner{ID: as.primaryTargetOwnerID},
	}
	if as.secondaryTargetType != "" {
		result.SecondaryTargetType = as.secondaryTargetType
		result.SecondaryTargetID = as.secondaryTargetID
		result.SecondaryTargetOwner = &pacta.Owner{ID: as.secondaryTargetOwnerID}
	}
	return result, nil
}

func notFoundOrUnauthorized[T ~string](actorInfo actorInfo, action pacta.AuditLogAction, primaryTargetType pacta.AuditLogTargetType, primaryTargetID T) error {
	return oapierr.NotFound("not found or unauthorized",
		zap.String("target_type", string(primaryTargetType)),
		zap.String("target_id", string(primaryTargetID)),
		zap.String("action", string(action)),
		zap.String("actor_id", string(actorInfo.UserID)),
		zap.String("actor_owner_id", string(actorInfo.OwnerID)))
}

func (s *Server) auditLogIfAuthorizedOrFail(ctx context.Context, status *authzStatus) error {
	if !status.isAuthorized {
		return notFoundOrUnauthorized(status.actorInfo, status.action, status.primaryTargetType, status.primaryTargetID)
	}
	al, err := status.ToAuditLog()
	if err != nil {
		return err
	}
	_, err = s.DB.CreateAuditLog(s.DB.NoTxn(ctx), al)
	if err != nil {
		return oapierr.Internal("creating audit log failed", zap.Error(err))
	}
	return nil
}

func (s *Server) auditLogForCreateEvent(ctx context.Context, actorInfo actorInfo, actorType pacta.AuditLogActorType, primaryTargetType pacta.AuditLogTargetType, primaryTargetID string) error {
	as := &authzStatus{
		primaryTargetID:       primaryTargetID,
		primaryTargetType:     primaryTargetType,
		primaryTargetOwnerID:  actorInfo.OwnerID,
		actorInfo:             actorInfo,
		action:                pacta.AuditLogAction_Create,
		isAuthorized:          true,
		authorizedAsActorType: &actorType,
	}
	return s.auditLogIfAuthorizedOrFail(ctx, as)
}

func allowIfAdmin(actorInfo actorInfo) (bool, *pacta.AuditLogActorType) {
	if actorInfo.IsAdmin {
		return true, ptr(pacta.AuditLogActorType_Admin)
	}
	if actorInfo.IsSuperAdmin {
		return true, ptr(pacta.AuditLogActorType_SuperAdmin)
	}
	return false, nil
}

func allowIfOwner(actorInfo actorInfo, targetOwnerID pacta.OwnerID) (bool, *pacta.AuditLogActorType) {
	if actorInfo.OwnerID == targetOwnerID {
		return true, ptr(pacta.AuditLogActorType_Owner)
	}
	return false, nil
}

func allowIfAdminOrOwner(actorInfo actorInfo, targetOwnerID pacta.OwnerID) (bool, *pacta.AuditLogActorType) {
	if actorInfo.OwnerID == targetOwnerID {
		return true, ptr(pacta.AuditLogActorType_Owner)
	}
	return allowIfAdmin(actorInfo)
}

const systemOwnedEntityOwner = "SYSTEM-OWNED"

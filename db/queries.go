package db

import (
	"fmt"
	"time"

	"github.com/RMI/pacta/pacta"
)

type Cursor string

type PageInfo struct {
	Cursor      Cursor
	HasNextPage bool
}

type AuditLogQuerySortBy string

const (
	AuditLogQuerySortBy_CreatedAt              AuditLogQuerySortBy = "created_at"
	AuditLogQuerySortBy_ActorType              AuditLogQuerySortBy = "actor_type"
	AuditLogQuerySortBy_ActorID                AuditLogQuerySortBy = "actor_id"
	AuditLogQuerySortBy_ActorOwnerID           AuditLogQuerySortBy = "actor_owner_id"
	AuditLogQuerySortBy_PrimaryTargetID        AuditLogQuerySortBy = "primary_target_id"
	AuditLogQuerySortBy_PrimaryTargetType      AuditLogQuerySortBy = "primary_target_type"
	AuditLogQuerySortBy_PrimaryTargetOwnerID   AuditLogQuerySortBy = "primary_target_owner_id"
	AuditLogQuerySortBy_SecondaryTargetID      AuditLogQuerySortBy = "secondary_target_id"
	AuditLogQuerySortBy_SecondaryTargetType    AuditLogQuerySortBy = "secondary_target_type"
	AuditLogQuerySortBy_SecondaryTargetOwnerID AuditLogQuerySortBy = "secondary_target_owner_id"
)

func ParseAuditLogQuerySortBy(s string) (AuditLogQuerySortBy, error) {
	switch s {
	case "created_at":
		return AuditLogQuerySortBy_CreatedAt, nil
	case "actor_type":
		return AuditLogQuerySortBy_ActorType, nil
	case "actor_id":
		return AuditLogQuerySortBy_ActorID, nil
	case "actor_owner_id":
		return AuditLogQuerySortBy_ActorOwnerID, nil
	case "primary_target_id":
		return AuditLogQuerySortBy_PrimaryTargetID, nil
	case "primary_target_type":
		return AuditLogQuerySortBy_PrimaryTargetType, nil
	case "primary_target_owner_id":
		return AuditLogQuerySortBy_PrimaryTargetOwnerID, nil
	case "secondary_target_id":
		return AuditLogQuerySortBy_SecondaryTargetID, nil
	case "secondary_target_type":
		return AuditLogQuerySortBy_SecondaryTargetType, nil
	case "secondary_target_owner_id":
		return AuditLogQuerySortBy_SecondaryTargetOwnerID, nil
	}
	return "", fmt.Errorf("unknown ParseAuditLogActorType: %q", s)
}

type AuditLogQuerySort struct {
	By        AuditLogQuerySortBy
	Ascending bool
}

type AuditLogQueryWhere struct {
	InID            []pacta.AuditLogID
	MinCreatedAt    time.Time
	MaxCreatedAt    time.Time
	InAction        []pacta.AuditLogAction
	InActorType     []pacta.AuditLogActorType
	InActorID       []string
	InActorOwnerID  []pacta.OwnerID
	InTargetType    []pacta.AuditLogTargetType
	InTargetID      []string
	InTargetOwnerID []pacta.OwnerID
}

type AuditLogQuery struct {
	Cursor Cursor
	Limit  int
	Wheres []*AuditLogQueryWhere
	Sorts  []*AuditLogQuerySort
}

type UserQuerySortBy string

const (
	UserQuerySortBy_CreatedAt UserQuerySortBy = "created_at"
)

type UserQuerySort struct {
	By        UserQuerySortBy
	Ascending bool
}

type UserQueryWhere struct {
	NameOrEmailLike string
}

type UserQuery struct {
	Cursor Cursor
	Limit  int
	Wheres []*UserQueryWhere
	Sorts  []*UserQuerySort
}

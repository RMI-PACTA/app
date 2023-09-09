package db

import (
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

type AuditLogQuerySort struct {
	By        AuditLogQuerySortBy
	Ascending bool
}

type AuditLogQueryWhere struct {
	InID            []pacta.AuditLogID
	MinCreatedAt    time.Time
	MaxCreatedAt    time.Time
	InActionType    []pacta.AuditLogAction
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

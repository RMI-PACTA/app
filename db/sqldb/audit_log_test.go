package sqldb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCreateAuditLog(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	cmpOpts := auditLogCmpOpts()
	al := &pacta.AuditLog{
		Action:               pacta.AuditLogAction_AddTo,
		ActorType:            pacta.AuditLogActorType_User,
		ActorID:              "user1",
		ActorOwner:           &pacta.Owner{ID: "owner1"},
		PrimaryTargetType:    pacta.AuditLogTargetType_Portfolio,
		PrimaryTargetID:      "portfolio-1",
		PrimaryTargetOwner:   &pacta.Owner{ID: "user2"},
		SecondaryTargetType:  pacta.AuditLogTargetType_PortfolioGroup,
		SecondaryTargetID:    "portfolio-group-1",
		SecondaryTargetOwner: &pacta.Owner{ID: "user3"},
	}
	id, err := tdb.CreateAuditLog(tx, al)
	if err != nil {
		t.Fatalf("error creating audit log: %v", err)
	}
	al.ID = id
	al.CreatedAt = time.Now().UTC()

	als, pi, err := tdb.AuditLogs(tx, &db.AuditLogQuery{
		Limit: 10,
		Wheres: []*db.AuditLogQueryWhere{{
			InID: []pacta.AuditLogID{id},
		}},
	})
	if err != nil {
		t.Fatalf("error getting audit log: %v", err)
	}

	if diff := cmp.Diff(al, als[0], cmpOpts); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}
	if diff := cmp.Diff(&db.PageInfo{HasNextPage: false, Cursor: "1"}, pi, cmpOpts); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}
}

func TestAuditLogActionConvertability(t *testing.T) {
	testAuditLogEnumConvertability(
		t,
		func(a pacta.AuditLogAction, al *pacta.AuditLog) { al.Action = a },
		func(al *pacta.AuditLog) pacta.AuditLogAction { return al.Action },
		pacta.AuditLogActionValues,
	)
}

func TestAuditLogActorTypeConvertability(t *testing.T) {
	testAuditLogEnumConvertability(
		t,
		func(a pacta.AuditLogActorType, al *pacta.AuditLog) { al.ActorType = a },
		func(al *pacta.AuditLog) pacta.AuditLogActorType { return al.ActorType },
		pacta.AuditLogActorTypeValues,
	)
}

func TestAuditLogTargetTypeConvertability(t *testing.T) {
	testAuditLogEnumConvertability(
		t,
		func(a pacta.AuditLogTargetType, al *pacta.AuditLog) { al.PrimaryTargetType = a },
		func(al *pacta.AuditLog) pacta.AuditLogTargetType { return al.PrimaryTargetType },
		pacta.AuditLogTargetTypeValues,
	)
}

func testAuditLogEnumConvertability[E comparable](t *testing.T, writeE func(E, *pacta.AuditLog), readE func(*pacta.AuditLog) E, values []E) {
	var zeroValue E
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	base := &pacta.AuditLog{
		Action:               pacta.AuditLogAction_AddTo,
		ActorType:            pacta.AuditLogActorType_User,
		ActorID:              "user1",
		ActorOwner:           &pacta.Owner{ID: "owner1"},
		PrimaryTargetType:    pacta.AuditLogTargetType_Portfolio,
		PrimaryTargetID:      "portfolio-1",
		PrimaryTargetOwner:   &pacta.Owner{ID: "user2"},
		SecondaryTargetType:  pacta.AuditLogTargetType_PortfolioGroup,
		SecondaryTargetID:    "portfolio-group-1",
		SecondaryTargetOwner: &pacta.Owner{ID: "user3"},
	}
	var id pacta.AuditLogID

	write := func(e E) error {
		al := base.Clone()
		writeE(e, al)
		id2, err := tdb.CreateAuditLog(tx, al)
		id = id2
		return err
	}
	read := func() (E, error) {
		als, _, err := tdb.AuditLogs(tx, &db.AuditLogQuery{
			Limit: 10,
			Wheres: []*db.AuditLogQueryWhere{{
				InID: []pacta.AuditLogID{id},
			}},
		})
		if err != nil {
			return zeroValue, fmt.Errorf("reading audit logs: %w", err)
		}
		return readE(als[0]), nil
	}

	testEnumConvertability(t, write, read, values)
}

func TestAuditSearch(t *testing.T) {
	beforeCreation := time.Now()
	action1 := pacta.AuditLogAction_AddTo
	action2 := pacta.AuditLogAction_Create
	actorType1 := pacta.AuditLogActorType_User
	actorType2 := pacta.AuditLogActorType_System
	actorID1 := "user1"
	actorID2 := "system2"
	actorOwner1 := &pacta.Owner{ID: "owner1"}
	actorOwner2 := &pacta.Owner{ID: "owner2"}
	targetType1 := pacta.AuditLogTargetType_Portfolio
	targetType2 := pacta.AuditLogTargetType_IncompleteUpload
	targetID1 := "portfolio-1"
	targetID2 := "incomplete-upload-2"
	targetOwner1 := &pacta.Owner{ID: "owner3"}
	targetOwner2 := &pacta.Owner{ID: "owner4"}

	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)

	alID1, err0 := tdb.CreateAuditLog(tx, &pacta.AuditLog{ActorType: actorType1, ActorID: actorID1, ActorOwner: actorOwner1, Action: action1, PrimaryTargetType: targetType1, PrimaryTargetID: targetID1, PrimaryTargetOwner: targetOwner1})
	alID2, err1 := tdb.CreateAuditLog(tx, &pacta.AuditLog{ActorType: actorType2, ActorID: actorID2, ActorOwner: actorOwner2, Action: action2, PrimaryTargetType: targetType2, PrimaryTargetID: targetID2, PrimaryTargetOwner: targetOwner2})
	alID3, err2 := tdb.CreateAuditLog(tx, &pacta.AuditLog{ActorType: actorType2, ActorID: actorID2, ActorOwner: actorOwner2, Action: action2, PrimaryTargetType: targetType2, PrimaryTargetID: "something", PrimaryTargetOwner: targetOwner1, SecondaryTargetType: targetType2, SecondaryTargetID: targetID2, SecondaryTargetOwner: targetOwner2})
	afterCreation := time.Now()
	noErrDuringSetup(t, err0, err1, err2)

	t.Run("Single Where Tests", func(t *testing.T) {
		cases := []struct {
			name     string
			where    *db.AuditLogQueryWhere
			expected []pacta.AuditLogID
		}{{
			name:     "By ID Singular",
			where:    &db.AuditLogQueryWhere{InID: []pacta.AuditLogID{alID2}},
			expected: []pacta.AuditLogID{alID2},
		}, {
			name:     "By ID Multiple",
			where:    &db.AuditLogQueryWhere{InID: []pacta.AuditLogID{alID2, alID3}},
			expected: []pacta.AuditLogID{alID2, alID3},
		}, {
			name:     "By Created At After Creation",
			where:    &db.AuditLogQueryWhere{MinCreatedAt: afterCreation},
			expected: []pacta.AuditLogID{},
		}, {
			name:     "By Created At Before Creation",
			where:    &db.AuditLogQueryWhere{MaxCreatedAt: beforeCreation},
			expected: []pacta.AuditLogID{},
		}, {
			name:     "By Created At Before Creation",
			where:    &db.AuditLogQueryWhere{MaxCreatedAt: afterCreation, MinCreatedAt: beforeCreation},
			expected: []pacta.AuditLogID{alID1, alID2, alID3},
		}, {
			name:     "By ActionType",
			where:    &db.AuditLogQueryWhere{InActionType: []pacta.AuditLogAction{action2}},
			expected: []pacta.AuditLogID{alID2, alID3},
		}, {
			name:     "By ActorType",
			where:    &db.AuditLogQueryWhere{InActorType: []pacta.AuditLogActorType{actorType1}},
			expected: []pacta.AuditLogID{alID1},
		}, {
			name:     "By ActorID",
			where:    &db.AuditLogQueryWhere{InActorID: []string{actorID2}},
			expected: []pacta.AuditLogID{alID2, alID3},
		}, {
			name:     "By TargetType",
			where:    &db.AuditLogQueryWhere{InTargetType: []pacta.AuditLogTargetType{targetType2}},
			expected: []pacta.AuditLogID{alID2, alID3},
		}, {
			name:     "By TargetID",
			where:    &db.AuditLogQueryWhere{InTargetID: []string{targetID2}},
			expected: []pacta.AuditLogID{alID2, alID3},
		}, {
			name:     "By TargetOwnerID Includes Both",
			where:    &db.AuditLogQueryWhere{InTargetOwnerID: []pacta.OwnerID{targetOwner2.ID}},
			expected: []pacta.AuditLogID{alID2, alID3},
		}, {
			name:     "By TargetOwnerID Includes Both Two",
			where:    &db.AuditLogQueryWhere{InTargetOwnerID: []pacta.OwnerID{targetOwner1.ID}},
			expected: []pacta.AuditLogID{alID1, alID3},
		}, {
			name:     "By TargetOwnerID Doesn't Duplicate",
			where:    &db.AuditLogQueryWhere{InTargetOwnerID: []pacta.OwnerID{targetOwner1.ID, targetOwner2.ID}},
			expected: []pacta.AuditLogID{alID1, alID2, alID3},
		}}

		for i, c := range cases {
			t.Run(fmt.Sprintf("case %d: %q", i, c.name), func(t *testing.T) {
				auditLogs, _, err := tdb.AuditLogs(tx, &db.AuditLogQuery{
					Limit:  10,
					Wheres: []*db.AuditLogQueryWhere{c.where},
				})
				if err != nil {
					t.Fatalf("getting audit logs: %v", err)
				}
				actual := make([]pacta.AuditLogID, len(auditLogs))
				for i, a := range auditLogs {
					actual[i] = a.ID
				}
				if diff := cmp.Diff(c.expected, actual, sortAuditLogIDs()); diff != "" {
					t.Errorf("unexpected diff:\n%s", diff)
				}
			})
		}
	})

	t.Run("Multiple Where Tests Are Conjunctive", func(t *testing.T) {
		cases := []struct {
			name     string
			where    []*db.AuditLogQueryWhere
			expected []pacta.AuditLogID
		}{{
			name: "All Match",
			where: []*db.AuditLogQueryWhere{
				&db.AuditLogQueryWhere{InID: []pacta.AuditLogID{alID1}},
				&db.AuditLogQueryWhere{MinCreatedAt: beforeCreation},
				&db.AuditLogQueryWhere{MaxCreatedAt: afterCreation},
				&db.AuditLogQueryWhere{InActionType: []pacta.AuditLogAction{action1}},
				&db.AuditLogQueryWhere{InActorType: []pacta.AuditLogActorType{actorType1}},
				&db.AuditLogQueryWhere{InActorID: []string{actorID1}},
				&db.AuditLogQueryWhere{InActorOwnerID: []pacta.OwnerID{actorOwner1.ID}},
				&db.AuditLogQueryWhere{InTargetType: []pacta.AuditLogTargetType{targetType1}},
				&db.AuditLogQueryWhere{InTargetID: []string{targetID1}},
				&db.AuditLogQueryWhere{InTargetOwnerID: []pacta.OwnerID{targetOwner1.ID}},
			},
			expected: []pacta.AuditLogID{alID1},
		}, {
			name: "One Does not Match",
			where: []*db.AuditLogQueryWhere{
				&db.AuditLogQueryWhere{InID: []pacta.AuditLogID{alID1}},
				&db.AuditLogQueryWhere{MinCreatedAt: beforeCreation},
				&db.AuditLogQueryWhere{MaxCreatedAt: afterCreation},
				&db.AuditLogQueryWhere{InActionType: []pacta.AuditLogAction{action1}},
				&db.AuditLogQueryWhere{InActorType: []pacta.AuditLogActorType{actorType2}},
				&db.AuditLogQueryWhere{InActorID: []string{actorID1}},
				&db.AuditLogQueryWhere{InActorOwnerID: []pacta.OwnerID{actorOwner1.ID}},
				&db.AuditLogQueryWhere{InTargetType: []pacta.AuditLogTargetType{targetType1}},
				&db.AuditLogQueryWhere{InTargetID: []string{targetID1}},
				&db.AuditLogQueryWhere{InTargetOwnerID: []pacta.OwnerID{targetOwner1.ID}},
			},
			expected: []pacta.AuditLogID{},
		}}

		for i, c := range cases {
			t.Run(fmt.Sprintf("case %d: %q", i, c.name), func(t *testing.T) {
				auditLogs, _, err := tdb.AuditLogs(tx, &db.AuditLogQuery{
					Limit:  10,
					Wheres: c.where,
				})
				if err != nil {
					t.Fatalf("getting audit logs: %v", err)
				}
				actual := make([]pacta.AuditLogID, len(auditLogs))
				for i, a := range auditLogs {
					actual[i] = a.ID
				}
				if diff := cmp.Diff(c.expected, actual, sortAuditLogIDs()); diff != "" {
					t.Errorf("unexpected diff:\n%s", diff)
				}
			})
		}
	})
}

func auditLogCmpOpts() cmp.Option {
	return cmp.Options{
		cmpopts.EquateEmpty(),
		cmpopts.EquateApproxTime(time.Second),
	}
}

func sortAuditLogIDs() cmp.Option {
	return cmpopts.SortSlices(func(a, b pacta.AuditLogID) bool {
		return string(a) < string(b)
	})
}

package sqldb

import (
	"context"
	"testing"
	"time"

	"github.com/RMI/pacta/pacta"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestGetOrCreateOwners(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	u := userForTesting(t, tdb)
	i := initiativeForTesting(t, tdb)

	uo, err := tdb.GetOrCreateOwnerForUser(tx, u.ID)
	if err != nil {
		t.Fatalf("creating owner for user: %v", err)
	}
	uo2, err := tdb.GetOrCreateOwnerForUser(tx, u.ID)
	if err != nil {
		t.Fatalf("creating owner for user: %v", err)
	}
	if uo2 != uo {
		t.Fatalf("expected owner id %q, got %q", uo, uo2)
	}

	io, err := tdb.GetOrCreateOwnerForInitiative(tx, i.ID)
	if err != nil {
		t.Fatalf("creating owner for initiative: %v", err)
	}
	io2, err := tdb.GetOrCreateOwnerForInitiative(tx, i.ID)
	if err != nil {
		t.Fatalf("creating owner for initiative: %v", err)
	}
	if io2 != io {
		t.Fatalf("expected owner id %q, got %q", io, io2)
	}
}

func TestReadOwners(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	u1 := userForTestingWithKey(t, tdb, "1")
	u2 := userForTestingWithKey(t, tdb, "2")
	ou1 := ownerUserForTesting(t, tdb, u1)
	ou2 := ownerUserForTesting(t, tdb, u2)
	i1 := initiativeForTesting(t, tdb)
	oi1 := ownerInitiativeForTesting(t, tdb, i1)
	cmpOpts := ownerCmpOpts()

	// Read by id
	actual, err := tdb.Owner(tx, ou1.ID)
	if err != nil {
		t.Fatalf("reading owner: %v", err)
	}
	expected := ou1
	if diff := cmp.Diff(actual, expected, cmpOpts); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}

	// Read by ids
	actualM, err := tdb.Owners(tx, []pacta.OwnerID{ou2.ID, ou1.ID, oi1.ID, oi1.ID})
	if err != nil {
		t.Fatalf("reading owners: %v", err)
	}
	expectedM := map[pacta.OwnerID]*pacta.Owner{
		oi1.ID: oi1,
		ou1.ID: ou1,
		ou2.ID: ou2,
	}
	if diff := cmp.Diff(actualM, expectedM, cmpOpts); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}
}

func ownerUserForTesting(t *testing.T, tdb *DB, u *pacta.User) *pacta.Owner {
	t.Helper()
	tx := tdb.NoTxn(context.Background())
	oid, err := tdb.GetOrCreateOwnerForUser(tx, u.ID)
	if err != nil {
		t.Fatalf("creating owner for user: %v", err)
	}
	o, err := tdb.Owner(tx, oid)
	if err != nil {
		t.Fatalf("getting owner: %v", err)
	}
	return o
}

func ownerInitiativeForTesting(t *testing.T, tdb *DB, i *pacta.Initiative) *pacta.Owner {
	t.Helper()
	tx := tdb.NoTxn(context.Background())
	oid, err := tdb.GetOrCreateOwnerForInitiative(tx, i.ID)
	if err != nil {
		t.Fatalf("creating owner for user: %v", err)
	}
	o, err := tdb.Owner(tx, oid)
	if err != nil {
		t.Fatalf("getting owner: %v", err)
	}
	return o
}

func ownerCmpOpts() cmp.Option {
	ownerIDLessFn := func(a, b pacta.OwnerID) bool {
		return a < b
	}
	ownerLessFn := func(a, b *pacta.Owner) bool {
		return a.ID < b.ID
	}
	return cmp.Options{
		cmpopts.SortSlices(ownerLessFn),
		cmpopts.SortMaps(ownerIDLessFn),
		cmpopts.EquateEmpty(),
		cmpopts.EquateApproxTime(time.Second),
	}
}

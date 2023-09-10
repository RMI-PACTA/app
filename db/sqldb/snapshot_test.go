package sqldb

import (
	"context"
	"testing"
	"time"

	"github.com/RMI/pacta/pacta"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSnapshotCreation(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	u := userForTesting(t, tdb)
	ou := ownerUserForTesting(t, tdb, u)
	p1 := portfolioForTestingWithKey(t, tdb, "1")
	p2 := portfolioForTestingWithKey(t, tdb, "2")
	p3 := portfolioForTestingWithKey(t, tdb, "3")
	p4 := portfolioForTestingWithKey(t, tdb, "4")
	// PG1 = [P1, P2, P3]
	pg1 := portfolioGroupForTesting(t, tdb, ou)
	portfolioGroupMembershipForTesting(t, tdb, pg1, p1)
	portfolioGroupMembershipForTesting(t, tdb, pg1, p2)
	portfolioGroupMembershipForTesting(t, tdb, pg1, p3)
	// PG2 = [P1, P4]
	pg2 := portfolioGroupForTesting(t, tdb, ou)
	portfolioGroupMembershipForTesting(t, tdb, pg2, p1)
	portfolioGroupMembershipForTesting(t, tdb, pg2, p4)
	// PG3 = []
	pg3 := portfolioGroupForTesting(t, tdb, ou)
	// I1 = [P2, P4]
	i1 := initiativeForTestingWithKey(t, tdb, "1")
	portfolioInitiativeForTesting(t, tdb, i1, p2)
	portfolioInitiativeForTesting(t, tdb, i1, p4)
	// I2 = []
	i2 := initiativeForTestingWithKey(t, tdb, "2")

	assertSnapshotContents := func(psID pacta.PortfolioSnapshotID, ps ...*pacta.Portfolio) {
		t.Helper()
		actual, err := tdb.PortfolioSnapshot(tx, psID)
		if err != nil {
			t.Fatalf("reading initiative: %v", err)
		}
		ids := make([]pacta.PortfolioID, len(ps))
		for i, p := range ps {
			ids[i] = p.ID
		}
		if diff := cmp.Diff(ids, actual.PortfolioIDs, portfolioSnapshotCmpOpts()); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	}

	pg1Snap, err := tdb.CreateSnapshotOfPortfolioGroup(tx, pg1.ID)
	if err != nil {
		t.Fatalf("creating snapshot: %v", err)
	}
	assertSnapshotContents(pg1Snap, p1, p2, p3)

	pg2Snap, err := tdb.CreateSnapshotOfPortfolioGroup(tx, pg2.ID)
	if err != nil {
		t.Fatalf("creating snapshot: %v", err)
	}
	assertSnapshotContents(pg2Snap, p1, p4)

	pg3Snap, err := tdb.CreateSnapshotOfPortfolioGroup(tx, pg3.ID)
	if err != nil {
		t.Fatalf("creating snapshot: %v", err)
	}
	assertSnapshotContents(pg3Snap)

	// Changing group membership should allow for a new snapshot to be created.
	portfolioGroupMembershipForTesting(t, tdb, pg1, p4)
	pg1Snap, err = tdb.CreateSnapshotOfPortfolioGroup(tx, pg1.ID)
	if err != nil {
		t.Fatalf("creating snapshot: %v", err)
	}
	assertSnapshotContents(pg1Snap, p1, p2, p3, p4)

	i1Snap, err := tdb.CreateSnapshotOfInitiative(tx, i1.ID)
	if err != nil {
		t.Fatalf("creating snapshot: %v", err)
	}
	assertSnapshotContents(i1Snap, p4, p2)

	i2Snap, err := tdb.CreateSnapshotOfInitiative(tx, i2.ID)
	if err != nil {
		t.Fatalf("creating snapshot: %v", err)
	}
	assertSnapshotContents(i2Snap)

	// Changing group membership should allow for a new snapshot to be created.
	portfolioInitiativeForTesting(t, tdb, i2, p4)
	i2Snap, err = tdb.CreateSnapshotOfInitiative(tx, i2.ID)
	if err != nil {
		t.Fatalf("creating snapshot: %v", err)
	}
	assertSnapshotContents(i2Snap, p4)
}

func portfolioSnapshotCmpOpts() cmp.Option {
	portfolioIDLessFn := func(a, b pacta.PortfolioID) bool {
		return a < b
	}
	return cmp.Options{
		cmpopts.SortSlices(portfolioIDLessFn),
		cmpopts.EquateEmpty(),
		cmpopts.EquateApproxTime(time.Second),
	}
}

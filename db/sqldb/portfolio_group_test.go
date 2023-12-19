package sqldb

import (
	"context"
	"testing"
	"time"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestPortfolioGroupCRUD(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	u1 := userForTestingWithKey(t, tdb, "1")
	u2 := userForTestingWithKey(t, tdb, "2")
	o1 := ownerUserForTesting(t, tdb, u1)
	o2 := ownerUserForTesting(t, tdb, u2)
	p1 := portfolioForTestingWithKey(t, tdb, "3")
	p2 := portfolioForTestingWithKey(t, tdb, "4")

	pg1 := &pacta.PortfolioGroup{
		Name:  "portfolio-group-name",
		Owner: &pacta.Owner{ID: o1.ID},
	}
	pgID1, err := tdb.CreatePortfolioGroup(tx, pg1)
	if err != nil {
		t.Fatalf("creating portfolio group: %w", err)
	}
	pg1.CreatedAt = time.Now()
	pg1.ID = pgID1

	pg2 := &pacta.PortfolioGroup{
		Name:  "portfolio-group-name",
		Owner: &pacta.Owner{ID: o1.ID},
	}
	pgID2, err := tdb.CreatePortfolioGroup(tx, pg2)
	if err != nil {
		t.Fatalf("creating portfolio group: %w", err)
	}
	pg2.CreatedAt = time.Now()
	pg2.ID = pgID2

	nName := "New Portfolio Group Name"
	nDesc := "New Portfolio Group Description"
	err = tdb.UpdatePortfolioGroup(tx, pg1.ID,
		db.SetPortfolioGroupName(nName),
		db.SetPortfolioGroupDescription(nDesc),
		db.SetPortfolioGroupOwner(o2.ID),
	)
	if err != nil {
		t.Fatalf("updating portfolio group: %v", err)
	}
	pg1.Name = nName
	pg1.Description = nDesc
	pg1.Owner = &pacta.Owner{ID: o2.ID}

	actual, err := tdb.PortfolioGroup(tx, pg1.ID)
	if diff := cmp.Diff(pg1, actual, portfolioGroupCmpOpts()); diff != "" {
		t.Fatalf("portfolio group mismatch (-want +got):\n%s", diff)
	}

	actuals, err := tdb.PortfolioGroups(tx, []pacta.PortfolioGroupID{pg1.ID, pg2.ID})
	expecteds := map[pacta.PortfolioGroupID]*pacta.PortfolioGroup{
		pg1.ID: pg1,
		pg2.ID: pg2,
	}
	if diff := cmp.Diff(expecteds, actuals, portfolioGroupCmpOpts()); diff != "" {
		t.Fatalf("portfolio group mismatch (-want +got):\n%s", diff)
	}

	if err := tdb.CreatePortfolioGroupMembership(tx, pg1.ID, p1.ID); err != nil {
		t.Fatalf("creating portfolio group membership: %v", err)
	}
	if err := tdb.CreatePortfolioGroupMembership(tx, pg1.ID, p2.ID); err != nil {
		t.Fatalf("creating portfolio group membership: %v", err)
	}
	if err := tdb.CreatePortfolioGroupMembership(tx, pg2.ID, p1.ID); err != nil {
		t.Fatalf("creating portfolio group membership: %v", err)
	}
	if err := tdb.CreatePortfolioGroupMembership(tx, pg2.ID, p2.ID); err != nil {
		t.Fatalf("creating portfolio group membership: %v", err)
	}
	if err := tdb.DeletePortfolioGroupMembership(tx, pg2.ID, p2.ID); err != nil {
		t.Fatalf("deleting portfolio group membership: %v", err)
	}

	actuals, err = tdb.PortfolioGroups(tx, []pacta.PortfolioGroupID{pg1.ID, pg2.ID})
	pg1.PortfolioGroupMemberships = []*pacta.PortfolioGroupMembership{{
		Portfolio: &pacta.Portfolio{ID: p1.ID},
		CreatedAt: time.Now(),
	}, {
		Portfolio: &pacta.Portfolio{ID: p2.ID},
		CreatedAt: time.Now(),
	}}
	pg2.PortfolioGroupMemberships = []*pacta.PortfolioGroupMembership{{
		Portfolio: &pacta.Portfolio{ID: p1.ID},
		CreatedAt: time.Now(),
	}}
	expecteds = map[pacta.PortfolioGroupID]*pacta.PortfolioGroup{
		pg1.ID: pg1,
		pg2.ID: pg2,
	}
	if diff := cmp.Diff(expecteds, actuals, portfolioGroupCmpOpts()); diff != "" {
		t.Fatalf("portfolio group mismatch (-want +got):\n%s", diff)
	}

	expectedP1 := p1.Clone()
	expectedP1.PortfolioGroupMemberships = []*pacta.PortfolioGroupMembership{{
		PortfolioGroup: &pacta.PortfolioGroup{ID: pg1.ID},
		CreatedAt:      time.Now(),
	}, {
		PortfolioGroup: &pacta.PortfolioGroup{ID: pg2.ID},
		CreatedAt:      time.Now(),
	}}

	err = tdb.DeletePortfolioGroup(tx, pg1.ID)
	if err != nil {
		t.Fatalf("delete portfolio group: %v", err)
	}
	actuals, err = tdb.PortfolioGroups(tx, []pacta.PortfolioGroupID{pg1.ID, pg2.ID})
	expecteds = map[pacta.PortfolioGroupID]*pacta.PortfolioGroup{
		pg2.ID: pg2,
	}
	if diff := cmp.Diff(expecteds, actuals, portfolioGroupCmpOpts()); diff != "" {
		t.Fatalf("portfolio group mismatch (-want +got):\n%s", diff)
	}
}

func TestPortfolioGroupMembership(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	u1 := userForTestingWithKey(t, tdb, "1")
	u2 := userForTestingWithKey(t, tdb, "2")
	o1 := ownerUserForTesting(t, tdb, u1)
	o2 := ownerUserForTesting(t, tdb, u2)
	pg1 := portfolioGroupForTesting(t, tdb, o1)
	pg2 := portfolioGroupForTesting(t, tdb, o2)
	pg3 := portfolioGroupForTesting(t, tdb, o2)
	p1 := portfolioForTestingWithKey(t, tdb, "3")
	p2 := portfolioForTestingWithKey(t, tdb, "4")
	p3 := portfolioForTestingWithKey(t, tdb, "5")

	assertPGMembership := func(pgID pacta.PortfolioGroupID, ps []*pacta.Portfolio) {
		idsOnly := make([]*pacta.Portfolio, len(ps))
		for i, p := range ps {
			idsOnly[i] = &pacta.Portfolio{ID: p.ID}
		}
		t.Helper()
		pg, err := tdb.PortfolioGroup(tx, pgID)
		if err != nil {
			t.Fatalf("getting portfolio group: %v", err)
		}
		ms := make([]*pacta.Portfolio, len(pg.PortfolioGroupMemberships))
		for i, m := range pg.PortfolioGroupMemberships {
			ms[i] = m.Portfolio
		}
		if diff := cmp.Diff(idsOnly, ms, portfolioCmpOpts()); diff != "" {
			t.Fatalf("portfolio group mismatch (-want +got):\n%s", diff)
		}
	}
	createMembership := func(pgID pacta.PortfolioGroupID, pID pacta.PortfolioID) {
		t.Helper()
		if err := tdb.CreatePortfolioGroupMembership(tx, pgID, pID); err != nil {
			t.Fatalf("creating portfolio group membership: %v", err)
		}
	}
	deleteMembership := func(pgID pacta.PortfolioGroupID, pID pacta.PortfolioID) {
		t.Helper()
		if err := tdb.DeletePortfolioGroupMembership(tx, pgID, pID); err != nil {
			t.Fatalf("creating portfolio group membership: %v", err)
		}
	}
	assertPGMembership(pg1.ID, []*pacta.Portfolio{})

	createMembership(pg1.ID, p1.ID)
	assertPGMembership(pg1.ID, []*pacta.Portfolio{p1})

	createMembership(pg1.ID, p2.ID)
	assertPGMembership(pg1.ID, []*pacta.Portfolio{p1, p2})

	createMembership(pg2.ID, p2.ID)
	createMembership(pg2.ID, p2.ID)
	createMembership(pg2.ID, p2.ID)
	assertPGMembership(pg2.ID, []*pacta.Portfolio{p2})

	createMembership(pg2.ID, p3.ID)
	assertPGMembership(pg2.ID, []*pacta.Portfolio{p2, p3})

	deleteMembership(pg2.ID, p2.ID)

	actual, err := tdb.PortfolioGroups(tx, []pacta.PortfolioGroupID{pg1.ID, pg2.ID, pg3.ID})
	if err != nil {
		t.Fatalf("getting portfolio group: %v", err)
	}
	pg1.PortfolioGroupMemberships = []*pacta.PortfolioGroupMembership{
		{Portfolio: &pacta.Portfolio{ID: p1.ID}, CreatedAt: time.Now()},
		{Portfolio: &pacta.Portfolio{ID: p2.ID}, CreatedAt: time.Now()},
	}
	pg2.PortfolioGroupMemberships = []*pacta.PortfolioGroupMembership{
		{Portfolio: &pacta.Portfolio{ID: p3.ID}, CreatedAt: time.Now()},
	}
	expected := map[pacta.PortfolioGroupID]*pacta.PortfolioGroup{
		pg1.ID: pg1,
		pg2.ID: pg2,
		pg3.ID: pg3,
	}
	if diff := cmp.Diff(expected, actual, portfolioCmpOpts()); diff != "" {
		t.Fatalf("portfolio group mismatch (-want +got):\n%s", diff)
	}
}

func portfolioGroupCmpOpts() cmp.Option {
	portfolioGroupIDLessFn := func(a, b pacta.PortfolioGroupID) bool {
		return a < b
	}
	portfolioGroupLessFn := func(a, b *pacta.PortfolioGroup) bool {
		return a.ID < b.ID
	}
	portfolioMembershipLessFn := func(a, b *pacta.PortfolioGroupMembership) bool {
		if a.Portfolio != nil && b.Portfolio != nil {
			return a.Portfolio.ID < b.Portfolio.ID
		}
		if a.PortfolioGroup != nil && b.PortfolioGroup != nil {
			return a.PortfolioGroup.ID < b.PortfolioGroup.ID
		}
		return false // Fundamentally uncomparable.
	}
	return cmp.Options{
		cmpopts.SortSlices(portfolioGroupLessFn),
		cmpopts.SortSlices(portfolioMembershipLessFn),
		cmpopts.SortMaps(portfolioGroupIDLessFn),
		cmpopts.EquateEmpty(),
		cmpopts.EquateApproxTime(time.Second),
	}
}

func portfolioGroupForTesting(t *testing.T, tdb *DB, owner *pacta.Owner) *pacta.PortfolioGroup {
	t.Helper()
	ctx := context.Background()
	tx := tdb.NoTxn(ctx)
	pg := &pacta.PortfolioGroup{
		Name:  "portfolio-group-name",
		Owner: &pacta.Owner{ID: owner.ID},
	}
	pgID, err := tdb.CreatePortfolioGroup(tx, pg)
	if err != nil {
		t.Fatalf("creating portfolio group: %w", err)
	}
	pg.ID = pgID
	pg.CreatedAt = time.Now()
	return pg
}

func portfolioGroupMembershipForTesting(t *testing.T, tdb *DB, pg *pacta.PortfolioGroup, p *pacta.Portfolio) {
	t.Helper()
	ctx := context.Background()
	tx := tdb.NoTxn(ctx)
	if err := tdb.CreatePortfolioGroupMembership(tx, pg.ID, p.ID); err != nil {
		t.Fatalf("creating portfolio group membership: %v", err)
	}
}

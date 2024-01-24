package sqldb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/RMI/pacta/pacta"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCreatePortfolioInitiativeMembership(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	iA := initiativeForTestingWithKey(t, tdb, "A")
	iB := initiativeForTestingWithKey(t, tdb, "B")
	p1 := portfolioForTestingWithKey(t, tdb, "1")
	p2 := portfolioForTestingWithKey(t, tdb, "2")
	u := userForTesting(t, tdb)
	cmpOpts := portfolioInitiativeCmpOpts()

	pimA1 := &pacta.PortfolioInitiativeMembership{
		Initiative: &pacta.Initiative{ID: iA.ID},
		Portfolio:  &pacta.Portfolio{ID: p1.ID},
		AddedBy:    &pacta.User{ID: u.ID},
	}
	err := tdb.CreatePortfolioInitiativeMembership(tx, pimA1)
	if err != nil {
		t.Fatalf("creating pim: %v", err)
	}
	pimA1.CreatedAt = time.Now()

	pimA2 := &pacta.PortfolioInitiativeMembership{
		Initiative: &pacta.Initiative{ID: iA.ID},
		Portfolio:  &pacta.Portfolio{ID: p2.ID},
		AddedBy:    &pacta.User{ID: u.ID},
	}
	err = tdb.CreatePortfolioInitiativeMembership(tx, pimA2)
	if err != nil {
		t.Fatalf("creating pim: %v", err)
	}
	pimA2.CreatedAt = time.Now()

	pimB2 := &pacta.PortfolioInitiativeMembership{
		Initiative: &pacta.Initiative{ID: iB.ID},
		Portfolio:  &pacta.Portfolio{ID: p2.ID},
		AddedBy:    &pacta.User{ID: u.ID},
	}
	err = tdb.CreatePortfolioInitiativeMembership(tx, pimB2)
	if err != nil {
		t.Fatalf("creating pim: %v", err)
	}
	pimB2.CreatedAt = time.Now()

	actual, err := tdb.PortfolioInitiativeMembershipsByPortfolio(tx, p2.ID)
	if err != nil {
		t.Fatalf("getting pims: %v", err)
	}
	if diff := cmp.Diff([]*pacta.PortfolioInitiativeMembership{pimA2, pimB2}, actual, cmpOpts); diff != "" {
		t.Errorf("unexpected diff (-want +got)\n%s", diff)
	}

	actual, err = tdb.PortfolioInitiativeMembershipsByInitiative(tx, iA.ID)
	if err != nil {
		t.Fatalf("getting pims: %v", err)
	}
	if diff := cmp.Diff([]*pacta.PortfolioInitiativeMembership{pimA1, pimA2}, actual, cmpOpts); diff != "" {
		t.Errorf("unexpected diff (-want +got)\n%s", diff)
	}

	err = tdb.DeletePortfolioInitiativeMembership(tx, p1.ID, iA.ID)
	if err != nil {
		t.Fatalf("deleting pim: %v", err)
	}

	actual, err = tdb.PortfolioInitiativeMembershipsByPortfolio(tx, p1.ID)
	if err != nil {
		t.Fatalf("getting pims: %v", err)
	}
	if diff := cmp.Diff([]*pacta.PortfolioInitiativeMembership{}, actual, cmpOpts); diff != "" {
		t.Errorf("unexpected diff (-want +got)\n%s", diff)
	}
}

func portfolioInitiativeCmpOpts() cmp.Option {
	initiativeUserRelationshipLessFn := func(a, b *pacta.PortfolioInitiativeMembership) bool {
		if a.Portfolio.ID < b.Portfolio.ID {
			return true
		}
		if a.Initiative.ID > b.Initiative.ID {
			return false
		}
		return a.Initiative.ID < b.Initiative.ID
	}
	return cmp.Options{
		cmpopts.SortSlices(initiativeUserRelationshipLessFn),
		cmpopts.EquateEmpty(),
		cmpopts.EquateApproxTime(time.Second),
	}
}

func portfolioInitiativeForTesting(t *testing.T, tdb *DB, i *pacta.Initiative, p *pacta.Portfolio) {
	u := userForTestingWithKey(t, tdb, fmt.Sprintf("for portfolio_initiative_membership %s %s", p.ID, i.ID))
	if err := tdb.CreatePortfolioInitiativeMembership(tdb.NoTxn(context.Background()), &pacta.PortfolioInitiativeMembership{
		AddedBy:    &pacta.User{ID: u.ID},
		Initiative: &pacta.Initiative{ID: i.ID},
		Portfolio:  &pacta.Portfolio{ID: p.ID},
	}); err != nil {
		t.Fatalf("creating portfolio_initiative_membership: %v", err)
	}
}

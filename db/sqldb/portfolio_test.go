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

func TestPortfolioCRUD(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	b := blobForTesting(t, tdb)
	u1 := userForTestingWithKey(t, tdb, "1")
	u2 := userForTestingWithKey(t, tdb, "2")
	o1 := ownerUserForTesting(t, tdb, u1)
	o2 := ownerUserForTesting(t, tdb, u2)

	p := &pacta.Portfolio{
		Name:         "portfolio-name",
		Description:  "portfolio-description",
		HoldingsDate: exampleHoldingsDate,
		Owner:        &pacta.Owner{ID: o1.ID},
		Blob:         &pacta.Blob{ID: b.ID},
		NumberOfRows: 10,
	}
	id, err := tdb.CreatePortfolio(tx, p)
	if err != nil {
		t.Fatalf("creating portfolio: %v", err)
	}
	p.ID = id
	p.CreatedAt = time.Now()

	actual, err := tdb.Portfolio(tx, p.ID)
	if err != nil {
		t.Fatalf("reading portfolio: %v", err)
	}
	if diff := cmp.Diff(p, actual, portfolioCmpOpts()); diff != "" {
		t.Fatalf("portfolio mismatch (-want +got):\n%s", diff)
	}

	ps, err := tdb.Portfolios(tx, []pacta.PortfolioID{p.ID, p.ID})
	if err != nil {
		t.Fatalf("reading portfolios: %w", err)
	}
	if diff := cmp.Diff(map[pacta.PortfolioID]*pacta.Portfolio{p.ID: p}, ps, portfolioCmpOpts()); diff != "" {
		t.Fatalf("portfolio mismatch (-want +got):\n%s", diff)
	}

	nName := "new-name"
	nDesc := "new-description"
	nRows := 20
	err = tdb.UpdatePortfolio(tx, p.ID,
		db.SetPortfolioName(nName),
		db.SetPortfolioDescription(nDesc),
		db.SetPortfolioHoldingsDate(exampleHoldingsDate2),
		db.SetPortfolioOwner(o2.ID),
		db.SetPortfolioAdminDebugEnabled(true),
		db.SetPortfolioNumberOfRows(nRows),
	)
	if err != nil {
		t.Fatalf("updating portfolio: %v", err)
	}
	p.Name = nName
	p.Description = nDesc
	p.HoldingsDate = exampleHoldingsDate2
	p.Owner = &pacta.Owner{ID: o2.ID}
	p.AdminDebugEnabled = true
	p.NumberOfRows = nRows

	actual, err = tdb.Portfolio(tx, p.ID)
	if err != nil {
		t.Fatalf("reading portfolio: %v", err)
	}
	if diff := cmp.Diff(p, actual, portfolioCmpOpts()); diff != "" {
		t.Fatalf("portfolio mismatch (-want +got):\n%s", diff)
	}
	pls, err := tdb.PortfoliosByOwner(tx, o2.ID)
	if err != nil {
		t.Fatalf("reading portfolios: %w", err)
	}
	if diff := cmp.Diff([]*pacta.Portfolio{p}, pls, portfolioCmpOpts()); diff != "" {
		t.Fatalf("portfolio mismatch (-want +got):\n%s", diff)
	}
	pls, err = tdb.PortfoliosByOwner(tx, o1.ID)
	if err != nil {
		t.Fatalf("reading portfolios: %w", err)
	}
	if diff := cmp.Diff([]*pacta.Portfolio{}, pls, portfolioCmpOpts()); diff != "" {
		t.Fatalf("portfolio mismatch (-want +got):\n%s", diff)
	}

	blobOwners, err := tdb.BlobOwners(tx, []pacta.BlobID{b.ID})
	if err != nil {
		t.Fatalf("reading blob owners: %v", err)
	}
	expectedBlobOwners := []*pacta.BlobOwnerInformation{{
		BlobID:            b.ID,
		OwnerID:           o2.ID,
		AdminDebugEnabled: true,
	}}
	if diff := cmp.Diff(expectedBlobOwners, blobOwners, portfolioCmpOpts()); diff != "" {
		t.Errorf("unexpected diff (+got -want): %v", diff)
	}

	buris, err := tdb.DeletePortfolio(tx, p.ID)
	if err != nil {
		t.Fatalf("deleting portfolio: %v", err)
	}
	if diff := cmp.Diff([]pacta.BlobURI{b.BlobURI}, buris); diff != "" {
		t.Fatalf("blob uri mismatch (-want +got):\n%s", diff)
	}

	_, err = tdb.BlobOwners(tx, []pacta.BlobID{b.ID})
	if err == nil {
		t.Fatalf("reading blob owners should have failed but was fine", err)
	}
}

// TODO(grady) write a thorough portfolio deletion test

func portfolioCmpOpts() cmp.Option {
	portfolioIDLessFn := func(a, b pacta.PortfolioID) bool {
		return a < b
	}
	portfolioLessFn := func(a, b *pacta.Portfolio) bool {
		return a.ID < b.ID
	}
	return cmp.Options{
		cmpopts.SortSlices(portfolioLessFn),
		cmpopts.SortMaps(portfolioIDLessFn),
		cmpopts.EquateEmpty(),
		cmpopts.EquateApproxTime(time.Second),
	}
}

func portfolioForTesting(t *testing.T, tdb *DB) *pacta.Portfolio {
	t.Helper()
	return portfolioForTestingWithKey(t, tdb, "only")
}

func portfolioForTestingWithKey(t *testing.T, tdb *DB, key string) *pacta.Portfolio {
	t.Helper()
	ctx := context.Background()
	tx := tdb.NoTxn(ctx)
	b := blobForTestingWithKey(t, tdb, key)
	u := userForTestingWithKey(t, tdb, key)
	o := ownerUserForTesting(t, tdb, u)

	p := &pacta.Portfolio{
		Name:         "portfolio-name-" + key,
		Description:  "portfolio-description-" + key,
		HoldingsDate: exampleHoldingsDate,
		Owner:        &pacta.Owner{ID: o.ID},
		Blob:         &pacta.Blob{ID: b.ID},
		NumberOfRows: 10,
	}
	id, err := tdb.CreatePortfolio(tx, p)
	if err != nil {
		t.Fatalf("creating portfolio: %v", err)
	}
	p.CreatedAt = time.Now()
	p.ID = id
	return p
}

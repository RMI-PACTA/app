package sqldb

import (
	"context"
	"testing"
	"time"

	"github.com/RMI/pacta/pacta"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestPortfolioCRUD(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	b := blobForTesting(t, tdb)
	u := userForTesting(t, tdb)
	o := ownerUserForTesting(t, tdb, u)

	p := &pacta.Portfolio{
		Name:         "portfolio-name",
		Description:  "portfolio-description",
		HoldingsDate: exampleHoldingsDate,
		Owner:        o,
		Blob:         b,
		NumberOfRows: 10,
	}
	err := tdb.CreatePortfolio(tx, p)
	if err != nil {
		t.Fatalf("creating portfolio: %v", err)
	}
	p.CreatedAt = time.Now()

	/*
		assert := func(i *pacta.Portfolio) {
			t.Helper()
			actual, err := tdb.Portfolio(tx, i.ID)
			if err != nil {
				t.Fatalf("reading portfolio: %v", err)
			}
			if diff := cmp.Diff(i, actual, portfolioCmpOpts()); diff != "" {
				t.Fatalf("portfolio mismatch (-want +got):\n%s", diff)
			}
			eM := map[pacta.PortfolioID]*pacta.Portfolio{i.ID: i}
			aM, err := tdb.Portfolios(tx, []pacta.PortfolioID{i.ID})
			if err != nil {
				t.Fatalf("reading portfolios: %v", err)
			}
			if diff := cmp.Diff(eM, aM, portfolioCmpOpts()); diff != "" {
				t.Fatalf("portfolio mismatch (-want +got):\n%s", diff)
			}
			actuals, err := tdb.AllPortfolios(tx)
			if err != nil {
				t.Fatalf("reading portfolios: %v", err)
			}
			if diff := cmp.Diff([]*pacta.Portfolio{i}, actuals, portfolioCmpOpts()); diff != "" {
				t.Fatalf("portfolio mismatch (-want +got):\n%s", diff)
			}
		}
		assert(i)

		i.Name = "new name"
		i.Affiliation = "new affiliation"
		i.PublicDescription = "new public decsription"
		i.InternalDescription = "new internal description"
		i.RequiresPortfolioToJoin = true
		i.Language = pacta.Language_EN
		err = tdb.UpdatePortfolio(tx, i.ID,
			db.SetPortfolioName(i.Name),
			db.SetPortfolioAffiliation(i.Affiliation),
			db.SetPortfolioPublicDescription(i.PublicDescription),
			db.SetPortfolioInternalDescription(i.InternalDescription),
			db.SetPortfolioRequiresPortfolioToJoin(i.RequiresPortfolioToJoin),
			db.SetPortfolioLanguage(pacta.Language_EN),
		)
		if err != nil {
			t.Fatalf("updating portfolio: %v", err)
		}
		assert(i)

		i.IsAcceptingNewMembers = true
		if err := tdb.UpdatePortfolio(tx, i.ID, db.SetPortfolioIsAcceptingNewMembers(true)); err != nil {
			t.Fatalf("updating portfolio: %v", err)
		}
		assert(i)

		i.IsAcceptingNewPortfolios = true
		if err := tdb.UpdatePortfolio(tx, i.ID, db.SetPortfolioIsAcceptingNewPortfolios(true)); err != nil {
			t.Fatalf("updating portfolio: %v", err)
		}
		assert(i)

		err = tdb.DeletePortfolio(tx, i.ID)
		if err != nil {
			t.Fatalf("delete portfolio: %v", err)
		}
	*/
}

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

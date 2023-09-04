package sqldb

/*
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



	i := &pacta.Portfolio{
		ID:           "portfolio-id",
		Language:     pacta.Language_DE,
		Name:         "portfolio-name",
		PACTAVersion: &pacta.PACTAVersion{ID: pvID},
	}
	err := tdb.CreatePortfolio(tx, i)
	if err != nil {
		t.Fatalf("creating portfolio: %v", err)
	}
	i.CreatedAt = time.Now()

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
}

func TestDeletePortfolio(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	pv := &pacta.PACTAVersion{
		Name:        "pacta version",
		Description: "pacta description",
		Digest:      "pacta digest",
	}
	pvID, err0 := tdb.CreatePACTAVersion(tx, pv)
	i := &pacta.Portfolio{
		ID:           "portfolio-id",
		Language:     pacta.Language_DE,
		Name:         "portfolio-name",
		PACTAVersion: &pacta.PACTAVersion{ID: pvID},
	}
	err1 := tdb.CreatePortfolio(tx, i)
	_, err2 := tdb.CreatePortfolioPortfolio(tx, &pacta.PortfolioPortfolio{
		Portfolio: &pacta.Portfolio{ID: i.ID},
	})
	uid, err3 := tdb.CreateUser(tx, &pacta.User{
		AuthnMechanism: pacta.AuthnMechanism_EmailAndPass,
		AuthnID:        "authn-id",
		EnteredEmail:   "entered-email",
		CanonicalEmail: "canonical-email",
		Name:           "User's Name",
	})
	iur := &pacta.PortfolioUserRelationship{
		User:       &pacta.User{ID: uid},
		Portfolio: &pacta.Portfolio{ID: i.ID},
	}
	err4 := tdb.PutPortfolioUserRelationship(tx, iur)
	noErrDuringSetup(t, err0, err1, err2, err3, err4)

	err := tdb.DeletePortfolio(tx, i.ID)
	if err != nil {
		t.Fatalf("delete portfolio: %v", err)
	}

	_, err = tdb.Portfolio(tx, i.ID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
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
*/

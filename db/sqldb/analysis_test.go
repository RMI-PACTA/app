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

func TestAnalysisCRUD(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	pv := pactaVersionForTesting(t, tdb)
	u1 := userForTestingWithKey(t, tdb, "User1")
	o1 := ownerUserForTesting(t, tdb, u1)
	u2 := userForTestingWithKey(t, tdb, "User2")
	o2 := ownerUserForTesting(t, tdb, u2)
	p1 := portfolioForTestingWithKey(t, tdb, "Portfolio1")
	p2 := portfolioForTestingWithKey(t, tdb, "Portfolio2")
	pg := portfolioGroupForTesting(t, tdb, o1)
	portfolioGroupMembershipForTesting(t, tdb, pg, p1)
	portfolioGroupMembershipForTesting(t, tdb, pg, p2)
	s := snapshotPortfolioGroupForTesting(t, tdb, pg)
	cmpOpts := analysisCmpOpts()

	iu := &pacta.Analysis{
		PortfolioSnapshot: s,
		PACTAVersion:      &pacta.PACTAVersion{ID: pv.ID},
		Name:              "analysis-name",
		Description:       "analysis-description",
		Owner:             &pacta.Owner{ID: o1.ID},
		AnalysisType:      pacta.AnalysisType_Audit,
	}
	id, err := tdb.CreateAnalysis(tx, iu)
	if err != nil {
		t.Fatalf("creating analysis: %v", err)
	}
	iu.ID = id
	iu.CreatedAt = time.Now()

	actual, err := tdb.Analysis(tx, iu.ID)
	if err != nil {
		t.Fatalf("reading analysis: %v", err)
	}
	if diff := cmp.Diff(iu, actual, cmpOpts); diff != "" {
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}

	ius, err := tdb.Analyses(tx, []pacta.AnalysisID{iu.ID, iu.ID, "nonsense"})
	if err != nil {
		t.Fatalf("reading analysiss: %w", err)
	}
	if diff := cmp.Diff(map[pacta.AnalysisID]*pacta.Analysis{iu.ID: iu}, ius, cmpOpts); diff != "" {
		t.Fatalf("analysis mismatch (-want +got):\n%s", diff)
	}

	nName := "new-name"
	nDesc := "new-description"
	ranAt := time.UnixMilli(111111111)
	completedAt := time.UnixMilli(222222222)
	failureCode := pacta.FailureCode_Unknown
	failureMessage := "failureMessage"
	err = tdb.UpdateAnalysis(tx, iu.ID,
		db.SetAnalysisName(nName),
		db.SetAnalysisDescription(nDesc),
		db.SetAnalysisRanAt(ranAt),
		db.SetAnalysisOwner(o2.ID),
		db.SetAnalysisFailureCode(failureCode),
		db.SetAnalysisCompletedAt(completedAt),
		db.SetAnalysisFailureMessage(failureMessage),
	)
	if err != nil {
		t.Fatalf("updating portfolio: %v", err)
	}
	iu.Name = nName
	iu.Description = nDesc
	iu.RanAt = ranAt
	iu.Owner = &pacta.Owner{ID: o2.ID}
	iu.FailureCode = failureCode
	iu.CompletedAt = completedAt
	iu.FailureMessage = failureMessage

	actual, err = tdb.Analysis(tx, iu.ID)
	if err != nil {
		t.Fatalf("reading portfolio: %v", err)
	}
	if diff := cmp.Diff(iu, actual, cmpOpts); diff != "" {
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}

	buris, err := tdb.DeleteAnalysis(tx, iu.ID)
	if err != nil {
		t.Fatalf("deleting analysis: %v", err)
	}
	if diff := cmp.Diff([]pacta.BlobURI{}, buris); diff != "" {
		t.Fatalf("blob uri mismatch (-want +got):\n%s", diff)
	}
}

func analysisCmpOpts() cmp.Option {
	return cmp.Options{
		cmpopts.EquateEmpty(),
		cmpopts.EquateApproxTime(time.Second),
	}
}

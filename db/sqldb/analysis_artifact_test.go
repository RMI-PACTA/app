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

func TestAnalysisArtifacts(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	pv := pactaVersionForTesting(t, tdb)
	u := userForTestingWithKey(t, tdb, "User1")
	o := ownerUserForTesting(t, tdb, u)
	p1 := portfolioForTestingWithKey(t, tdb, "Portfolio1")
	p2 := portfolioForTestingWithKey(t, tdb, "Portfolio2")
	pg := portfolioGroupForTesting(t, tdb, o)
	portfolioGroupMembershipForTesting(t, tdb, pg, p1)
	portfolioGroupMembershipForTesting(t, tdb, pg, p2)
	s := snapshotPortfolioGroupForTesting(t, tdb, pg)
	b1 := blobForTestingWithKey(t, tdb, "blob1")
	b2 := blobForTestingWithKey(t, tdb, "blob2")
	b3 := blobForTestingWithKey(t, tdb, "blob3")
	aid, err := tdb.CreateAnalysis(tx, &pacta.Analysis{
		PortfolioSnapshot: &pacta.PortfolioSnapshot{ID: s.ID},
		PACTAVersion:      &pacta.PACTAVersion{ID: pv.ID},
		Name:              "analysis-name",
		Description:       "analysis-description",
		Owner:             &pacta.Owner{ID: o.ID},
		AnalysisType:      pacta.AnalysisType_Audit,
	})
	if err != nil {
		t.Fatalf("creating analysis: %v", err)
	}
	cmpOpts := analysisArtifactCmpOpts()

	aa1 := &pacta.AnalysisArtifact{
		AnalysisID: aid,
		Blob:       &pacta.Blob{ID: b1.ID},
	}
	aa1.ID, err = tdb.CreateAnalysisArtifact(tx, aa1)
	if err != nil {
		t.Fatalf("creating analysis artifact: %v", err)
	}
	aa2 := &pacta.AnalysisArtifact{
		AnalysisID: aid,
		Blob:       &pacta.Blob{ID: b2.ID},
	}
	aa2.ID, err = tdb.CreateAnalysisArtifact(tx, aa2)
	if err != nil {
		t.Fatalf("creating analysis artifact: %v", err)
	}
	aa3 := &pacta.AnalysisArtifact{
		AnalysisID: aid,
		Blob:       &pacta.Blob{ID: b3.ID},
	}
	aa3.ID, err = tdb.CreateAnalysisArtifact(tx, aa3)
	if err != nil {
		t.Fatalf("creating analysis artifact: %v", err)
	}

	if err := tdb.UpdateAnalysisArtifact(tx, aa2.ID, db.SetAnalysisArtifactAdminDebugEnabled(true)); err != nil {
		t.Fatalf("updating analysis artifact: %v", err)
	}
	if err := tdb.UpdateAnalysisArtifact(tx, aa3.ID, db.SetAnalysisArtifactSharedToPublic(true)); err != nil {
		t.Fatalf("updating analysis artifact: %v", err)
	}

	actualAA1, err := tdb.AnalysisArtifact(tx, aa1.ID)
	if err != nil {
		t.Fatalf("reading analysis artifact: %v", err)
	}
	if diff := cmp.Diff(aa1, actualAA1, cmpOpts); diff != "" {
		t.Errorf("unexpected analysis artifact: %v", diff)
	}

	actualArtifacts, err := tdb.AnalysisArtifactsForAnalysis(tx, aid)
	if err != nil {
		t.Fatalf("reading analysis artifacts: %v", err)
	}
	if diff := cmp.Diff(actualArtifacts, []*pacta.AnalysisArtifact{aa1, aa2, aa3}, cmpOpts); diff != "" {
		t.Errorf("unexpected diff (+got -want): %v", diff)
	}

	listedArtifacts, err := tdb.AnalysisArtifacts(tx, []pacta.AnalysisArtifactID{aa1.ID, aa3.ID})
	if err != nil {
		t.Fatalf("reading analysis artifacts: %v", err)
	}
	if diff := cmp.Diff([]*pacta.AnalysisArtifact{aa1, aa3}, listedArtifacts, cmpOpts); diff != "" {
		t.Errorf("unexpected diff (+got -want): %v", diff)
	}

	buris, err := tdb.DeleteAnalysis(tx, aid)
	if err != nil {
		t.Fatalf("deleting analysis: %v", err)
	}
	if diff := cmp.Diff([]pacta.BlobURI{b1.BlobURI, b2.BlobURI, b3.BlobURI}, buris, cmpOpts); diff != "" {
		t.Errorf("unexpected diff (+got -want): %v", diff)
	}
}

func analysisArtifactCmpOpts() cmp.Option {
	blobURILessFn := func(a, b pacta.BlobURI) bool {
		return a < b
	}
	aaLessFn := func(a, b *pacta.AnalysisArtifact) bool {
		return a.ID < b.ID
	}
	return cmp.Options{
		cmpopts.EquateEmpty(),
		cmpopts.EquateApproxTime(time.Second),
		cmpopts.SortSlices(blobURILessFn),
		cmpopts.SortSlices(aaLessFn),
	}
}

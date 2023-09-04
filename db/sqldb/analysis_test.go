package sqldb

/*
func TestAnalysisCRUD(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	pv := pactaVersionForTesting(t, tdb)
	a := &pacta.Analysis{
		ID:           "initiative-id",
		Language:     pacta.Language_DE,
		Name:         "initiative-name",
		PACTAVersion: &pacta.PACTAVersion{ID: pv.ID},
	}
	err := tdb.CreateAnalysis(tx, i)
	if err != nil {
		t.Fatalf("creating initiative: %v", err)
	}
	i.CreatedAt = time.Now()

	assert := func(i *pacta.Analysis) {
		t.Helper()
		actual, err := tdb.Analysis(tx, i.ID)
		if err != nil {
			t.Fatalf("reading initiative: %v", err)
		}
		if diff := cmp.Diff(i, actual, initiativeCmpOpts()); diff != "" {
			t.Fatalf("initiative mismatch (-want +got):\n%s", diff)
		}
		eM := map[pacta.AnalysisID]*pacta.Analysis{i.ID: i}
		aM, err := tdb.Analysiss(tx, []pacta.AnalysisID{i.ID})
		if err != nil {
			t.Fatalf("reading initiatives: %v", err)
		}
		if diff := cmp.Diff(eM, aM, initiativeCmpOpts()); diff != "" {
			t.Fatalf("initiative mismatch (-want +got):\n%s", diff)
		}
		actuals, err := tdb.AllAnalysiss(tx)
		if err != nil {
			t.Fatalf("reading initiatives: %v", err)
		}
		if diff := cmp.Diff([]*pacta.Analysis{i}, actuals, initiativeCmpOpts()); diff != "" {
			t.Fatalf("initiative mismatch (-want +got):\n%s", diff)
		}
	}
	assert(i)

	i.Name = "new name"
	i.Affiliation = "new affiliation"
	i.PublicDescription = "new public decsription"
	i.InternalDescription = "new internal description"
	i.RequiresInvitationToJoin = true
	i.Language = pacta.Language_EN
	err = tdb.UpdateAnalysis(tx, i.ID,
		db.SetAnalysisName(i.Name),
		db.SetAnalysisAffiliation(i.Affiliation),
		db.SetAnalysisPublicDescription(i.PublicDescription),
		db.SetAnalysisInternalDescription(i.InternalDescription),
		db.SetAnalysisRequiresInvitationToJoin(i.RequiresInvitationToJoin),
		db.SetAnalysisLanguage(pacta.Language_EN),
	)
	if err != nil {
		t.Fatalf("updating initiative: %v", err)
	}
	assert(i)

	i.IsAcceptingNewMembers = true
	if err := tdb.UpdateAnalysis(tx, i.ID, db.SetAnalysisIsAcceptingNewMembers(true)); err != nil {
		t.Fatalf("updating initiative: %v", err)
	}
	assert(i)

	i.IsAcceptingNewPortfolios = true
	if err := tdb.UpdateAnalysis(tx, i.ID, db.SetAnalysisIsAcceptingNewPortfolios(true)); err != nil {
		t.Fatalf("updating initiative: %v", err)
	}
	assert(i)

	err = tdb.DeleteAnalysis(tx, i.ID)
	if err != nil {
		t.Fatalf("delete initiative: %v", err)
	}
}

func TestDeleteAnalysis(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	i := initiativeForTesting(t, tdb)
	u := userForTesting(t, tdb)
	_, err0 := tdb.CreateAnalysisInvitation(tx, &pacta.AnalysisInvitation{
		Analysis: &pacta.Analysis{ID: i.ID},
	})
	iur := &pacta.AnalysisUserRelationship{
		User:       &pacta.User{ID: u.ID},
		Analysis: &pacta.Analysis{ID: i.ID},
	}
	err1 := tdb.PutAnalysisUserRelationship(tx, iur)
	noErrDuringSetup(t, err0, err1)

	err := tdb.DeleteAnalysis(tx, i.ID)
	if err != nil {
		t.Fatalf("delete initiative: %v", err)
	}

	_, err = tdb.Analysis(tx, i.ID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func initiativeCmpOpts() cmp.Option {
	initiativeIDLessFn := func(a, b pacta.AnalysisID) bool {
		return a < b
	}
	initiativeLessFn := func(a, b *pacta.Analysis) bool {
		return a.ID < b.ID
	}
	return cmp.Options{
		cmpopts.SortSlices(initiativeLessFn),
		cmpopts.SortMaps(initiativeIDLessFn),
		cmpopts.EquateEmpty(),
		cmpopts.EquateApproxTime(time.Second),
	}
}

func initiativeForTesting(t *testing.T, tdb *DB) *pacta.Analysis {
	t.Helper()
	pv := pactaVersionForTesting(t, tdb)
	i := &pacta.Analysis{
		ID:           "initiative-id",
		Language:     pacta.Language_DE,
		Name:         "initiative-name",
		PACTAVersion: &pacta.PACTAVersion{ID: pv.ID},
	}
	ctx := context.Background()
	tx := tdb.NoTxn(ctx)
	err := tdb.CreateAnalysis(tx, i)
	if err != nil {
		t.Fatalf("creating initiative: %v", err)
	}
	i.CreatedAt = time.Now()
	return i
}
*/

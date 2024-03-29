package sqldb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestInitiativeCRUD(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	pv := pactaVersionForTesting(t, tdb)
	i := &pacta.Initiative{
		ID:           "initiative-id",
		Language:     pacta.Language_DE,
		Name:         "initiative-name",
		PACTAVersion: &pacta.PACTAVersion{ID: pv.ID},
	}
	err := tdb.CreateInitiative(tx, i)
	if err != nil {
		t.Fatalf("creating initiative: %v", err)
	}
	i.CreatedAt = time.Now()

	assert := func(i *pacta.Initiative) {
		t.Helper()
		actual, err := tdb.Initiative(tx, i.ID)
		if err != nil {
			t.Fatalf("reading initiative: %v", err)
		}
		if diff := cmp.Diff(i, actual, initiativeCmpOpts()); diff != "" {
			t.Fatalf("initiative mismatch (-want +got):\n%s", diff)
		}
		eM := map[pacta.InitiativeID]*pacta.Initiative{i.ID: i}
		aM, err := tdb.Initiatives(tx, []pacta.InitiativeID{i.ID})
		if err != nil {
			t.Fatalf("reading initiatives: %v", err)
		}
		if diff := cmp.Diff(eM, aM, initiativeCmpOpts()); diff != "" {
			t.Fatalf("initiative mismatch (-want +got):\n%s", diff)
		}
		actuals, err := tdb.AllInitiatives(tx)
		if err != nil {
			t.Fatalf("reading initiatives: %v", err)
		}
		if diff := cmp.Diff([]*pacta.Initiative{i}, actuals, initiativeCmpOpts()); diff != "" {
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
	err = tdb.UpdateInitiative(tx, i.ID,
		db.SetInitiativeName(i.Name),
		db.SetInitiativeAffiliation(i.Affiliation),
		db.SetInitiativePublicDescription(i.PublicDescription),
		db.SetInitiativeInternalDescription(i.InternalDescription),
		db.SetInitiativeRequiresInvitationToJoin(i.RequiresInvitationToJoin),
		db.SetInitiativeLanguage(pacta.Language_EN),
	)
	if err != nil {
		t.Fatalf("updating initiative: %v", err)
	}
	assert(i)

	i.IsAcceptingNewMembers = true
	if err := tdb.UpdateInitiative(tx, i.ID, db.SetInitiativeIsAcceptingNewMembers(true)); err != nil {
		t.Fatalf("updating initiative: %v", err)
	}
	assert(i)

	i.IsAcceptingNewPortfolios = true
	if err := tdb.UpdateInitiative(tx, i.ID, db.SetInitiativeIsAcceptingNewPortfolios(true)); err != nil {
		t.Fatalf("updating initiative: %v", err)
	}
	assert(i)

	buris, err := tdb.DeleteInitiative(tx, i.ID)
	if err != nil {
		t.Fatalf("delete initiative: %v", err)
	}
	if len(buris) != 0 {
		t.Fatalf("expected no deleted buris, got %d", len(buris))
	}
}

func TestDeleteInitiative(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	i := initiativeForTesting(t, tdb)
	u := userForTesting(t, tdb)
	_, err0 := tdb.CreateInitiativeInvitation(tx, &pacta.InitiativeInvitation{
		Initiative: &pacta.Initiative{ID: i.ID},
	})
	iur := &pacta.InitiativeUserRelationship{
		User:       &pacta.User{ID: u.ID},
		Initiative: &pacta.Initiative{ID: i.ID},
	}
	err1 := tdb.PutInitiativeUserRelationship(tx, iur)
	noErrDuringSetup(t, err0, err1)

	buris, err := tdb.DeleteInitiative(tx, i.ID)
	if err != nil {
		t.Fatalf("delete initiative: %v", err)
	}

	if len(buris) != 0 {
		t.Fatalf("expected no buris but got %+v", buris)
	}

	_, err = tdb.Initiative(tx, i.ID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func initiativeCmpOpts() cmp.Option {
	initiativeIDLessFn := func(a, b pacta.InitiativeID) bool {
		return a < b
	}
	initiativeLessFn := func(a, b *pacta.Initiative) bool {
		return a.ID < b.ID
	}
	return cmp.Options{
		cmpopts.SortSlices(initiativeLessFn),
		cmpopts.SortMaps(initiativeIDLessFn),
		cmpopts.EquateEmpty(),
		cmpopts.EquateApproxTime(time.Second),
	}
}

func initiativeForTesting(t *testing.T, tdb *DB) *pacta.Initiative {
	t.Helper()
	return initiativeForTestingWithKey(t, tdb, "only")
}

func initiativeForTestingWithKey(t *testing.T, tdb *DB, key string) *pacta.Initiative {
	t.Helper()
	pv := pactaVersionForTestingWithKey(t, tdb, key)
	i := &pacta.Initiative{
		ID:           pacta.InitiativeID(fmt.Sprintf("initiative-id-%s", key)),
		Language:     pacta.Language_DE,
		Name:         "initiative-name",
		PACTAVersion: &pacta.PACTAVersion{ID: pv.ID},
	}
	ctx := context.Background()
	tx := tdb.NoTxn(ctx)
	err := tdb.CreateInitiative(tx, i)
	if err != nil {
		t.Fatalf("creating initiative: %v", err)
	}
	i.CreatedAt = time.Now()
	return i
}

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

func TestInitiativeCRUD(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	pv := &pacta.PACTAVersion{
		Name:        "pacta version",
		Description: "pacta description",
		Digest:      "pacta digest",
	}
	pvID, err0 := tdb.CreatePACTAVersion(tx, pv)
	noErrDuringSetup(t, err0)

	i := &pacta.Initiative{
		ID:           "initiative-id",
		Language:     pacta.Language_DE,
		Name:         "initiative-name",
		PACTAVersion: &pacta.PACTAVersion{ID: pvID},
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

	err = tdb.DeleteInitiative(tx, i.ID)
	if err != nil {
		t.Fatalf("delete initiative: %v", err)
	}
}

func TestDeleteInitiative(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	pv := &pacta.PACTAVersion{
		Name:        "pacta version",
		Description: "pacta description",
		Digest:      "pacta digest",
	}
	pvID, err0 := tdb.CreatePACTAVersion(tx, pv)
	i := &pacta.Initiative{
		ID:           "initiative-id",
		Language:     pacta.Language_DE,
		Name:         "initiative-name",
		PACTAVersion: &pacta.PACTAVersion{ID: pvID},
	}
	err1 := tdb.CreateInitiative(tx, i)
	_, err2 := tdb.CreateInitiativeInvitation(tx, &pacta.InitiativeInvitation{
		Initiative: &pacta.Initiative{ID: i.ID},
	})
	uid, err3 := tdb.CreateUser(tx, &pacta.User{
		AuthnMechanism: pacta.AuthnMechanism_EmailAndPass,
		AuthnID:        "authn-id",
		EnteredEmail:   "entered-email",
		CanonicalEmail: "canonical-email",
		Name:           "User's Name",
	})
	iur := &pacta.InitiativeUserRelationship{
		User:       &pacta.User{ID: uid},
		Initiative: &pacta.Initiative{ID: i.ID},
	}
	err4 := tdb.PutInitiativeUserRelationship(tx, iur)
	noErrDuringSetup(t, err0, err1, err2, err3, err4)

	err := tdb.DeleteInitiative(tx, i.ID)
	if err != nil {
		t.Fatalf("delete initiative: %v", err)
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

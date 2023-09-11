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

func TestCreateInitiativeInvitation(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	i := initiativeForTesting(t, tdb)

	var presetID pacta.InitiativeInvitationID = "PresetID"
	ii := &pacta.InitiativeInvitation{
		ID:         presetID,
		Initiative: &pacta.Initiative{ID: i.ID},
	}
	id, err := tdb.CreateInitiativeInvitation(tx, ii)
	if err != nil {
		t.Fatalf("creating initiative_invitation: %v", err)
	}
	if id != presetID {
		t.Fatalf("expected %q got %q", presetID, id)
	}

	// Create should fail with the same initiative invitation id
	_, err = tdb.CreateInitiativeInvitation(tx, ii)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Create without a preset id
	ii2 := &pacta.InitiativeInvitation{
		Initiative: &pacta.Initiative{ID: i.ID},
	}
	id2, err := tdb.CreateInitiativeInvitation(tx, ii2)
	if err != nil {
		t.Fatalf("creating initiative_invitation: %v", err)
	}
	if id2 == "" || id2 == id {
		t.Fatalf("expected a new id, got %q", id2)
	}

	// Read By ID
	actual, err := tdb.InitiativeInvitation(tx, id2)
	if err != nil {
		t.Fatalf("getting initiative invitation: %v", err)
	}
	ii2.CreatedAt = time.Now()
	if diff := cmp.Diff(ii2, actual, initiativeInvitationCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}

	// Read By Initiative
	actualIIs, err := tdb.InitiativeInvitationsByInitiative(tx, i.ID)
	if err != nil {
		t.Fatalf("getting initiative invitations: %v", err)
	}
	ii.CreatedAt = time.Now()
	if diff := cmp.Diff([]*pacta.InitiativeInvitation{ii, ii2}, actualIIs, initiativeInvitationCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}
}

func TestUpdateInitiativeInvitation(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	i := initiativeForTesting(t, tdb)
	u := userForTesting(t, tdb)
	iiid, err0 := tdb.CreateInitiativeInvitation(tx, &pacta.InitiativeInvitation{
		Initiative: &pacta.Initiative{ID: i.ID},
	})
	noErrDuringSetup(t, err0)

	err := tdb.UpdateInitiativeInvitation(tx, iiid,
		db.SetInitiativeInvitationUsedBy(u.ID),
		db.SetInitiativeInvitationUsedAt(time.Now()))
	if err != nil {
		t.Fatalf("update initiative invitation: %v", err)
	}

	actual, err := tdb.InitiativeInvitation(tx, iiid)
	if err != nil {
		t.Fatalf("getting initiative invitation: %v", err)
	}
	expected := &pacta.InitiativeInvitation{
		ID:         iiid,
		Initiative: &pacta.Initiative{ID: i.ID},
		CreatedAt:  time.Now(),
		UsedAt:     time.Now(),
		UsedBy:     &pacta.User{ID: u.ID},
	}
	if diff := cmp.Diff(expected, actual, initiativeInvitationCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}
}

func TestDeleteInitiativeInvitation(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	i := initiativeForTesting(t, tdb)
	iiid, err0 := tdb.CreateInitiativeInvitation(tx, &pacta.InitiativeInvitation{
		Initiative: &pacta.Initiative{ID: i.ID},
	})
	noErrDuringSetup(t, err0)

	err := tdb.DeleteInitiativeInvitation(tx, iiid)
	if err != nil {
		t.Fatalf("delete initiative invitation: %v", err)
	}

	_, err = tdb.InitiativeInvitation(tx, iiid)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func initiativeInvitationCmpOpts() cmp.Option {
	initiativeInvitationIDLessFn := func(a, b pacta.InitiativeInvitationID) bool {
		return a < b
	}
	initiativeInvitationLessFn := func(a, b *pacta.InitiativeInvitation) bool {
		return a.ID < b.ID
	}
	return cmp.Options{
		cmpopts.SortSlices(initiativeInvitationLessFn),
		cmpopts.SortMaps(initiativeInvitationIDLessFn),
		cmpopts.EquateEmpty(),
		cmpopts.EquateApproxTime(time.Second),
	}
}

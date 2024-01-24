package sqldb

import (
	"context"
	"testing"
	"time"

	"github.com/RMI/pacta/pacta"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCreateInitiativeUserRelationship(t *testing.T) {
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
	uid, err2 := tdb.createUser(tx, &pacta.User{
		CanonicalEmail: "canon",
		EnteredEmail:   "entered",
		AuthnMechanism: pacta.AuthnMechanism_EmailAndPass,
		AuthnID:        "A",
	})
	noErrDuringSetup(t, err0, err1, err2)

	iur := &pacta.InitiativeUserRelationship{
		User:       &pacta.User{ID: uid},
		Initiative: &pacta.Initiative{ID: i.ID},
	}

	err := tdb.PutInitiativeUserRelationship(tx, iur)
	if err != nil {
		t.Fatalf("creating initiative_user_relationship: %v", err)
	}

	// Read By ID
	actual, err := tdb.InitiativeUserRelationship(tx, i.ID, uid)
	if err != nil {
		t.Fatalf("getting initiative invitation: %v", err)
	}
	expected := iur.Clone()
	expected.UpdatedAt = time.Now()
	if diff := cmp.Diff(expected, actual, initiativeUserRelationshipCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}
}

func TestUpdateInitiativeUserRelationship(t *testing.T) {
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
	uid, err2 := tdb.createUser(tx, &pacta.User{
		CanonicalEmail: "canon",
		EnteredEmail:   "entered",
		AuthnMechanism: pacta.AuthnMechanism_EmailAndPass,
		AuthnID:        "A",
	})
	iur := &pacta.InitiativeUserRelationship{
		User:       &pacta.User{ID: uid},
		Initiative: &pacta.Initiative{ID: i.ID},
	}
	err3 := tdb.PutInitiativeUserRelationship(tx, iur)
	noErrDuringSetup(t, err0, err1, err2, err3)

	iur.Manager = true
	iur.UpdatedAt = time.Now()
	err := tdb.PutInitiativeUserRelationship(tx, iur)
	if err != nil {
		t.Fatalf("update initiative user relationship: %v", err)
	}
	actual, err := tdb.InitiativeUserRelationship(tx, iur.Initiative.ID, iur.User.ID)
	if err != nil {
		t.Fatalf("getting initiative user relationship: %v", err)
	}
	if diff := cmp.Diff(iur, actual, initiativeUserRelationshipCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}

	iur.Member = true
	iur.Manager = false
	iur.UpdatedAt = time.Now()
	err = tdb.PutInitiativeUserRelationship(tx, iur)
	if err != nil {
		t.Fatalf("update initiative user relationship: %v", err)
	}
	actual, err = tdb.InitiativeUserRelationship(tx, iur.Initiative.ID, iur.User.ID)
	if err != nil {
		t.Fatalf("getting initiative user relationship: %v", err)
	}
	if diff := cmp.Diff(iur, actual, initiativeUserRelationshipCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}
}

func TestListInitiativeUserRelationships(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	pv := &pacta.PACTAVersion{
		Name:        "pacta version",
		Description: "pacta description",
		Digest:      "pacta digest",
	}
	pvID, err0 := tdb.CreatePACTAVersion(tx, pv)
	i1 := &pacta.Initiative{
		ID:           "initiative-id",
		Language:     pacta.Language_DE,
		Name:         "initiative-name",
		PACTAVersion: &pacta.PACTAVersion{ID: pvID},
	}
	err1 := tdb.CreateInitiative(tx, i1)
	i2 := &pacta.Initiative{
		ID:           "initiative-id-2",
		Language:     pacta.Language_EN,
		Name:         "initiative-name-2",
		PACTAVersion: &pacta.PACTAVersion{ID: pvID},
	}
	err2 := tdb.CreateInitiative(tx, i2)
	u1, err3 := tdb.createUser(tx, &pacta.User{
		CanonicalEmail: "canon",
		EnteredEmail:   "entered",
		AuthnMechanism: pacta.AuthnMechanism_EmailAndPass,
		AuthnID:        "A",
	})
	u2, err4 := tdb.createUser(tx, &pacta.User{
		CanonicalEmail: "canon2",
		EnteredEmail:   "entered2",
		AuthnMechanism: pacta.AuthnMechanism_EmailAndPass,
		AuthnID:        "B",
	})
	iurI1U1 := &pacta.InitiativeUserRelationship{
		User:       &pacta.User{ID: u1},
		Initiative: &pacta.Initiative{ID: i1.ID},
		UpdatedAt:  time.Now(),
	}
	err5 := tdb.PutInitiativeUserRelationship(tx, iurI1U1)
	iurI2U1 := &pacta.InitiativeUserRelationship{
		User:       &pacta.User{ID: u1},
		Initiative: &pacta.Initiative{ID: i2.ID},
		Member:     true,
		UpdatedAt:  time.Now(),
	}
	err6 := tdb.PutInitiativeUserRelationship(tx, iurI2U1)
	iurI2U2 := &pacta.InitiativeUserRelationship{
		User:       &pacta.User{ID: u2},
		Initiative: &pacta.Initiative{ID: i2.ID},
		Manager:    true,
		UpdatedAt:  time.Now(),
	}
	err7 := tdb.PutInitiativeUserRelationship(tx, iurI2U2)
	noErrDuringSetup(t, err0, err1, err2, err3, err4, err5, err6, err7)

	actual, err := tdb.InitiativeUserRelationshipsByInitiative(tx, i2.ID)
	if err != nil {
		t.Fatalf("getting initiative user relationships: %v", err)
	}
	if diff := cmp.Diff([]*pacta.InitiativeUserRelationship{iurI2U1, iurI2U2}, actual, initiativeUserRelationshipCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}

	actual, err = tdb.InitiativeUserRelationshipsByUser(tx, u2)
	if err != nil {
		t.Fatalf("getting initiative user relationships: %v", err)
	}
	if diff := cmp.Diff([]*pacta.InitiativeUserRelationship{iurI2U2}, actual, initiativeUserRelationshipCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}

	actual, err = tdb.InitiativeUserRelationshipsByUser(tx, u1)
	if err != nil {
		t.Fatalf("getting initiative user relationships: %v", err)
	}
	if diff := cmp.Diff([]*pacta.InitiativeUserRelationship{iurI2U1, iurI1U1}, actual, initiativeUserRelationshipCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}
}

func initiativeUserRelationshipCmpOpts() cmp.Option {
	initiativeUserRelationshipLessFn := func(a, b *pacta.InitiativeUserRelationship) bool {
		if a.User.ID < b.User.ID {
			return true
		}
		if a.User.ID > b.User.ID {
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

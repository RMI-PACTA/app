package sqldb

import (
	"context"
	"testing"
	"time"

	"github.com/RMI/pacta"
	"github.com/RMI/pacta/db"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	email := "user@example.com"
	userID, err := tdb.CreateUser(tx, pacta.EmailAndPass, pacta.UserID(email), "User's Name", email)
	if err != nil {
		t.Fatalf("creating user: %v", err)
	}

	actual, err := tdb.User(tx, userID)
	if err != nil {
		t.Fatalf("getting user: %v", err)
	}
	expected := &pacta.User{
		ID:                userID,
		Name:              "User's Name",
		CreatedAt:         time.Now(),
		AuthnProviderType: pacta.EmailAndPass,
		AuthnProviderID:   pacta.UserID(email),
		Email:             email,
	}
	if diff := cmp.Diff(expected, actual, userCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	emailA := "userA@example.com"
	userID, err := tdb.CreateUser(tx, pacta.EmailAndPass, pacta.UserID(emailA), "User's Name", emailA)
	if err != nil {
		t.Fatalf("creating user: %v", err)
	}

	emailB := "userB@example.com"
	nameA := "Prince"
	err = tdb.UpdateUser(tx, userID, db.SetUserEmail(emailB), db.SetUserName(nameA))
	if err != nil {
		t.Fatalf("update user 1: %v", err)
	}

	nameB := "The artist formerly known as Prince"
	err = tdb.UpdateUser(tx, userID, db.SetUserName(nameB))
	if err != nil {
		t.Fatalf("update user 2: %v", err)
	}

	actual, err := tdb.User(tx, userID)
	if err != nil {
		t.Fatalf("getting user: %v", err)
	}
	expected := &pacta.User{
		ID:                userID,
		Name:              nameB,
		CreatedAt:         time.Now(),
		AuthnProviderType: pacta.EmailAndPass,
		AuthnProviderID:   pacta.UserID(emailA),
		Email:             emailB,
	}
	if diff := cmp.Diff(expected, actual, userCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}
}

func TestListUsers(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	emailA := "userA@example.com"
	emailB := "userB@example.com"
	emailC := "userC@example.com"
	fbIDB := "FB:UserB"
	googleIDC := "Google:UserC"
	nameA := "R2D2"
	nameB := "C3P0"
	userIDA, err0 := tdb.CreateUser(tx, pacta.EmailAndPass, pacta.UserID(emailA), "User A", emailA)
	userIDB, err1 := tdb.CreateUser(tx, pacta.Facebook, pacta.UserID(fbIDB), "User B", emailB)
	userIDC, err2 := tdb.CreateUser(tx, pacta.Google, pacta.UserID(googleIDC), "User C", emailC)
	err3 := tdb.UpdateUser(tx, userIDA, db.SetUserName(nameA))
	err4 := tdb.UpdateUser(tx, userIDB, db.SetUserName(nameB))
	noErrDuringSetup(t, err0, err1, err2, err3, err4)

	actual, err := tdb.Users(tx)
	if err != nil {
		t.Fatalf("listing users: %v", err)
	}
	expected := []*pacta.User{{
		ID:                userIDA,
		Name:              nameA,
		CreatedAt:         time.Now(),
		AuthnProviderType: pacta.EmailAndPass,
		AuthnProviderID:   pacta.UserID(emailA),
		Email:             emailA,
	}, {
		ID:                userIDB,
		Name:              nameB,
		CreatedAt:         time.Now(),
		AuthnProviderType: pacta.Facebook,
		AuthnProviderID:   pacta.UserID(fbIDB),
		Email:             emailB,
	}, {
		ID:                userIDC,
		Name:              "User C",
		CreatedAt:         time.Now(),
		AuthnProviderType: pacta.Google,
		AuthnProviderID:   pacta.UserID(googleIDC),
		Email:             emailC,
	}}
	if diff := cmp.Diff(expected, actual, userCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}
}

func userCmpOpts() cmp.Option {
	userIDLessFn := func(a, b pacta.UserID) bool {
		return a < b
	}
	groupLessFn := func(a, b *pacta.User) bool {
		return a.ID < b.ID
	}
	return cmp.Options{
		cmpopts.EquateEmpty(),
		cmpopts.EquateApproxTime(time.Second),
		cmpopts.SortSlices(userIDLessFn),
		cmpopts.SortSlices(groupLessFn),
	}
}

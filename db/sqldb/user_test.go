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

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	u := &pacta.User{
		AuthnMechanism: pacta.AuthnMechanism_EmailAndPass,
		AuthnID:        "authn-id",
		EnteredEmail:   "entered-email",
		CanonicalEmail: "canonical-email",
		Name:           "User's Name",
	}
	userID, err := tdb.CreateUser(tx, u)
	if err != nil {
		t.Fatalf("creating user: %v", err)
	}
	u.CreatedAt = time.Now()
	u.ID = userID

	// Read By ID
	actual, err := tdb.User(tx, userID)
	if err != nil {
		t.Fatalf("getting user: %v", err)
	}
	if diff := cmp.Diff(u, actual, userCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}

	// Read by Authn
	actual, err = tdb.UserByAuthn(tx, u.AuthnMechanism, u.AuthnID)
	if err != nil {
		t.Fatalf("getting user by authn: %w", err)
	}
	if diff := cmp.Diff(u, actual, userCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}

	// Read by id list
	aMap, err := tdb.Users(tx, []pacta.UserID{"somenonsense", userID})
	if err != nil {
		t.Fatalf("getting users: %w", err)
	}
	eMap := map[pacta.UserID]*pacta.User{userID: u}
	if diff := cmp.Diff(eMap, aMap, userCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}

	// Create should fail with the same authn
	u.ID = ""
	u2 := u.Clone()
	u2.EnteredEmail = "entered email 2"
	u2.CanonicalEmail = "canonical email 2"
	_, err = tdb.CreateUser(tx, u2)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Create should fail with the same canonicalEmail
	u3 := u.Clone()
	u3.EnteredEmail = "entered email 3"
	u3.AuthnID = "AUthn id 3"
	_, err = tdb.CreateUser(tx, u3)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Create should fail with the same entered email
	u4 := u.Clone()
	u4.AuthnID = "authn id 3"
	u4.CanonicalEmail = "canonical email 4"
	_, err = tdb.CreateUser(tx, u4)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Create should succeed if each are different.
	u5 := u.Clone()
	u5.EnteredEmail = "entered email 5"
	u5.AuthnID = "AUthn id 5"
	u5.CanonicalEmail = "canonical email 5"
	_, err = tdb.CreateUser(tx, u5)
	if err != nil {
		t.Fatal("expected success but got: %w", err)
	}
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	u := &pacta.User{
		AuthnMechanism: pacta.AuthnMechanism_EmailAndPass,
		AuthnID:        "authn-id",
		EnteredEmail:   "entered-email",
		CanonicalEmail: "canonical-email",
		Name:           "User's Name",
	}
	userID, err0 := tdb.CreateUser(tx, u)
	noErrDuringSetup(t, err0)
	u.CreatedAt = time.Now()
	u.ID = userID

	nameA := "Prince"
	lang := pacta.Language_DE
	err := tdb.UpdateUser(tx, userID, db.SetUserAdmin(true), db.SetUserName(nameA), db.SetUserPreferredLanguage(lang))
	if err != nil {
		t.Fatalf("update user 1: %v", err)
	}
	actual, err := tdb.User(tx, userID)
	if err != nil {
		t.Fatalf("getting user: %v", err)
	}
	u.Admin = true
	u.Name = nameA
	u.PreferredLanguage = lang
	if diff := cmp.Diff(u, actual, userCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}

	nameB := "The artist formerly known as Prince"
	err = tdb.UpdateUser(tx, userID, db.SetUserName(nameB), db.SetUserSuperAdmin(true), db.SetUserAdmin(false))
	if err != nil {
		t.Fatalf("update user 2: %v", err)
	}

	actual, err = tdb.User(tx, userID)
	if err != nil {
		t.Fatalf("getting user: %v", err)
	}
	u.Name = nameB
	u.SuperAdmin = true
	u.Admin = false
	if diff := cmp.Diff(u, actual, userCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}
}

func TestListUsers(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	nameA := "R2D2"
	nameB := "C3P0"
	userA := &pacta.User{
		Name:           "original name",
		AuthnMechanism: pacta.AuthnMechanism_EmailAndPass,
		AuthnID:        "AAA",
		CanonicalEmail: "canon",
		EnteredEmail:   "enterentered1",
	}
	userB := &pacta.User{
		Name:           "name b original",
		AuthnMechanism: pacta.AuthnMechanism_EmailAndPass,
		AuthnID:        "BBB",
		CanonicalEmail: "cnanon",
		EnteredEmail:   "entered2",
	}
	userC := &pacta.User{
		Name:           "User C",
		AuthnMechanism: pacta.AuthnMechanism_EmailAndPass,
		AuthnID:        "CCC",
		CanonicalEmail: "cannnnon",
		EnteredEmail:   "enter3",
	}
	userIDA, err0 := tdb.CreateUser(tx, userA)
	userA.ID = userIDA
	userA.CreatedAt = time.Now()
	userIDB, err1 := tdb.CreateUser(tx, userB)
	userB.ID = userIDB
	userB.CreatedAt = time.Now()
	userIDC, err2 := tdb.CreateUser(tx, userC)
	userC.ID = userIDC
	userC.CreatedAt = time.Now()
	err3 := tdb.UpdateUser(tx, userIDA, db.SetUserName(nameA))
	userA.Name = nameA
	err4 := tdb.UpdateUser(tx, userIDB, db.SetUserName(nameB))
	userB.Name = nameB
	noErrDuringSetup(t, err0, err1, err2, err3, err4)

	actual, err := tdb.Users(tx, []pacta.UserID{userIDA, userIDB, userIDC, "some nonsense"})
	if err != nil {
		t.Fatalf("listing users: %v", err)
	}
	expected := map[pacta.UserID]*pacta.User{
		userIDA: userA,
		userIDB: userB,
		userIDC: userC,
	}
	if diff := cmp.Diff(expected, actual, userCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}
}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	u := &pacta.User{
		AuthnMechanism: pacta.AuthnMechanism_EmailAndPass,
		AuthnID:        "authn-id",
		EnteredEmail:   "entered-email",
		CanonicalEmail: "canonical-email",
		Name:           "User's Name",
	}
	userID, err0 := tdb.CreateUser(tx, u)
	noErrDuringSetup(t, err0)

	err := tdb.DeleteUser(tx, userID)
	if err != nil {
		t.Fatalf("deleting user: %v", err)
	}

	// Read By ID
	_, err = tdb.User(tx, userID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Read by Authn
	_, err = tdb.UserByAuthn(tx, u.AuthnMechanism, u.AuthnID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Read by id list
	aMap, err := tdb.Users(tx, []pacta.UserID{"somenonsense", userID, "something else"})
	if err != nil {
		t.Fatalf("getting users: %w", err)
	}
	eMap := map[pacta.UserID]*pacta.User{}
	if diff := cmp.Diff(eMap, aMap, userCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}
}

func userCmpOpts() cmp.Option {
	userIDLessFn := func(a, b pacta.UserID) bool {
		return a < b
	}
	userLessFn := func(a, b *pacta.User) bool {
		return a.ID < b.ID
	}
	return cmp.Options{
		cmpopts.EquateEmpty(),
		cmpopts.EquateApproxTime(time.Second),
		cmpopts.SortSlices(userIDLessFn),
		cmpopts.SortSlices(userLessFn),
		cmpopts.SortMaps(userIDLessFn),
	}
}

func userForTesting(t *testing.T, tdb *DB) *pacta.User {
	t.Helper()
	return userForTestingWithKey(t, tdb, "only")
}

func userForTestingWithKey(t *testing.T, tdb *DB, key string) *pacta.User {
	t.Helper()
	u := &pacta.User{
		CanonicalEmail: fmt.Sprintf("canoncal-email-%s@example.com", key),
		EnteredEmail:   fmt.Sprintf("entered-email-%s+helloworld@example.com", key),
		AuthnMechanism: pacta.AuthnMechanism_EmailAndPass,
		AuthnID:        fmt.Sprintf("authn-id-%s", key),
	}
	ctx := context.Background()
	tx := tdb.NoTxn(ctx)
	uid, err := tdb.CreateUser(tx, u)
	if err != nil {
		t.Fatalf("creating user: %v", err)
	}
	u.ID = uid
	return u
}

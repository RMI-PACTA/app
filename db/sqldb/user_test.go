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
	userID, err := tdb.createUser(tx, u)
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
		t.Fatalf("getting user by authn: %v", err)
	}
	if diff := cmp.Diff(u, actual, userCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}

	// Read by id list
	aMap, err := tdb.Users(tx, []pacta.UserID{"somenonsense", userID})
	if err != nil {
		t.Fatalf("getting users: %v", err)
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
	_, err = tdb.createUser(tx, u2)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Create should fail with the same canonicalEmail
	u3 := u.Clone()
	u3.EnteredEmail = "entered email 3"
	u3.AuthnID = "AUthn id 3"
	_, err = tdb.createUser(tx, u3)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Create should fail with the same entered email
	u4 := u.Clone()
	u4.AuthnID = "authn id 3"
	u4.CanonicalEmail = "canonical email 4"
	_, err = tdb.createUser(tx, u4)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Create should succeed if each are different.
	u5 := u.Clone()
	u5.EnteredEmail = "entered email 5"
	u5.AuthnID = "AUthn id 5"
	u5.CanonicalEmail = "canonical email 5"
	_, err = tdb.createUser(tx, u5)
	if err != nil {
		t.Fatalf("expected success but got: %v", err)
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
	userID, err0 := tdb.createUser(tx, u)
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

func TestQueryUsers(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	userA := &pacta.User{
		Name:           "Assitant Regional Manager Schrute",
		AuthnMechanism: pacta.AuthnMechanism_EmailAndPass,
		AuthnID:        "AAA",
		CanonicalEmail: "dwight@dm.com",
		EnteredEmail:   "something-else",
	}
	userB := &pacta.User{
		Name:           "Jim",
		AuthnMechanism: pacta.AuthnMechanism_EmailAndPass,
		AuthnID:        "BBB",
		CanonicalEmail: "jim@dm.com",
		EnteredEmail:   "entered2",
	}
	userC := &pacta.User{
		Name:           "DWIGHT SCHRUTE, FARMER",
		AuthnMechanism: pacta.AuthnMechanism_EmailAndPass,
		AuthnID:        "CCC",
		CanonicalEmail: "beets-for-sale-northern-pa@gmail.com",
		EnteredEmail:   "entered3",
	}
	userIDA, err0 := tdb.createUser(tx, userA)
	userA.ID = userIDA
	userA.CreatedAt = time.Now()
	userIDB, err1 := tdb.createUser(tx, userB)
	userB.ID = userIDB
	userB.CreatedAt = time.Now()
	userIDC, err2 := tdb.createUser(tx, userC)
	userC.ID = userIDC
	userC.CreatedAt = time.Now()
	noErrDuringSetup(t, err0, err1, err2)

	testCases := []struct {
		name           string
		query          *db.UserQuery
		expected       []pacta.UserID
		expectedMore   bool
		expectedCursor string
	}{
		{
			name: "Sort Asc",
			query: &db.UserQuery{
				Sorts: []*db.UserQuerySort{{By: db.UserQuerySortBy_CreatedAt, Ascending: true}},
				Limit: 4,
			},
			expected:     []pacta.UserID{userIDA, userIDB, userIDC},
			expectedMore: false,
		},
		{
			name: "Sort Desc",
			query: &db.UserQuery{
				Sorts: []*db.UserQuerySort{{By: db.UserQuerySortBy_CreatedAt, Ascending: false}},
				Limit: 4,
			},
			expected:     []pacta.UserID{userIDC, userIDB, userIDA},
			expectedMore: false,
		},
		{
			name: "Limit Enforced",
			query: &db.UserQuery{
				Sorts: []*db.UserQuerySort{{By: db.UserQuerySortBy_CreatedAt, Ascending: false}},
				Limit: 2,
			},
			expected:       []pacta.UserID{userIDC, userIDB},
			expectedMore:   true,
			expectedCursor: "2",
		},
		{
			name: "With Cursor",
			query: &db.UserQuery{
				Sorts:  []*db.UserQuerySort{{By: db.UserQuerySortBy_CreatedAt, Ascending: false}},
				Limit:  2,
				Cursor: "2",
			},
			expected:     []pacta.UserID{userIDA},
			expectedMore: false,
		},
		{
			name: "Dwight LC",
			query: &db.UserQuery{
				Sorts:  []*db.UserQuerySort{{By: db.UserQuerySortBy_CreatedAt, Ascending: true}},
				Limit:  4,
				Wheres: []*db.UserQueryWhere{{NameOrEmailLike: "dwight"}},
			},
			expected:     []pacta.UserID{userIDA, userIDC},
			expectedMore: false,
		},
		{
			name: "Schrute",
			query: &db.UserQuery{
				Sorts:  []*db.UserQuerySort{{By: db.UserQuerySortBy_CreatedAt, Ascending: true}},
				Limit:  4,
				Wheres: []*db.UserQueryWhere{{NameOrEmailLike: "schrute"}},
			},
			expected:     []pacta.UserID{userIDA, userIDC},
			expectedMore: false,
		},
		{
			name: "Dwight Partial",
			query: &db.UserQuery{
				Sorts:  []*db.UserQuerySort{{By: db.UserQuerySortBy_CreatedAt, Ascending: true}},
				Limit:  4,
				Wheres: []*db.UserQueryWhere{{NameOrEmailLike: "wigh"}},
			},
			expected:     []pacta.UserID{userIDA, userIDC},
			expectedMore: false,
		},
		{
			name: "Dwight Spongebob",
			query: &db.UserQuery{
				Sorts:  []*db.UserQuerySort{{By: db.UserQuerySortBy_CreatedAt, Ascending: true}},
				Limit:  4,
				Wheres: []*db.UserQueryWhere{{NameOrEmailLike: "dWiGhT"}},
			},
			expected:     []pacta.UserID{userIDA, userIDC},
			expectedMore: false,
		},
		{
			name: "Dunder Miflin",
			query: &db.UserQuery{
				Sorts:  []*db.UserQuerySort{{By: db.UserQuerySortBy_CreatedAt, Ascending: true}},
				Limit:  4,
				Wheres: []*db.UserQueryWhere{{NameOrEmailLike: "dm.com"}},
			},
			expected:     []pacta.UserID{userIDA, userIDB},
			expectedMore: false,
		},
		{
			name: "Jim",
			query: &db.UserQuery{
				Limit:  4,
				Wheres: []*db.UserQueryWhere{{NameOrEmailLike: "jim"}},
			},
			expected:     []pacta.UserID{userIDB},
			expectedMore: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			users, pi, err := tdb.QueryUsers(nil, tc.query)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			userIDs := make([]pacta.UserID, len(users))
			for i, user := range users {
				userIDs[i] = user.ID
			}
			if diff := cmp.Diff(tc.expected, userIDs); diff != "" {
				t.Errorf("Expected and actual users do not match. Expected: %v, Actual: %v:\n%s\nDecodingMap = %+v", tc.expected, userIDs, diff, map[pacta.UserID]string{userIDA: "Dwight DM", userIDB: "Jim DM", userIDC: "Dwight Personal"})
			}
			if pi.HasNextPage != tc.expectedMore {
				t.Errorf("Expected HasNextPage to be %v, got %v", tc.expectedMore, pi.HasNextPage)
			}
			if tc.expectedCursor != "" && string(pi.Cursor) != tc.expectedCursor {
				t.Errorf("Expected cursor to be %v, got %v", tc.expectedCursor, pi.Cursor)
			}
		})
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
	userIDA, err0 := tdb.createUser(tx, userA)
	userA.ID = userIDA
	userA.CreatedAt = time.Now()
	userIDB, err1 := tdb.createUser(tx, userB)
	userB.ID = userIDB
	userB.CreatedAt = time.Now()
	userIDC, err2 := tdb.createUser(tx, userC)
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
	userID, err0 := tdb.createUser(tx, u)
	noErrDuringSetup(t, err0)

	_, err := tdb.DeleteUser(tx, userID)
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
		t.Fatalf("getting users: %v", err)
	}
	eMap := map[pacta.UserID]*pacta.User{}
	if diff := cmp.Diff(eMap, aMap, userCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}
}

func TestAuthnMechanismPersistability(t *testing.T) {
	testUserEnumConvertability(
		t,
		func(e pacta.AuthnMechanism, u *pacta.User) { u.AuthnMechanism = e },
		func(u *pacta.User) pacta.AuthnMechanism { return u.AuthnMechanism },
		pacta.AuthnMechanismValues,
	)
}

func TestLanguagePersistability(t *testing.T) {
	testUserEnumConvertability(
		t,
		func(e pacta.Language, u *pacta.User) { u.PreferredLanguage = e },
		func(u *pacta.User) pacta.Language { return u.PreferredLanguage },
		pacta.LanguageValues,
	)
}

func testUserEnumConvertability[E comparable](t *testing.T, writeE func(E, *pacta.User), readE func(*pacta.User) E, values []E) {
	var zeroValue E
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	base := &pacta.User{
		AuthnMechanism: pacta.AuthnMechanism_EmailAndPass,
		AuthnID:        "authn-id",
		EnteredEmail:   "entered-email",
		CanonicalEmail: "canonical-email",
		Name:           "User's Name",
	}
	var id pacta.UserID
	iteration := 0

	write := func(e E) error {
		u := base.Clone()
		u.EnteredEmail = fmt.Sprintf("entered-email-%d", iteration)
		u.AuthnID = fmt.Sprintf("authn-id-%d", iteration)
		u.CanonicalEmail = fmt.Sprintf("canonical-email-%d", iteration)
		writeE(e, u)
		iteration++
		id2, err := tdb.createUser(tx, u)
		id = id2
		return err
	}
	read := func() (E, error) {
		u, err := tdb.User(tx, id)
		if err != nil {
			return zeroValue, err
		}
		return readE(u), nil
	}

	testEnumConvertability(t, write, read, values)
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
	canonicalEmail := fmt.Sprintf("canoncal-email-%s@example.com", key)
	enteredEmail := fmt.Sprintf("entered-email-%s+helloworld@example.com", key)
	authnMechanism := pacta.AuthnMechanism_EmailAndPass
	authnID := fmt.Sprintf("authn-id-%s", key)
	ctx := context.Background()
	tx := tdb.NoTxn(ctx)
	user, err := tdb.GetOrCreateUserByAuthn(tx, authnMechanism, authnID, enteredEmail, canonicalEmail)
	if err != nil {
		t.Fatalf("creating user: %v", err)
	}
	return user
}

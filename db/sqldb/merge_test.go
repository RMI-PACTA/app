package sqldb

import (
	"context"
	"fmt"
	"testing"

	"github.com/RMI/pacta/pacta"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestMergeUsers(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	a := userForTestingWithKey(t, tdb, "A")
	b := userForTestingWithKey(t, tdb, "B")
	c := userForTestingWithKey(t, tdb, "C")
	d := userForTestingWithKey(t, tdb, "D")
	e := userForTestingWithKey(t, tdb, "E")
	f := userForTestingWithKey(t, tdb, "F")

	usersToIDs := func(users []*pacta.User) []string {
		ids := []string{}
		for _, u := range users {
			ids = append(ids, string(u.ID))
		}
		return ids
	}
	cmpOpts := cmpopts.SortSlices(func(a, b string) bool { return a < b })
	runTests := func(name string, tests []struct {
		in       []*pacta.User
		expected []*pacta.User
	}) {
		t.Helper()
		for _, test := range tests {
			t.Run(fmt.Sprintf("%s_case_%v", name, test.in), func(t *testing.T) {
				got, err := tdb.findAllMergedUsers(tx, usersToIDs(test.in))
				if err != nil {
					t.Fatalf("get users: %v", err)
				}
				if diff := cmp.Diff(usersToIDs(test.expected), got, cmpOpts); diff != "" {
					t.Fatalf("users mismatch (-want +got):\n%s", diff)
				}
			})

		}
	}

	runTests("Pre-Merge", []struct {
		in       []*pacta.User
		expected []*pacta.User
	}{
		{
			in:       []*pacta.User{},
			expected: []*pacta.User{},
		}, {
			in:       []*pacta.User{a},
			expected: []*pacta.User{a},
		}, {
			in:       []*pacta.User{c, c, c},
			expected: []*pacta.User{c},
		}, {
			in:       []*pacta.User{e},
			expected: []*pacta.User{e},
		}, {
			in:       []*pacta.User{a, b},
			expected: []*pacta.User{a, b},
		}, {
			in:       []*pacta.User{a, b, c, d, e, f},
			expected: []*pacta.User{a, b, c, d, e, f},
		},
	})

	// Merge A and B
	if err := tdb.RecordUserMerge(tx, a.ID, b.ID, "AdminID"); err != nil {
		t.Fatalf("record merge: %v", err)
	}
	// Merge C and D
	if err := tdb.RecordUserMerge(tx, c.ID, d.ID, "AdminID"); err != nil {
		t.Fatalf("record merge: %v", err)
	}

	runTests("Post-First-Merge", []struct {
		in       []*pacta.User
		expected []*pacta.User
	}{
		{
			in:       []*pacta.User{},
			expected: []*pacta.User{},
		}, {
			in:       []*pacta.User{a},
			expected: []*pacta.User{a, b},
		}, {
			in:       []*pacta.User{e},
			expected: []*pacta.User{e},
		}, {
			in:       []*pacta.User{c, c, c},
			expected: []*pacta.User{c, d},
		}, {
			in:       []*pacta.User{a, b},
			expected: []*pacta.User{a, b},
		}, {
			in:       []*pacta.User{a, b, c, d, e, f},
			expected: []*pacta.User{a, b, c, d, e, f},
		},
	})

	// Merge B and E
	if err := tdb.RecordUserMerge(tx, b.ID, e.ID, "AdminID"); err != nil {
		t.Fatalf("record merge: %v", err)
	}
	// Merge B and F
	if err := tdb.RecordUserMerge(tx, f.ID, b.ID, "AdminID"); err != nil {
		t.Fatalf("record merge: %v", err)
	}

	runTests("Post-Second-Merge", []struct {
		in       []*pacta.User
		expected []*pacta.User
	}{
		{
			in:       []*pacta.User{},
			expected: []*pacta.User{},
		}, {
			in:       []*pacta.User{a},
			expected: []*pacta.User{a, b, e, f},
		}, {
			in:       []*pacta.User{e},
			expected: []*pacta.User{a, b, e, f},
		}, {
			in:       []*pacta.User{c, c, c},
			expected: []*pacta.User{c, d},
		}, {
			in:       []*pacta.User{a, b},
			expected: []*pacta.User{a, b, e, f},
		}, {
			in:       []*pacta.User{a, b, c, d, e, f},
			expected: []*pacta.User{a, b, c, d, e, f},
		},
	})

	// Merge D and E
	if err := tdb.RecordUserMerge(tx, d.ID, e.ID, "AdminID"); err != nil {
		t.Fatalf("record merge: %v", err)
	}

	runTests("Post-Third-Merge", []struct {
		in       []*pacta.User
		expected []*pacta.User
	}{
		{
			in:       []*pacta.User{},
			expected: []*pacta.User{},
		}, {
			in:       []*pacta.User{a},
			expected: []*pacta.User{a, b, c, d, e, f},
		}, {
			in:       []*pacta.User{e},
			expected: []*pacta.User{a, b, c, d, e, f},
		}, {
			in:       []*pacta.User{c, c, c},
			expected: []*pacta.User{a, b, c, d, e, f},
		}, {
			in:       []*pacta.User{a, b},
			expected: []*pacta.User{a, b, c, d, e, f},
		}, {
			in:       []*pacta.User{a, b, c, d, e, f},
			expected: []*pacta.User{a, b, c, d, e, f},
		},
	})
}

func TestMergeOwners(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	uA := userForTestingWithKey(t, tdb, "A")
	uB := userForTestingWithKey(t, tdb, "B")
	uC := userForTestingWithKey(t, tdb, "C")
	uD := userForTestingWithKey(t, tdb, "D")
	uE := userForTestingWithKey(t, tdb, "E")
	uF := userForTestingWithKey(t, tdb, "F")
	a, err0 := tdb.GetOwnerForUser(tx, uA.ID)
	b, err1 := tdb.GetOwnerForUser(tx, uB.ID)
	c, err2 := tdb.GetOwnerForUser(tx, uC.ID)
	d, err3 := tdb.GetOwnerForUser(tx, uD.ID)
	e, err4 := tdb.GetOwnerForUser(tx, uE.ID)
	f, err5 := tdb.GetOwnerForUser(tx, uF.ID)
	noErrDuringSetup(t, err0, err1, err2, err3, err4, err5)

	cmpOpts := cmpopts.SortSlices(func(a, b pacta.OwnerID) bool { return a < b })
	runTests := func(name string, tests []struct {
		in       []pacta.OwnerID
		expected []pacta.OwnerID
	}) {
		t.Helper()
		for _, test := range tests {
			t.Run(fmt.Sprintf("%s_case_%v", name, test.in), func(t *testing.T) {
				got, err := tdb.findAllMergedOwners(tx, test.in)
				if err != nil {
					t.Fatalf("get owners: %v", err)
				}
				if diff := cmp.Diff(test.expected, got, cmpOpts); diff != "" {
					t.Fatalf("owners mismatch (-want +got):\n%s", diff)
				}
			})

		}
	}

	runTests("Pre-Merge", []struct {
		in       []pacta.OwnerID
		expected []pacta.OwnerID
	}{
		{
			in:       []pacta.OwnerID{},
			expected: []pacta.OwnerID{},
		}, {
			in:       []pacta.OwnerID{a},
			expected: []pacta.OwnerID{a},
		}, {
			in:       []pacta.OwnerID{c, c, c},
			expected: []pacta.OwnerID{c},
		}, {
			in:       []pacta.OwnerID{e},
			expected: []pacta.OwnerID{e},
		}, {
			in:       []pacta.OwnerID{a, b},
			expected: []pacta.OwnerID{a, b},
		}, {
			in:       []pacta.OwnerID{a, b, c, d, e, f},
			expected: []pacta.OwnerID{a, b, c, d, e, f},
		},
	})

	// Merge A and B
	if err := tdb.RecordOwnerMerge(tx, a, b, "AdminID"); err != nil {
		t.Fatalf("record merge: %v", err)
	}
	// Merge C and D
	if err := tdb.RecordOwnerMerge(tx, c, d, "AdminID"); err != nil {
		t.Fatalf("record merge: %v", err)
	}

	runTests("Post-First-Merge", []struct {
		in       []pacta.OwnerID
		expected []pacta.OwnerID
	}{
		{
			in:       []pacta.OwnerID{},
			expected: []pacta.OwnerID{},
		}, {
			in:       []pacta.OwnerID{a},
			expected: []pacta.OwnerID{a, b},
		}, {
			in:       []pacta.OwnerID{e},
			expected: []pacta.OwnerID{e},
		}, {
			in:       []pacta.OwnerID{c, c, c},
			expected: []pacta.OwnerID{c, d},
		}, {
			in:       []pacta.OwnerID{a, b},
			expected: []pacta.OwnerID{a, b},
		}, {
			in:       []pacta.OwnerID{a, b, c, d, e, f},
			expected: []pacta.OwnerID{a, b, c, d, e, f},
		},
	})

	// Merge B and E
	if err := tdb.RecordOwnerMerge(tx, b, e, "AdminID"); err != nil {
		t.Fatalf("record merge: %v", err)
	}
	// Merge B and F
	if err := tdb.RecordOwnerMerge(tx, f, b, "AdminID"); err != nil {
		t.Fatalf("record merge: %v", err)
	}

	runTests("Post-Second-Merge", []struct {
		in       []pacta.OwnerID
		expected []pacta.OwnerID
	}{
		{
			in:       []pacta.OwnerID{},
			expected: []pacta.OwnerID{},
		}, {
			in:       []pacta.OwnerID{a},
			expected: []pacta.OwnerID{a, b, e, f},
		}, {
			in:       []pacta.OwnerID{e},
			expected: []pacta.OwnerID{a, b, e, f},
		}, {
			in:       []pacta.OwnerID{c, c, c},
			expected: []pacta.OwnerID{c, d},
		}, {
			in:       []pacta.OwnerID{a, b},
			expected: []pacta.OwnerID{a, b, e, f},
		}, {
			in:       []pacta.OwnerID{a, b, c, d, e, f},
			expected: []pacta.OwnerID{a, b, c, d, e, f},
		},
	})

	// Merge D and E
	if err := tdb.RecordOwnerMerge(tx, d, e, "AdminID"); err != nil {
		t.Fatalf("record merge: %v", err)
	}

	runTests("Post-Third-Merge", []struct {
		in       []pacta.OwnerID
		expected []pacta.OwnerID
	}{
		{
			in:       []pacta.OwnerID{},
			expected: []pacta.OwnerID{},
		}, {
			in:       []pacta.OwnerID{a},
			expected: []pacta.OwnerID{a, b, c, d, e, f},
		}, {
			in:       []pacta.OwnerID{e},
			expected: []pacta.OwnerID{a, b, c, d, e, f},
		}, {
			in:       []pacta.OwnerID{c, c, c},
			expected: []pacta.OwnerID{a, b, c, d, e, f},
		}, {
			in:       []pacta.OwnerID{a, b},
			expected: []pacta.OwnerID{a, b, c, d, e, f},
		}, {
			in:       []pacta.OwnerID{a, b, c, d, e, f},
			expected: []pacta.OwnerID{a, b, c, d, e, f},
		},
	})
}

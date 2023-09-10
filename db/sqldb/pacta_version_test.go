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

func TestCreatePACTAVersion(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	pv := &pacta.PACTAVersion{
		Name:        "pacta version",
		Description: "pacta description",
		Digest:      "pacta digest",
	}
	pvid, err := tdb.CreatePACTAVersion(tx, pv)
	if err != nil {
		t.Fatalf("creating pacta_version: %v", err)
	}
	pv.CreatedAt = time.Now()
	pv.ID = pvid

	// Read By ID
	actual, err := tdb.PACTAVersion(tx, pvid)
	if err != nil {
		t.Fatalf("getting pacta version: %v", err)
	}
	if diff := cmp.Diff(pv, actual, pactaVersionCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}

	// Create should fail with the same digest
	pv.ID = ""
	pv.CreatedAt = time.Time{}
	pv2 := pv.Clone()
	_, err = tdb.CreatePACTAVersion(tx, pv2)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Create should succeed with a different digest
	pv3 := pv.Clone()
	pv3.Digest = "different digest"
	_, err = tdb.CreatePACTAVersion(tx, pv3)
	if err != nil {
		t.Fatalf("creating pacta_version: %v", err)
	}
}

func TestUpdatePACTAVersion(t *testing.T) {
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
	pv.CreatedAt = time.Now()
	pv.ID = pvID

	name := "New Name"
	desc := "New Description"
	digest := "New Digest"
	err := tdb.UpdatePACTAVersion(tx, pvID, db.SetPACTAVersionDescription(desc), db.SetPACTAVersionDigest(digest), db.SetPACTAVersionName(name))
	if err != nil {
		t.Fatalf("update pacta version: %v", err)
	}
	actual, err := tdb.PACTAVersion(tx, pvID)
	if err != nil {
		t.Fatalf("getting pacta version: %v", err)
	}
	pv.Name = name
	pv.Digest = digest
	pv.Description = desc
	if diff := cmp.Diff(pv, actual, pactaVersionCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}
}

func TestDefaultPACTAVersion(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	pvA := &pacta.PACTAVersion{
		Name:        "pv 1",
		Digest:      "111",
		Description: "Pacta Version 1",
	}
	pvB := &pacta.PACTAVersion{
		Name:        "pv 2",
		Digest:      "222",
		Description: "Pacta Version 2",
	}
	pvC := &pacta.PACTAVersion{
		Name:        "pv 3",
		Digest:      "333",
		Description: "Pacta Version 3",
	}
	pvIDA, err0 := tdb.CreatePACTAVersion(tx, pvA)
	pvA.ID = pvIDA
	pvA.CreatedAt = time.Now()
	pvIDB, err1 := tdb.CreatePACTAVersion(tx, pvB)
	pvB.ID = pvIDB
	pvB.CreatedAt = time.Now()
	pvIDC, err2 := tdb.CreatePACTAVersion(tx, pvC)
	pvC.ID = pvIDC
	pvC.CreatedAt = time.Now()
	noErrDuringSetup(t, err0, err1, err2)

	d, err := tdb.DefaultPACTAVersion(tx)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	actual, err := tdb.PACTAVersions(tx)
	if err != nil {
		t.Fatalf("listing pacta versions: %v", err)
	}
	expected := []*pacta.PACTAVersion{pvA, pvB, pvC}
	if diff := cmp.Diff(expected, actual, pactaVersionCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}

	err = tdb.SetDefaultPACTAVersion(tx, pvB.ID)
	if err != nil {
		t.Fatalf("setting primary pacta version: %v", err)
	}
	pvB.IsDefault = true
	actual, err = tdb.PACTAVersions(tx)
	if diff := cmp.Diff(expected, actual, pactaVersionCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}
	d, err = tdb.DefaultPACTAVersion(tx)
	if err != nil {
		t.Fatalf("getting default pacta version: %v", err)
	}
	if diff := cmp.Diff(pvB, d, pactaVersionCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}

	err = tdb.SetDefaultPACTAVersion(tx, pvC.ID)
	if err != nil {
		t.Fatalf("setting primary pacta version: %v", err)
	}
	pvB.IsDefault = false
	pvC.IsDefault = true
	actual, err = tdb.PACTAVersions(tx)
	if diff := cmp.Diff(expected, actual, pactaVersionCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}
	d, err = tdb.DefaultPACTAVersion(tx)
	if err != nil {
		t.Fatalf("getting default pacta version: %v", err)
	}
	if diff := cmp.Diff(pvC, d, pactaVersionCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}
}

func TestPACTAVersionDeletion(t *testing.T) {
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
	pv.CreatedAt = time.Now()
	pv.ID = pvID

	err := tdb.DeletePACTAVersion(tx, pvID)
	if err != nil {
		t.Fatalf("deleting pacta version: %v", err)
	}

	// Read By ID
	_, err = tdb.PACTAVersion(tx, pvID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Read by id list
	m, err := tdb.PACTAVersions(tx)
	if err != nil {
		t.Fatalf("getting pacta_versions: %w", err)
	}
	if diff := cmp.Diff(m, []*pacta.PACTAVersion{}, pactaVersionCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}
}

func pactaVersionCmpOpts() cmp.Option {
	pactaVersionIDLessFn := func(a, b pacta.PACTAVersionID) bool {
		return a < b
	}
	pactaVersionLessFn := func(a, b *pacta.PACTAVersion) bool {
		return a.ID < b.ID
	}
	return cmp.Options{
		cmpopts.SortSlices(pactaVersionLessFn),
		cmpopts.SortMaps(pactaVersionIDLessFn),
		cmpopts.EquateEmpty(),
		cmpopts.EquateApproxTime(time.Second),
	}
}

func pactaVersionForTesting(t *testing.T, tdb *DB) *pacta.PACTAVersion {
	t.Helper()
	return pactaVersionForTestingWithKey(t, tdb, "only")
}

func pactaVersionForTestingWithKey(t *testing.T, tdb *DB, key string) *pacta.PACTAVersion {
	t.Helper()
	pv := &pacta.PACTAVersion{
		Name:        "pacta version",
		Description: "pacta description",
		Digest:      fmt.Sprintf("pacta digest %s", key),
	}
	tx := tdb.NoTxn(context.Background())
	pvID, err := tdb.CreatePACTAVersion(tx, pv)
	if err != nil {
		t.Fatalf("creating pacta_version: %v", err)
	}
	pv.ID = pvID
	return pv
}

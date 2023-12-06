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

func TestIncompleteUploadCRUD(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	b := blobForTesting(t, tdb)
	u1 := userForTestingWithKey(t, tdb, "1")
	o1 := ownerUserForTesting(t, tdb, u1)
	u2 := userForTestingWithKey(t, tdb, "2")
	o2 := ownerUserForTesting(t, tdb, u2)
	cmpOpts := incompleteUploadCmpOpts()

	iu := &pacta.IncompleteUpload{
		Name:         "i-u-name",
		Description:  "i-u-description",
		HoldingsDate: exampleHoldingsDate,
		Owner:        &pacta.Owner{ID: o1.ID},
		Blob:         &pacta.Blob{ID: b.ID},
	}
	id, err := tdb.CreateIncompleteUpload(tx, iu)
	if err != nil {
		t.Fatalf("creating incomplete upload: %v", err)
	}
	iu.ID = id
	iu.CreatedAt = time.Now()

	actual, err := tdb.IncompleteUpload(tx, iu.ID)
	if err != nil {
		t.Fatalf("reading incomplete_upload: %v", err)
	}
	if diff := cmp.Diff(iu, actual, cmpOpts); diff != "" {
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}

	ius, err := tdb.IncompleteUploads(tx, []pacta.IncompleteUploadID{iu.ID, iu.ID, "nonsense"})
	if err != nil {
		t.Fatalf("reading incomplete_uploads: %w", err)
	}
	if diff := cmp.Diff(map[pacta.IncompleteUploadID]*pacta.IncompleteUpload{iu.ID: iu}, ius, cmpOpts); diff != "" {
		t.Fatalf("incomplete_upload mismatch (-want +got):\n%s", diff)
	}

	nName := "new-name"
	nDesc := "new-description"
	ranAt := time.UnixMilli(111111111)
	completedAt := time.UnixMilli(222222222)
	failureCode := pacta.FailureCode_Unknown
	failureMessage := "failureMessage"
	hd := exampleHoldingsDate2
	err = tdb.UpdateIncompleteUpload(tx, iu.ID,
		db.SetIncompleteUploadName(nName),
		db.SetIncompleteUploadDescription(nDesc),
		db.SetIncompleteUploadRanAt(ranAt),
		db.SetIncompleteUploadOwner(o2.ID),
		db.SetIncompleteUploadFailureCode(failureCode),
		db.SetIncompleteUploadCompletedAt(completedAt),
		db.SetIncompleteUploadAdminDebugEnabled(true),
		db.SetIncompleteUploadFailureMessage(failureMessage),
		db.SetIncompleteUploadHoldingsDate(hd),
	)
	if err != nil {
		t.Fatalf("updating incomplete upload: %v", err)
	}
	iu.Name = nName
	iu.Description = nDesc
	iu.RanAt = ranAt
	iu.Owner = &pacta.Owner{ID: o2.ID}
	iu.FailureCode = failureCode
	iu.CompletedAt = completedAt
	iu.AdminDebugEnabled = true
	iu.FailureMessage = failureMessage
	iu.HoldingsDate = hd

	actual, err = tdb.IncompleteUpload(tx, iu.ID)
	if err != nil {
		t.Fatalf("reading incomplete upload: %v", err)
	}
	if diff := cmp.Diff(iu, actual, cmpOpts); diff != "" {
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}

	iusbyo, err := tdb.IncompleteUploadsByOwner(tx, o2.ID)
	if err != nil {
		t.Fatalf("reading incomplete_uploads by owner: %v", err)
	}
	if diff := cmp.Diff([]*pacta.IncompleteUpload{iu}, iusbyo, cmpOpts); diff != "" {
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}
	iusbyo, err = tdb.IncompleteUploadsByOwner(tx, o1.ID)
	if err != nil {
		t.Fatalf("reading incomplete_uploads by owner: %v", err)
	}
	if diff := cmp.Diff([]*pacta.IncompleteUpload{}, iusbyo, cmpOpts); diff != "" {
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}

	buris, err := tdb.DeleteIncompleteUpload(tx, iu.ID)
	if err != nil {
		t.Fatalf("deleting incompleteUpload: %v", err)
	}
	if diff := cmp.Diff(b.BlobURI, buris); diff != "" {
		t.Fatalf("blob uri mismatch (-want +got):\n%s", diff)
	}
}

func TestFailureCodePersistability(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	b := blobForTesting(t, tdb)
	u := userForTestingWithKey(t, tdb, "only")
	o := ownerUserForTesting(t, tdb, u)

	iu := &pacta.IncompleteUpload{
		Name:         "i-u-name",
		Description:  "i-u-description",
		HoldingsDate: exampleHoldingsDate,
		Owner:        &pacta.Owner{ID: o.ID},
		Blob:         &pacta.Blob{ID: b.ID},
	}
	id, err := tdb.CreateIncompleteUpload(tx, iu)
	if err != nil {
		t.Fatalf("creating incomplete upload: %v", err)
	}
	iu.ID = id
	iu.CreatedAt = time.Now()

	write := func(ft pacta.FailureCode) error {
		return tdb.UpdateIncompleteUpload(tx, iu.ID, db.SetIncompleteUploadFailureCode(ft))
	}
	read := func() (pacta.FailureCode, error) {
		iu, err := tdb.IncompleteUpload(tx, id)
		if err != nil {
			return "", err
		}
		return iu.FailureCode, nil
	}

	testEnumConvertability(t, write, read, pacta.FailureCodeValues)
}

func incompleteUploadCmpOpts() cmp.Option {
	return cmp.Options{
		cmpopts.EquateEmpty(),
		cmpopts.EquateApproxTime(time.Second),
	}
}

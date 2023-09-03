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

func TestCreateBlob(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	b := &pacta.Blob{
		FileType: pacta.FileType_CSV,
		BlobURI:  pacta.BlobURI("blob-uri"),
		FileName: "example-spreadsheet",
	}
	bid, err := tdb.CreateBlob(tx, b)
	if err != nil {
		t.Fatalf("creating blob: %v", err)
	}
	b.CreatedAt = time.Now()
	b.ID = bid

	// Read By ID
	actual, err := tdb.Blob(tx, bid)
	if err != nil {
		t.Fatalf("getting blob: %v", err)
	}
	if diff := cmp.Diff(b, actual, blobCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}

	// Read by list
	m, err := tdb.Blobs(tx, []pacta.BlobID{bid, "something else"})
	if err != nil {
		t.Fatalf("getting blobs: %v", err)
	}
	if diff := cmp.Diff(m, map[pacta.BlobID]*pacta.Blob{bid: b}, blobCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}

	// Create should fail with the same uri
	b.ID = ""
	b.CreatedAt = time.Time{}
	b2 := b.Clone()
	_, err = tdb.CreateBlob(tx, b2)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Create should succeed with a different uri
	b3 := b.Clone()
	b3.BlobURI = "different blob URI"
	_, err = tdb.CreateBlob(tx, b3)
	if err != nil {
		t.Fatalf("creating blob: %v", err)
	}
}

func TestUpdateBlob(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	b := &pacta.Blob{
		FileType: pacta.FileType_CSV,
		BlobURI:  pacta.BlobURI("blob-uri"),
		FileName: "example-spreadsheet",
	}
	bID, err0 := tdb.CreateBlob(tx, b)
	noErrDuringSetup(t, err0)
	b.CreatedAt = time.Now()
	b.ID = bID

	name := "New Name"
	err := tdb.UpdateBlob(tx, bID, db.SetBlobFileName(name))
	if err != nil {
		t.Fatalf("update blob: %v", err)
	}
	actual, err := tdb.Blob(tx, bID)
	if err != nil {
		t.Fatalf("getting blob: %v", err)
	}
	b.FileName = name
	if diff := cmp.Diff(b, actual, blobCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}
}

func TestBlobDeletion(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	blobURI := pacta.BlobURI("blob-uri")
	b := &pacta.Blob{
		FileName: "file name",
		FileType: pacta.FileType_HTML,
		BlobURI:  blobURI,
	}
	bID, err0 := tdb.CreateBlob(tx, b)
	noErrDuringSetup(t, err0)

	returned, err := tdb.DeleteBlob(tx, bID)
	if err != nil {
		t.Fatalf("deleting blob: %v", err)
	}

	if returned != blobURI {
		t.Fatalf("expected blob URI %q, got %q", blobURI, returned)
	}

	// Read By ID
	_, err = tdb.Blob(tx, bID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Read by list
	m, err := tdb.Blobs(tx, []pacta.BlobID{bID, "something else"})
	if err != nil {
		t.Fatalf("getting blobs: %v", err)
	}
	if diff := cmp.Diff(m, map[pacta.BlobID]*pacta.Blob{}, blobCmpOpts()); diff != "" {
		t.Fatalf("unexpected diff (-want +got)\n%s", diff)
	}
}

func TestFileTypePersistability(t *testing.T) {
	ctx := context.Background()
	tdb := createDBForTesting(t)
	tx := tdb.NoTxn(ctx)
	base := &pacta.Blob{
		BlobURI:  "blob-uri",
		FileType: pacta.FileType_HTML,
		FileName: "file-name",
	}
	var id pacta.BlobID
	iteration := 0

	write := func(ft pacta.FileType) error {
		blob := base.Clone()
		blob.BlobURI = pacta.BlobURI(fmt.Sprintf("blob-uri-%d", iteration))
		blob.FileType = ft
		iteration++
		id2, err := tdb.CreateBlob(tx, blob)
		id = id2
		return err
	}
	read := func() (pacta.FileType, error) {
		b, err := tdb.Blob(tx, id)
		if err != nil {
			return "", err
		}
		return b.FileType, nil
	}

	testEnumConvertability(t, write, read, pacta.FileTypeValues)
}

func blobCmpOpts() cmp.Option {
	blobIDLessFn := func(a, b pacta.BlobID) bool {
		return a < b
	}
	blobLessFn := func(a, b *pacta.Blob) bool {
		return a.ID < b.ID
	}
	return cmp.Options{
		cmpopts.SortSlices(blobLessFn),
		cmpopts.SortMaps(blobIDLessFn),
		cmpopts.EquateEmpty(),
		cmpopts.EquateApproxTime(time.Second),
	}
}

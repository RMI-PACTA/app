package sqldb

import (
	"fmt"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
	"github.com/jackc/pgx/v5"
)

const blobIDNamespace = "blob"
const blobSelectColumns = `
	blob.id,
	blob.blob_uri,
	blob.file_type,
	blob.file_name,
	blob.created_at`

func (d *DB) Blob(tx db.Tx, id pacta.BlobID) (*pacta.Blob, error) {
	rows, err := d.query(tx, `
		SELECT `+blobSelectColumns+`
		FROM blob 
		WHERE id = $1;`, id)
	if err != nil {
		return nil, fmt.Errorf("querying blob: %w", err)
	}
	blobs, err := rowsToBlobs(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to blobs: %w", err)
	}
	return exactlyOne("blob", id, blobs)
}

func (d *DB) Blobs(tx db.Tx, ids []pacta.BlobID) (map[pacta.BlobID]*pacta.Blob, error) {
	rows, err := d.query(tx, `
		SELECT `+blobSelectColumns+`
		FROM blob 
		WHERE id IN `+createWhereInFmt(len(ids))+`;`, idsToInterface(ids)...)
	if err != nil {
		return nil, fmt.Errorf("querying blobs: %w", err)
	}
	blobs, err := rowsToBlobs(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to blobs: %w", err)
	}
	result := make(map[pacta.BlobID]*pacta.Blob, len(blobs))
	for _, blob := range blobs {
		result[blob.ID] = blob
	}
	return result, nil
}

func (d *DB) CreateBlob(tx db.Tx, b *pacta.Blob) (pacta.BlobID, error) {
	if err := validateBlobForCreation(b); err != nil {
		return "", fmt.Errorf("validating blob for creation: %w", err)
	}
	id := pacta.BlobID(d.randomID(blobIDNamespace))
	err := d.exec(tx, `
		INSERT INTO blob 
			(id, blob_uri, file_type, file_name)
			VALUES
			($1, $2, $3, $4);
	`, id, b.BlobURI, b.FileType, b.FileName)
	if err != nil {
		return "", fmt.Errorf("creating blob row: %w", err)
	}
	return id, nil
}

func (d *DB) UpdateBlob(tx db.Tx, id pacta.BlobID, mutations ...db.UpdateBlobFn) error {
	err := d.RunOrContinueTransaction(tx, func(db.Tx) error {
		b, err := d.Blob(tx, id)
		if err != nil {
			return fmt.Errorf("reading blob: %w", err)
		}
		for i, m := range mutations {
			err := m(b)
			if err != nil {
				return fmt.Errorf("running %d-th mutation: %w", i, err)
			}
		}
		err = d.putBlob(tx, b)
		if err != nil {
			return fmt.Errorf("putting blob: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("updating blob: %w", err)
	}
	return nil
}

func (d *DB) DeleteBlob(tx db.Tx, id pacta.BlobID) (pacta.BlobURI, error) {
	var buri pacta.BlobURI
	row := d.queryRow(tx, `DELETE FROM blob WHERE id = $1 RETURNING blob_uri;`, id)
	err := row.Scan(&buri)
	if err != nil {
		return "", fmt.Errorf("when deleting blob: %w", err)
	}
	return buri, nil
}

func (db *DB) putBlob(tx db.Tx, b *pacta.Blob) error {
	err := db.exec(tx, `
		UPDATE blob SET
			file_name = $2
		WHERE id = $1;
		`, b.ID, b.FileName)
	if err != nil {
		return fmt.Errorf("updating blob writable fields: %w", err)
	}
	return nil
}

func rowsToBlobs(rows pgx.Rows) ([]*pacta.Blob, error) {
	return allRows("blob", rows, rowToBlob)
}

func rowToBlob(row rowScanner) (*pacta.Blob, error) {
	b := &pacta.Blob{}
	fileType := ""
	err := row.Scan(
		&b.ID,
		&b.BlobURI,
		&fileType,
		&b.FileName,
		&b.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning into blob: %w", err)
	}
	ft, err := pacta.ParseFileType(fileType)
	if err != nil {
		return nil, fmt.Errorf("parsing blob file_type: %w", err)
	}
	b.FileType = ft
	return b, nil
}

func validateBlobForCreation(b *pacta.Blob) error {
	if b.ID != "" {
		return fmt.Errorf("blob already has an ID")
	}
	if b.BlobURI == "" {
		return fmt.Errorf("blob missing BlobURI")
	}
	if b.FileType == "" {
		return fmt.Errorf("blob missing FileType")
	}
	if b.FileName == "" {
		return fmt.Errorf("blob missing FileName")
	}
	return nil
}

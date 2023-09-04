package sqldb

import (
	"errors"
	"fmt"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const incompleteUploadSelectColumns = `
	incomplete_upload.id,
	incomplete_upload.owner_id,
	incomplete_upload.admin_debug_enabled,
	incomplete_upload.blob_id,
	incomplete_upload.name,
	incomplete_upload.description,
	incomplete_upload.holdings_date,
	incomplete_upload.created_at,
	incomplete_upload.ran_at,
	incomplete_upload.completed_at,
	incomplete_upload.failure_code,
	incomplete_upload.failure_message
`

func (d *DB) IncompleteUpload(tx db.Tx, id pacta.IncompleteUploadID) (*pacta.IncompleteUpload, error) {
	rows, err := d.query(tx, `
		SELECT `+incompleteUploadSelectColumns+`
		FROM incomplete_upload 
		WHERE id = $1;`, id)
	if err != nil {
		return nil, fmt.Errorf("querying incomplete_upload: %w", err)
	}
	pvs, err := rowsToIncompleteUploads(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to incomplete_uploads: %w", err)
	}
	return exactlyOne("incomplete_upload", id, pvs)
}

func (d *DB) IncompleteUploads(tx db.Tx, ids []pacta.IncompleteUploadID) (map[pacta.IncompleteUploadID]*pacta.IncompleteUpload, error) {
	rows, err := d.query(tx, `
		SELECT `+incompleteUploadSelectColumns+`
		FROM incomplete_upload 
		WHERE id IN `+createWhereInFmt(len(ids))+`;`, idsToInterface(ids)...)
	if err != nil {
		return nil, fmt.Errorf("querying incomplete_uploads: %w", err)
	}
	pvs, err := rowsToIncompleteUploads(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to incomplete_uploads: %w", err)
	}
	result := make(map[pacta.IncompleteUploadID]*pacta.IncompleteUpload, len(pvs))
	for _, pv := range pvs {
		result[pv.ID] = pv
	}
	return result, nil
}

func (d *DB) CreateIncompleteUpload(tx db.Tx, i *pacta.IncompleteUpload) error {
	if err := validateIncompleteUploadForCreation(i); err != nil {
		return fmt.Errorf("validating incomplete_upload for creation: %w", err)
	}
	hd, err := validateHoldingsDate(i.HoldingsDate)
	if err != nil {
		return fmt.Errorf("validating holdings date: %w", err)
	}
	i.ID = pacta.IncompleteUploadID(d.randomID("iu"))
	err = d.exec(tx, `
		INSERT INTO incomplete_upload 
			(id, owner_id, admin_debug_enabled, blob_id, name, description, holdings_date)
			VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
		i.ID, i.Owner.ID, i.AdminDebugEnabled, i.Blob.ID, i.Name, i.Description, hd)
	if err != nil {
		return fmt.Errorf("creating incomplete_upload: %w", err)
	}
	return nil
}

func (d *DB) UpdateIncompleteUpload(tx db.Tx, id pacta.IncompleteUploadID, mutations ...db.UpdateIncompleteUploadFn) error {
	err := d.RunOrContinueTransaction(tx, func(db.Tx) error {
		init, err := d.IncompleteUpload(tx, id)
		if err != nil {
			return fmt.Errorf("reading incomplete_upload: %w", err)
		}
		for i, m := range mutations {
			err := m(init)
			if err != nil {
				return fmt.Errorf("running %d-th mutation: %w", i, err)
			}
		}
		err = d.putIncompleteUpload(tx, init)
		if err != nil {
			return fmt.Errorf("putting incomplete_upload: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("updating incomplete_upload: %w", err)
	}
	return nil
}

func (d *DB) DeleteIncompleteUpload(tx db.Tx, id pacta.IncompleteUploadID) (pacta.BlobURI, error) {
	var buri pacta.BlobURI
	err := d.RunOrContinueTransaction(tx, func(tx db.Tx) error {
		iu, err := d.IncompleteUpload(tx, id)
		if err != nil {
			return fmt.Errorf("reading incomplete_upload: %w", err)
		}
		buri, err = d.DeleteBlob(tx, iu.Blob.ID)
		if err != nil {
			return fmt.Errorf("deleting blob: %w", err)
		}
		err = d.exec(tx, `DELETE FROM incomplete_upload WHERE id = $1;`, id)
		if err != nil {
			return fmt.Errorf("deleting incomplete_upload: %w", err)
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("performing incomplete_upload deletion: %w", err)
	}
	return buri, nil
}

func rowToIncompleteUpload(row rowScanner) (*pacta.IncompleteUpload, error) {
	iu := &pacta.IncompleteUpload{Owner: &pacta.Owner{}, Blob: &pacta.Blob{}}
	var (
		failureCode, failureMessage pgtype.Text
		hd, ranAt, completedAt      pgtype.Timestamptz
	)
	err := row.Scan(
		&iu.ID,
		&iu.Owner.ID,
		&iu.AdminDebugEnabled,
		&iu.Blob.ID,
		&iu.Name,
		&iu.Description,
		&hd,
		&iu.CreatedAt,
		&ranAt,
		&completedAt,
		&failureCode,
		&failureMessage,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning into incomplete_upload: %w", err)
	}
	iu.HoldingsDate, err = decodeHoldingsDate(hd)
	if err != nil {
		return nil, fmt.Errorf("decoding holdings date: %w", err)
	}
	if failureCode.Valid {
		iu.FailureCode, err = pacta.ParseFailureCode(failureCode.String)
		if err != nil {
			return nil, fmt.Errorf("parsing failure code: %w", err)
		}
	}
	if ranAt.Valid {
		iu.RanAt = ranAt.Time
	}
	if completedAt.Valid {
		iu.CompletedAt = completedAt.Time
	}
	if failureMessage.Valid {
		iu.FailureMessage = failureMessage.String
	}
	return iu, nil
}

func rowsToIncompleteUploads(rows pgx.Rows) ([]*pacta.IncompleteUpload, error) {
	return allRows("incomplete_upload", rows, rowToIncompleteUpload)
}

func (db *DB) putIncompleteUpload(tx db.Tx, iu *pacta.IncompleteUpload) error {
	hd, err := validateHoldingsDate(iu.HoldingsDate)
	if err != nil {
		return fmt.Errorf("validating holdings date: %w", err)
	}
	err = db.exec(tx, `
		UPDATE incomplete_upload SET
			owner_id = $2,
			admin_debug_enabled = $3,
			name = $4, 
			description = $5,
			holdings_date = $6,
			ran_at = $7,
			completed_at = $8,
			failure_code = $9,
			failure_message = $10
		WHERE id = $1;
		`, iu.Owner.ID, iu.AdminDebugEnabled, iu.Name, iu.Description,
		hd, timeToNilable(iu.RanAt), timeToNilable(iu.CompletedAt),
		strToNilable(iu.FailureCode), strToNilable(iu.FailureMessage))
	if err != nil {
		return fmt.Errorf("updating incomplete_upload writable fields: %w", err)
	}
	return nil
}

func validateIncompleteUploadForCreation(p *pacta.IncompleteUpload) error {
	if p.ID != "" {
		return errors.New("incomplete_upload id must be empty")
	}
	if p.Owner == nil || p.Owner.ID == "" {
		return errors.New("incomplete_upload must contain a non-nil owner with a present ID")
	}
	if p.Blob == nil || p.Blob.ID == "" {
		return errors.New("incomplete_upload must contain a non-nil blob with a present ID")
	}
	if p.Name == "" {
		return errors.New("incomplete_upload name must be present")
	}
	if !p.CreatedAt.IsZero() {
		return fmt.Errorf("incomplete_upload created_at must be zero")
	}
	return nil
}

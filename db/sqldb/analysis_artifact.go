package sqldb

import (
	"fmt"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
	"github.com/jackc/pgx/v5"
)

const analysisArtifactIDNamespace = "analysis_artifact"
const analysisArtifactSelectColumns = `
	analysis_artifact.id,
	analysis_artifact.analysis_id,
	analysis_artifact.blob_id,
	analysis_artifact.admin_debug_enabled,
	analysis_artifact.shared_to_public
`

func (d *DB) AnalysisArtifact(tx db.Tx, id pacta.AnalysisArtifactID) (*pacta.AnalysisArtifact, error) {
	rows, err := d.query(tx, `
		SELECT `+analysisArtifactSelectColumns+`
		FROM analysis_artifact 
		WHERE id = $1;`, id)
	if err != nil {
		return nil, fmt.Errorf("querying analysis_artifact: %w", err)
	}
	analysis_artifacts, err := rowsToAnalysisArtifacts(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to analysis_artifacts: %w", err)
	}
	return exactlyOne("analysis_artifact", id, analysis_artifacts)
}

func (d *DB) AnalysisArtifacts(tx db.Tx, id []pacta.AnalysisArtifactID) (map[pacta.AnalysisArtifactID]*pacta.AnalysisArtifact, error) {
	rows, err := d.query(tx, `
		SELECT `+analysisArtifactSelectColumns+`
		FROM analysis_artifact 
		WHERE id IN `+createWhereInFmt(len(id))+`;`, idsToInterface(id)...)
	if err != nil {
		return nil, fmt.Errorf("querying analysis_artifacts: %w", err)
	}
	aas, err := rowsToAnalysisArtifacts(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to analysis_artifacts: %w", err)
	}
	result := make(map[pacta.AnalysisArtifactID]*pacta.AnalysisArtifact, len(aas))
	for _, aa := range aas {
		result[aa.ID] = aa
	}
	return result, nil
}

func (d *DB) AnalysisArtifactsForAnalysis(tx db.Tx, id pacta.AnalysisID) ([]*pacta.AnalysisArtifact, error) {
	rows, err := d.query(tx, `
		SELECT `+analysisArtifactSelectColumns+`
		FROM analysis_artifact 
		WHERE analysis_id = $1;`, id)
	if err != nil {
		return nil, fmt.Errorf("querying analysis_artifacts: %w", err)
	}
	return rowsToAnalysisArtifacts(rows)
}

func (d *DB) CreateAnalysisArtifact(tx db.Tx, a *pacta.AnalysisArtifact) (pacta.AnalysisArtifactID, error) {
	if err := validateAnalysisArtifactForCreation(a); err != nil {
		return "", fmt.Errorf("validating analysis_artifact for creation: %w", err)
	}
	id := pacta.AnalysisArtifactID(d.randomID(analysisArtifactIDNamespace))
	err := d.exec(tx, `
		INSERT INTO analysis_artifact 
			(id, analysis_id, blob_id, admin_debug_enabled, shared_to_public)
			VALUES
			($1, $2, $3, $4);
	`, id, a.AnalysisID, a.Blob.ID, a.AdminDebugEnabled, a.SharedToPublic)
	if err != nil {
		return "", fmt.Errorf("creating analysis_artifact row: %w", err)
	}
	return id, nil
}

func (d *DB) UpdateAnalysisArtifact(tx db.Tx, id pacta.AnalysisArtifactID, mutations ...db.UpdateAnalysisArtifactFn) error {
	err := d.RunOrContinueTransaction(tx, func(db.Tx) error {
		b, err := d.AnalysisArtifact(tx, id)
		if err != nil {
			return fmt.Errorf("reading analysis_artifact: %w", err)
		}
		for i, m := range mutations {
			err := m(b)
			if err != nil {
				return fmt.Errorf("running %d-th mutation: %w", i, err)
			}
		}
		err = d.putAnalysisArtifact(tx, b)
		if err != nil {
			return fmt.Errorf("putting analysis_artifact: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("updating analysis_artifact: %w", err)
	}
	return nil
}

func (d *DB) DeleteAnalysisArtifact(tx db.Tx, id pacta.AnalysisArtifactID) (pacta.BlobURI, error) {
	var buri pacta.BlobURI
	err := d.RunOrContinueTransaction(tx, func(tx db.Tx) error {
		var blobID pacta.BlobID
		row := d.queryRow(tx, `DELETE FROM analysis_artifact WHERE id = $1 RETURNING blob_id;`, id)
		err := row.Scan(&blobID)
		if err != nil {
			return fmt.Errorf("when deleting analysis_artifact: %w", err)
		}
		buri, err = d.DeleteBlob(tx, blobID)
		if err != nil {
			return fmt.Errorf("when deleting blob: %w", err)
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("deleting analysis_artifact: %w", err)
	}
	return buri, nil
}

func (db *DB) putAnalysisArtifact(tx db.Tx, a *pacta.AnalysisArtifact) error {
	err := db.exec(tx, `
		UPDATE analysis_artifact SET
			admin_debug_enabled = $2,
			shared_to_public = $3
		WHERE id = $1;
		`, a.ID, a.AdminDebugEnabled, a.SharedToPublic)
	if err != nil {
		return fmt.Errorf("updating analysis_artifact writable fields: %w", err)
	}
	return nil
}

func rowsToAnalysisArtifacts(rows pgx.Rows) ([]*pacta.AnalysisArtifact, error) {
	return allRows("analysis_artifact", rows, rowToAnalysisArtifact)
}

func rowToAnalysisArtifact(row rowScanner) (*pacta.AnalysisArtifact, error) {
	a := &pacta.AnalysisArtifact{Blob: &pacta.Blob{}}
	err := row.Scan(
		&a.ID,
		&a.AnalysisID,
		&a.Blob.ID,
		&a.AdminDebugEnabled,
		&a.SharedToPublic,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning into analysis_artifact: %w", err)
	}
	return a, nil
}

func validateAnalysisArtifactForCreation(a *pacta.AnalysisArtifact) error {
	if a.ID != "" {
		return fmt.Errorf("analysis_artifact already has an ID")
	}
	if a.Blob == nil || a.Blob.ID == "" {
		return fmt.Errorf("analysis_artifact is missing blob")
	}
	if a.AnalysisID == "" {
		return fmt.Errorf("analysis_artifact is missing analysis_id")
	}
	return nil
}

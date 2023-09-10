package sqldb

import (
	"fmt"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const analyseSelectColumns = `
	analysis.id,
	analysis.analysis_type,
	analysis.owner_id,
	analysis.pacta_version_id,
	analysis.portfolio_snapshot_id,
	analysis.name,
	analysis.description,
	analysis.created_at,
	analysis.ran_at,
	analysis.completed_at,
	analysis.failure_code,
	analysis.failure_message`

func (d *DB) Analysis(tx db.Tx, id pacta.AnalysisID) (*pacta.Analysis, error) {
	rows, err := d.query(tx, `
		SELECT `+analyseSelectColumns+`
		FROM analysis 
		WHERE id = $1;`, id)
	if err != nil {
		return nil, fmt.Errorf("querying analysis: %w", err)
	}
	pvs, err := rowsToAnalyses(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to analyses: %w", err)
	}
	return exactlyOne("analysis", id, pvs)
}

func (d *DB) Analyses(tx db.Tx, ids []pacta.AnalysisID) (map[pacta.AnalysisID]*pacta.Analysis, error) {
	ids = dedupeIDs(ids)
	rows, err := d.query(tx, `
		SELECT `+analyseSelectColumns+`
		FROM analysis 
		WHERE id IN `+createWhereInFmt(len(ids))+`;`, idsToInterface(ids)...)
	if err != nil {
		return nil, fmt.Errorf("querying analyses: %w", err)
	}
	pvs, err := rowsToAnalyses(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to analyses: %w", err)
	}
	result := make(map[pacta.AnalysisID]*pacta.Analysis, len(pvs))
	for _, pv := range pvs {
		result[pv.ID] = pv
	}
	return result, nil
}

func (d *DB) CreateAnalysis(tx db.Tx, a *pacta.Analysis) (pacta.AnalysisID, error) {
	if err := validateAnalysisForCreation(a); err != nil {
		return "", fmt.Errorf("validating analysis for creation: %w", err)
	}
	a.ID = pacta.AnalysisID(d.randomID("analysis"))
	err := d.exec(tx, `
		INSERT INTO analysis 
			(id, analysis_type, owner_id, pacta_version_id, portfolio_snapshot_id, name, description)
			VALUES
			($1, $2, $3, $4, $5, $6, $7);`,
		a.ID, a.AnalysisType, a.Owner.ID, a.PACTAVersion.ID, a.PortfolioSnapshot.ID, a.Name, a.Description)
	if err != nil {
		return "", fmt.Errorf("creating analysis: %w", err)
	}
	return a.ID, nil
}

func (d *DB) UpdateAnalysis(tx db.Tx, id pacta.AnalysisID, mutations ...db.UpdateAnalysisFn) error {
	err := d.RunOrContinueTransaction(tx, func(db.Tx) error {
		init, err := d.Analysis(tx, id)
		if err != nil {
			return fmt.Errorf("reading analysis: %w", err)
		}
		for i, m := range mutations {
			err := m(init)
			if err != nil {
				return fmt.Errorf("running %d-th mutation: %w", i, err)
			}
		}
		err = d.putAnalysis(tx, init)
		if err != nil {
			return fmt.Errorf("putting analysis: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("updating analysis: %w", err)
	}
	return nil
}

func (d *DB) DeleteAnalysis(tx db.Tx, id pacta.AnalysisID) ([]pacta.BlobURI, error) {
	buris := []pacta.BlobURI{}
	err := d.RunOrContinueTransaction(tx, func(tx db.Tx) error {
		a, err := d.Analysis(tx, id)
		if err != nil {
			return fmt.Errorf("reading analysis: %w", err)
		}
		aas, err := d.AnalysisArtifactsForAnalysis(tx, id)
		if err != nil {
			return fmt.Errorf("reading analysis artifacts: %w", err)
		}
		for _, aa := range aas {
			buri, err := d.DeleteAnalysisArtifact(tx, aa.ID)
			if err != nil {
				return fmt.Errorf("deleting analysis artifact %q: %w", aa.ID, err)
			}
			buris = append(buris, buri)
		}
		err = d.exec(tx, `DELETE FROM analysis WHERE id = $1;`, id)
		if err != nil {
			return fmt.Errorf("deleting analysis: %w", err)
		}
		err = d.exec(tx, `DELETE FROM portfolio_snapshot WHERE id = $1;`, a.PortfolioSnapshot.ID)
		if err != nil {
			return fmt.Errorf("deleting analysis_invitations: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("performing analysis sub-entity deletion: %w", err)
	}
	return buris, nil
}

func rowToAnalysis(row rowScanner) (*pacta.Analysis, error) {
	a := &pacta.Analysis{
		Owner:             &pacta.Owner{},
		PACTAVersion:      &pacta.PACTAVersion{},
		PortfolioSnapshot: &pacta.PortfolioSnapshot{},
	}
	var (
		aType                       string
		failureCode, failureMessage pgtype.Text
		ranAt, completedAt          pgtype.Timestamptz
	)
	err := row.Scan(
		&a.ID,
		&aType,
		&a.Owner.ID,
		&a.PACTAVersion.ID,
		&a.PortfolioSnapshot.ID,
		&a.Name,
		&a.Description,
		&a.CreatedAt,
		&ranAt,
		&completedAt,
		&failureCode,
		&failureMessage,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning into analysis: %w", err)
	}
	a.AnalysisType, err = pacta.ParseAnalysisType(aType)
	if err != nil {
		return nil, fmt.Errorf("parsing analysis type: %w", err)
	}
	if failureCode.Valid {
		a.FailureCode, err = pacta.ParseFailureCode(failureCode.String)
		if err != nil {
			return nil, fmt.Errorf("parsing failure code: %w", err)
		}
	}
	if ranAt.Valid {
		a.RanAt = ranAt.Time
	}
	if completedAt.Valid {
		a.CompletedAt = completedAt.Time
	}
	if failureMessage.Valid {
		a.FailureMessage = failureMessage.String
	}
	return a, nil
}

func rowsToAnalyses(rows pgx.Rows) ([]*pacta.Analysis, error) {
	return allRows("analysis", rows, rowToAnalysis)
}

func (db *DB) putAnalysis(tx db.Tx, a *pacta.Analysis) error {
	err := db.exec(tx, `
		UPDATE analysis SET
			owner_id = $2,
			name = $3, 
			description = $4,
			ran_at = $5,
			completed_at = $6,
			failure_code = $7,
			failure_message = $8
		WHERE id = $1;
		`,
		a.ID,
		a.Owner.ID,
		a.Name,
		a.Description,
		timeToNilable(a.RanAt),
		timeToNilable(a.CompletedAt),
		strToNilable(a.FailureCode),
		strToNilable(a.FailureMessage),
	)
	if err != nil {
		return fmt.Errorf("updating analysis writable fields: %w", err)
	}
	return nil
}

func validateAnalysisForCreation(a *pacta.Analysis) error {
	if a.ID != "" {
		return fmt.Errorf("analysis id must be absent")
	}
	if !a.CreatedAt.IsZero() {
		return fmt.Errorf("analysis created_at must be zero")
	}
	if a.Owner == nil || a.Owner.ID == "" {
		return fmt.Errorf("analysis owner must be present")
	}
	if a.PACTAVersion == nil || a.PACTAVersion.ID == "" {
		return fmt.Errorf("analysis pacta_version must be present")
	}
	if a.PortfolioSnapshot == nil || a.PortfolioSnapshot.ID == "" {
		return fmt.Errorf("analysis portfolio_snapshot must be present")
	}
	if a.AnalysisType == "" {
		return fmt.Errorf("analysis type must be present")
	}
	return nil
}

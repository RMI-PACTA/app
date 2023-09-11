package sqldb

import (
	"fmt"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
	"github.com/jackc/pgx/v5"
)

const pactaVersionIDNamespace = "pv"
const pactaVersionSelectColumns = `
	pacta_version.id,
	pacta_version.name,
	pacta_version.description,
	pacta_version.digest,
	pacta_version.created_at,
	COALESCE(pacta_version.is_default, false)`

func (d *DB) PACTAVersion(tx db.Tx, id pacta.PACTAVersionID) (*pacta.PACTAVersion, error) {
	rows, err := d.query(tx, `
		SELECT `+pactaVersionSelectColumns+`
		FROM pacta_version 
		WHERE id = $1;`, id)
	if err != nil {
		return nil, fmt.Errorf("querying pacta_version: %w", err)
	}
	pvs, err := rowsToPACTAVersions(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to pacta versions: %w", err)
	}
	return exactlyOne("pacta_version", id, pvs)
}

func (d *DB) DefaultPACTAVersion(tx db.Tx) (*pacta.PACTAVersion, error) {
	rows, err := d.query(tx, `
		SELECT `+pactaVersionSelectColumns+`
		FROM pacta_version 
		WHERE is_default;`)
	if err != nil {
		return nil, fmt.Errorf("querying primary pacta_version: %w", err)
	}
	pvs, err := rowsToPACTAVersions(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to pacta versions: %w", err)
	}
	return exactlyOne("pacta_version", "is_default", pvs)
}

func (d *DB) PACTAVersions(tx db.Tx) ([]*pacta.PACTAVersion, error) {
	rows, err := d.query(tx, `
		SELECT `+pactaVersionSelectColumns+`
		FROM pacta_version;`)
	if err != nil {
		return nil, fmt.Errorf("querying pacta_version: %w", err)
	}
	pvs, err := mapRows("pacta_version", rows, rowToPACTAVersion)
	if err != nil {
		return nil, fmt.Errorf("translating rows to pacta versions: %w", err)
	}
	return pvs, nil
}

func (d *DB) CreatePACTAVersion(tx db.Tx, pv *pacta.PACTAVersion) (pacta.PACTAVersionID, error) {
	if err := validatePACTAVersionForCreation(pv); err != nil {
		return "", fmt.Errorf("failed pacta_version validation: %w", err)
	}
	id := pacta.PACTAVersionID(d.randomID(pactaVersionIDNamespace))
	err := d.exec(tx, `
		INSERT INTO pacta_version 
			(id, name, description, digest, is_default)
			VALUES
			($1, $2, $3, $4, $5);
	`, id, pv.Name, pv.Description, pv.Digest, nil)
	if err != nil {
		return "", fmt.Errorf("creating pacta_version: %w", err)
	}
	return id, nil
}

func (d *DB) SetDefaultPACTAVersion(tx db.Tx, id pacta.PACTAVersionID) error {
	err := d.RunOrContinueTransaction(tx, func(db.Tx) error {
		err := d.exec(tx, `UPDATE pacta_version SET is_default = NULL;`)
		if err != nil {
			return fmt.Errorf("setting all pacta_versions is_default to null: %w", err)
		}
		err = d.exec(tx, `UPDATE pacta_version SET is_default = TRUE WHERE id = $1;`, id)
		if err != nil {
			return fmt.Errorf("setting pacta version with id %q is_default: %w", id, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("updating primary pacta_version: %w", err)
	}
	return nil
}

func (d *DB) UpdatePACTAVersion(tx db.Tx, id pacta.PACTAVersionID, mutations ...db.UpdatePACTAVersionFn) error {
	err := d.RunOrContinueTransaction(tx, func(db.Tx) error {
		pv, err := d.PACTAVersion(tx, id)
		if err != nil {
			return fmt.Errorf("reading pacta_version: %w", err)
		}
		for i, m := range mutations {
			err := m(pv)
			if err != nil {
				return fmt.Errorf("running %d-th mutation: %w", i, err)
			}
		}
		err = d.putPACTAVersion(tx, pv)
		if err != nil {
			return fmt.Errorf("putting pacta_version: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("updating pacta_version: %w", err)
	}
	return nil
}

func (d *DB) DeletePACTAVersion(tx db.Tx, id pacta.PACTAVersionID) error {
	err := d.RunOrContinueTransaction(tx, func(db.Tx) error {
		// TODO - include appropriate handling for deleting a pacta_version that is used.
		err := d.exec(tx, `DELETE FROM pacta_version WHERE id = $1;`, id)
		if err != nil {
			return fmt.Errorf("deleting pacta_version: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("performing pacta_version deletion: %w", err)
	}
	return nil
}

func rowToPACTAVersion(row rowScanner) (*pacta.PACTAVersion, error) {
	p := &pacta.PACTAVersion{}
	err := row.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Digest,
		&p.CreatedAt,
		&p.IsDefault,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning into pacta_version: %w", err)
	}
	return p, err
}

func rowsToPACTAVersions(rows pgx.Rows) ([]*pacta.PACTAVersion, error) {
	return mapRows("pacta_version", rows, rowToPACTAVersion)
}

func (d *DB) putPACTAVersion(tx db.Tx, pv *pacta.PACTAVersion) error {
	var isDefault *bool
	if pv.IsDefault {
		isDefault = &pv.IsDefault
	}
	err := d.exec(tx, `
		UPDATE pacta_version SET
			name = $2,
			description = $3,
			digest = $4,
			is_default = $5
		WHERE id = $1;
		`, pv.ID, pv.Name, pv.Description, pv.Digest, isDefault)
	if err != nil {
		return fmt.Errorf("updating pacta_version writable fields: %w", err)
	}
	return nil
}

func validatePACTAVersionForCreation(pv *pacta.PACTAVersion) error {
	if pv.IsDefault {
		return fmt.Errorf("cannot create a default pacta_version - use SetDefaultPACTAVersion() instead")
	}
	if !pv.CreatedAt.IsZero() {
		return fmt.Errorf("cannot set created_at on creation")
	}
	if pv.Digest == "" {
		return fmt.Errorf("digest is required")
	}
	if pv.Name == "" {
		return fmt.Errorf("name is required")
	}
	if pv.ID != "" {
		return fmt.Errorf("cannot set id on creation")
	}
	return nil
}

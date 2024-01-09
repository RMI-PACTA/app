package sqldb

import (
	"fmt"
	"time"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
	"github.com/jackc/pgx/v5"
)

const initiativeUserRelationshipSelectColumns = `
	initiative_user_relationship.initiative_id,
	initiative_user_relationship.user_id,
	initiative_user_relationship.member,
	initiative_user_relationship.manager,
	initiative_user_relationship.updated_at
`

func (d *DB) InitiativeUserRelationship(tx db.Tx, iid pacta.InitiativeID, uid pacta.UserID) (*pacta.InitiativeUserRelationship, error) {
	rows, err := d.query(tx, `
		SELECT `+initiativeUserRelationshipSelectColumns+`
		FROM initiative_user_relationship 
		WHERE initiative_id = $1 AND user_id = $2;`, iid, uid)
	if err != nil {
		return nil, fmt.Errorf("querying initiative_user_relationship: %w", err)
	}
	iurs, err := rowsToInitiativeUserRelationships(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to initiative_user_relationships: %w", err)
	}
	return exactlyOne("initiative_user_relationship", fmt.Sprintf("%s:%s", iid, uid), iurs)
}

func (d *DB) InitiativeUserRelationshipsByUser(tx db.Tx, uid pacta.UserID) ([]*pacta.InitiativeUserRelationship, error) {
	rows, err := d.query(tx, `
		SELECT `+initiativeUserRelationshipSelectColumns+`
		FROM initiative_user_relationship 
		WHERE user_id = $1;`, uid)
	if err != nil {
		return nil, fmt.Errorf("querying initiative_user_relationship: %w", err)
	}
	iurs, err := rowsToInitiativeUserRelationships(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to initiative_user_relationships: %w", err)
	}
	return iurs, nil
}

func (d *DB) InitiativeUserRelationshipsByInitiative(tx db.Tx, iid pacta.InitiativeID) ([]*pacta.InitiativeUserRelationship, error) {
	rows, err := d.query(tx, `
		SELECT `+initiativeUserRelationshipSelectColumns+`
		FROM initiative_user_relationship 
		WHERE initiative_id = $1;`, iid)
	if err != nil {
		return nil, fmt.Errorf("querying initiative_user_relationship: %w", err)
	}
	iurs, err := rowsToInitiativeUserRelationships(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to initiative_user_relationships: %w", err)
	}
	return iurs, nil
}

func (d *DB) PutInitiativeUserRelationship(tx db.Tx, iur *pacta.InitiativeUserRelationship) error {
	if err := validateInitiativeUserRelationshipForPut(iur); err != nil {
		return fmt.Errorf("validating initiative_user_relationship for put: %w", err)
	}
	updatedAt := time.Now()
	err := d.exec(tx, `
		INSERT INTO initiative_user_relationship
			(initiative_id, user_id, member, manager, updated_at) 
			VALUES
			($1, $2, $3, $4, $5)
		ON CONFLICT (initiative_id, user_id) DO UPDATE SET
			member = $3,
			manager = $4,
			updated_at = $5;
		`, iur.Initiative.ID, iur.User.ID, iur.Member, iur.Manager, updatedAt)
	if err != nil {
		return fmt.Errorf("updating initiative_user_relationship writable fields: %w", err)
	}
	return nil
}

func (d *DB) UpdateInitiativeUserRelationship(tx db.Tx, iid pacta.InitiativeID, uid pacta.UserID, mutations ...db.UpdateInitiativeUserRelationshipFn) error {
	iur, err := d.InitiativeUserRelationship(tx, iid, uid)
	if err != nil {
		if db.IsNotFound(err) {
			iur = &pacta.InitiativeUserRelationship{
				Initiative: &pacta.Initiative{ID: iid},
				User:       &pacta.User{ID: uid},
			}
		} else {
			return fmt.Errorf("getting initiative_user_relationship: %w", err)
		}
	}
	for _, mutation := range mutations {
		if err := mutation(iur); err != nil {
			return fmt.Errorf("applying mutation to initiative_user_relationship: %w", err)
		}
	}
	iur.UpdatedAt = time.Now()
	return d.PutInitiativeUserRelationship(tx, iur)
}

func validateInitiativeUserRelationshipForPut(iur *pacta.InitiativeUserRelationship) error {
	if iur.User == nil || iur.User.ID == "" {
		return fmt.Errorf("user_id is required")
	}
	if iur.Initiative == nil || iur.Initiative.ID == "" {
		return fmt.Errorf("initiative_id is required")
	}
	return nil
}

func rowToInitiativeUserRelationship(row rowScanner) (*pacta.InitiativeUserRelationship, error) {
	iuv := &pacta.InitiativeUserRelationship{
		Initiative: &pacta.Initiative{},
		User:       &pacta.User{},
	}
	err := row.Scan(
		&iuv.Initiative.ID,
		&iuv.User.ID,
		&iuv.Member,
		&iuv.Manager,
		&iuv.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning into initiative_user_relationship: %w", err)
	}
	return iuv, nil
}

func rowsToInitiativeUserRelationships(rows pgx.Rows) ([]*pacta.InitiativeUserRelationship, error) {
	return mapRows("initiative_user_relationship", rows, rowToInitiativeUserRelationship)
}

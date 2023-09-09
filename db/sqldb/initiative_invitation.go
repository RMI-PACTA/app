package sqldb

import (
	"fmt"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const initiativeInvitiationIDNamespace = "iivte"
const initiativeIntivtationSelectColumns = `
	initiative_invitation.id,
	initiative_invitation.created_at,
	initiative_invitation.used_at,
	initiative_invitation.initiative_id,
	initiative_invitation.used_by_user_id`

func (d *DB) InitiativeInvitation(tx db.Tx, id pacta.InitiativeInvitationID) (*pacta.InitiativeInvitation, error) {
	rows, err := d.query(tx, `
		SELECT `+initiativeIntivtationSelectColumns+`
		FROM initiative_invitation 
		WHERE id = $1;`, id)
	if err != nil {
		return nil, fmt.Errorf("querying initiative_invitation: %w", err)
	}
	iis, err := rowsToInitiativeInvitations(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to initiative_invitations: %w", err)
	}
	return exactlyOne("initiative_invitation", id, iis)
}

func (d *DB) InitiativeInvitationsByInitiative(tx db.Tx, iid pacta.InitiativeID) ([]*pacta.InitiativeInvitation, error) {
	rows, err := d.query(tx, `
		SELECT `+initiativeIntivtationSelectColumns+`
		FROM initiative_invitation 
		WHERE initiative_id = $1;`, iid)
	if err != nil {
		return nil, fmt.Errorf("querying initiative_invitation: %w", err)
	}
	iis, err := rowsToInitiativeInvitations(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to initiative_invitations: %w", err)
	}
	return iis, nil
}

func (d *DB) CreateInitiativeInvitation(tx db.Tx, ii *pacta.InitiativeInvitation) (pacta.InitiativeInvitationID, error) {
	if err := validateInitiativeInvitationForCreation(ii); err != nil {
		return "", fmt.Errorf("validating initiative_invitation for creation: %w", err)
	}
	// Unlike other entities, II's can have user-set IDs so as to appear semantically meaningful.
	if ii.ID == "" {
		ii.ID = pacta.InitiativeInvitationID(d.randomID(initiativeInvitiationIDNamespace))
	}
	err := d.exec(tx, `
		INSERT INTO initiative_invitation
			(id, initiative_id)					
			VALUES
			($1, $2)`, ii.ID, ii.Initiative.ID)
	if err != nil {
		return "", fmt.Errorf("creating initiative_invitation: %w", err)
	}
	return ii.ID, nil
}

func (d *DB) UpdateInitiativeInvitation(tx db.Tx, id pacta.InitiativeInvitationID, mutations ...db.UpdateInitiativeInvitationFn) error {
	err := d.RunOrContinueTransaction(tx, func(db.Tx) error {
		ii, err := d.InitiativeInvitation(tx, id)
		if err != nil {
			return fmt.Errorf("reading initiative_invitation: %w", err)
		}
		for i, m := range mutations {
			err := m(ii)
			if err != nil {
				return fmt.Errorf("running %d-th mutation: %w", i, err)
			}
		}
		err = d.putInitiativeInvitation(tx, ii)
		if err != nil {
			return fmt.Errorf("putting initiative_invitation: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("updating initiative_invitation: %w", err)
	}
	return nil
}

func (d *DB) DeleteInitiativeInvitation(tx db.Tx, id pacta.InitiativeInvitationID) error {
	if err := d.exec(tx, `DELETE FROM initiative_invitation WHERE id = $1;`, id); err != nil {
		return fmt.Errorf("deleting initiative_invitation: %w", err)
	}
	return nil
}

func rowToInitiativeInvitation(row rowScanner) (*pacta.InitiativeInvitation, error) {
	ii := &pacta.InitiativeInvitation{Initiative: &pacta.Initiative{}}
	ubid := pgtype.Text{}
	t := pgtype.Timestamptz{}
	err := row.Scan(
		&ii.ID,
		&ii.CreatedAt,
		&t,
		&ii.Initiative.ID,
		&ubid,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning into initiative_invitation: %w", err)
	}
	if ubid.Valid {
		ii.UsedBy = &pacta.User{ID: pacta.UserID(ubid.String)}
	}
	if t.Valid {
		ii.UsedAt = t.Time
	}
	return ii, nil
}

func rowsToInitiativeInvitations(rows pgx.Rows) ([]*pacta.InitiativeInvitation, error) {
	return allRows("pacta_version", rows, rowToInitiativeInvitation)
}

func validateInitiativeInvitationForCreation(ii *pacta.InitiativeInvitation) error {
	if ii.Initiative == nil || ii.Initiative.ID == "" {
		return fmt.Errorf("InitiativeInvitation.Initiative.ID must not be empty")
	}
	if !ii.UsedAt.IsZero() {
		return fmt.Errorf("InitiativeInvitation.UsedAt must be zero")
	}
	if ii.UsedBy != nil {
		return fmt.Errorf("InitiativeInvitation.UsedBy must be nil")
	}
	return nil
}

func (db *DB) putInitiativeInvitation(tx db.Tx, ii *pacta.InitiativeInvitation) error {
	var uid pgtype.Text
	if ii.UsedBy != nil {
		uid.Valid = true
		uid.String = string(ii.UsedBy.ID)
	}
	err := db.exec(tx, `
		UPDATE initiative_invitation SET
			used_at = $2,
			used_by_user_id = $3
		WHERE id = $1;
		`, ii.ID, ii.UsedAt, uid)
	if err != nil {
		return fmt.Errorf("updating initiative_invitation writable fields: %w", err)
	}
	return nil
}

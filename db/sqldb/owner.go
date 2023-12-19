package sqldb

import (
	"fmt"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const ownerIDNamespace = "own"
const ownerSelectColumns = `
	owner.id,
	owner.user_id,
	owner.initiative_id
`

func (d *DB) Owners(tx db.Tx, ids []pacta.OwnerID) (map[pacta.OwnerID]*pacta.Owner, error) {
	ids = dedupeIDs(ids)
	rows, err := d.query(tx, `
		SELECT `+ownerSelectColumns+`
		FROM owner
		WHERE owner.id IN `+createWhereInFmt(len(ids))+`;`, idsToInterface(ids)...)
	if err != nil {
		return nil, fmt.Errorf("querying owners: %w", err)
	}
	return rowsToOwners(rows)
}

func (d *DB) Owner(tx db.Tx, id pacta.OwnerID) (*pacta.Owner, error) {
	owners, err := d.Owners(tx, []pacta.OwnerID{id})
	if err != nil {
		return nil, fmt.Errorf("querying owner: %w", err)
	}
	return exactlyOneFromMap("owner", id, owners)
}

func (d *DB) ownerByUser(tx db.Tx, id pacta.UserID) (*pacta.Owner, error) {
	rows, err := d.query(tx, `SELECT `+ownerSelectColumns+` FROM owner WHERE owner.user_id = $1;`, id)
	if err != nil {
		return nil, fmt.Errorf("querying owners: %w", err)
	}
	owners, err := rowsToOwners(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to owners: %w", err)
	}
	for _, owner := range owners {
		if owner.User != nil && owner.User.ID == id {
			return owner, nil
		}
	}
	return nil, db.NotFound(id, "ownerByUserId")
}

func (d *DB) ownerByInitiative(tx db.Tx, id pacta.InitiativeID) (*pacta.Owner, error) {
	rows, err := d.query(tx, `SELECT `+ownerSelectColumns+` FROM owner WHERE owner.initiative_id = $1;`, id)
	if err != nil {
		return nil, fmt.Errorf("querying owners: %w", err)
	}
	owners, err := rowsToOwners(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to owners: %w", err)
	}
	for _, owner := range owners {
		if owner.Initiative != nil && owner.Initiative.ID == id {
			return owner, nil
		}
	}
	return nil, db.NotFound(id, "ownerByInitiativeId")
}

func (d *DB) GetOwnerForUser(tx db.Tx, uID pacta.UserID) (pacta.OwnerID, error) {
	var ownerID pacta.OwnerID
	err := d.RunOrContinueTransaction(tx, func(tx db.Tx) error {
		owner, err := d.ownerByUser(tx, uID)
		if err == nil {
			ownerID = owner.ID
			return nil
		}
		if !db.IsNotFound(err) {
			return fmt.Errorf("user owner not found: %w", err)
		}
		return fmt.Errorf("looking up owner: %w", err)
	})
	if err != nil {
		return "", fmt.Errorf("getting or creating owner for initiative: %w", err)
	}
	return ownerID, nil
}

func (d *DB) GetOrCreateOwnerForInitiative(tx db.Tx, iID pacta.InitiativeID) (pacta.OwnerID, error) {
	var ownerID pacta.OwnerID
	err := d.RunOrContinueTransaction(tx, func(tx db.Tx) error {
		owner, err := d.ownerByInitiative(tx, iID)
		if err == nil {
			ownerID = owner.ID
			return nil
		}
		if !db.IsNotFound(err) {
			return fmt.Errorf("querying owner by user: %w", err)
		}
		ownerID, err = d.createOwner(tx, &pacta.Owner{Initiative: &pacta.Initiative{ID: iID}})
		if err != nil {
			return fmt.Errorf("creating owner: %w", err)
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("getting or creating owner for initiative: %w", err)
	}
	return ownerID, nil
}

func (d *DB) createOwner(tx db.Tx, o *pacta.Owner) (pacta.OwnerID, error) {
	if err := validateOwnerForCreation(o); err != nil {
		return "", fmt.Errorf("validating owner for creation: %w", err)
	}
	var uID *pacta.UserID
	var iID *pacta.InitiativeID
	if o.User != nil && o.User.ID != "" {
		uID = &o.User.ID
	}
	if o.Initiative != nil && o.Initiative.ID != "" {
		iID = &o.Initiative.ID
	}
	id := pacta.OwnerID(d.randomID(ownerIDNamespace))
	err := d.exec(tx, `
		INSERT INTO owner 
			(id, user_id, initiative_id)
			VALUES
			($1, $2, $3);
	`, id, uID, iID)
	if err != nil {
		return "", fmt.Errorf("creating owner row: %w", err)
	}
	return id, nil
}

func rowsToOwners(rows pgx.Rows) (map[pacta.OwnerID]*pacta.Owner, error) {
	owners, err := mapRows("owner", rows, rowToOwner)
	if err != nil {
		return nil, fmt.Errorf("translating rows to owners: %w", err)
	}
	result := make(map[pacta.OwnerID]*pacta.Owner, len(owners))
	for _, owner := range owners {
		result[owner.ID] = owner
	}
	return result, nil
}

func rowToOwner(row rowScanner) (*pacta.Owner, error) {
	o := &pacta.Owner{}
	var uID, iID pgtype.Text
	err := row.Scan(
		&o.ID,
		&uID,
		&iID,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning into owner: %w", err)
	}
	if uID.Valid {
		o.User = &pacta.User{ID: pacta.UserID(uID.String)}
	}
	if iID.Valid {
		o.Initiative = &pacta.Initiative{ID: pacta.InitiativeID(iID.String)}
	}
	return o, nil
}

func validateOwnerForCreation(o *pacta.Owner) error {
	if o.ID != "" {
		return fmt.Errorf("owner already has an ID")
	}
	hasUser := o.User != nil && o.User.ID != ""
	hasInitiative := o.Initiative != nil && o.Initiative.ID != ""
	if !hasUser && !hasInitiative {
		return fmt.Errorf("owner missing User or Initiative")
	}
	if hasUser && hasInitiative {
		return fmt.Errorf("owner has both User and Initiative")
	}
	return nil
}

// TODO(grady) take on owner deletion

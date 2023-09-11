package sqldb

import (
	"fmt"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const userIDNamespace = "user"
const userSelectColumns = `
	pacta_user.id,
	pacta_user.authn_mechanism,
	pacta_user.authn_id,
	pacta_user.entered_email,
	pacta_user.canonical_email,
	pacta_user.admin,
	pacta_user.super_admin,
	pacta_user.name,
	pacta_user.preferred_language,
	pacta_user.created_at`

func (d *DB) User(tx db.Tx, id pacta.UserID) (*pacta.User, error) {
	rows, err := d.query(tx, `
		SELECT `+userSelectColumns+`
		FROM pacta_user 
		WHERE id = $1;`, id)
	if err != nil {
		return nil, fmt.Errorf("querying user: %w", err)
	}
	us, err := rowsToUsers(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to users: %w", err)
	}
	return exactlyOne("user", id, us)
}

func (d *DB) UserByAuthn(tx db.Tx, authnMechanism pacta.AuthnMechanism, authnID string) (*pacta.User, error) {
	rows, err := d.query(tx, `
		SELECT `+userSelectColumns+`
		FROM pacta_user 
		WHERE authn_mechanism = $1 AND authn_id = $2;`, authnMechanism, authnID)
	if err != nil {
		return nil, fmt.Errorf("querying user: %w", err)
	}
	us, err := rowsToUsers(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to users: %w", err)
	}
	return exactlyOne("user", fmt.Sprintf("%s:%s", authnMechanism, authnID), us)
}

func (d *DB) Users(tx db.Tx, ids []pacta.UserID) (map[pacta.UserID]*pacta.User, error) {
	ids = dedupeIDs(ids)
	rows, err := d.query(tx, `
		SELECT `+userSelectColumns+`
		FROM pacta_user 
		WHERE id IN `+createWhereInFmt(len(ids))+`;`, idsToInterface(ids)...)
	if err != nil {
		return nil, fmt.Errorf("querying users: %w", err)
	}
	us, err := rowsToUsers(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to users: %w", err)
	}
	result := make(map[pacta.UserID]*pacta.User)
	for _, u := range us {
		result[u.ID] = u
	}
	return result, nil
}

func (d *DB) CreateUser(tx db.Tx, u *pacta.User) (pacta.UserID, error) {
	if err := validateUserForCreation(u); err != nil {
		return "", fmt.Errorf("validating user for creation: %w", err)
	}
	var pl pgtype.Text
	if u.PreferredLanguage != "" {
		pl.Valid = true
		pl.String = string(u.PreferredLanguage)
	}
	id := pacta.UserID(d.randomID(userIDNamespace))
	err := d.exec(tx, `
		INSERT INTO pacta_user 
			(id, authn_mechanism, authn_id, entered_email, canonical_email, admin, super_admin, name, preferred_language)
			VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9);
		`, id, u.AuthnMechanism, u.AuthnID, u.EnteredEmail, u.CanonicalEmail, false, false, u.Name, pl)
	if err != nil {
		return "", fmt.Errorf("creating pacta_user row for %q: %w", id, err)
	}
	return id, nil
}

func (d *DB) UpdateUser(tx db.Tx, id pacta.UserID, mutations ...db.UpdateUserFn) error {
	err := d.RunOrContinueTransaction(tx, func(db.Tx) error {
		u, err := d.User(tx, id)
		if err != nil {
			return fmt.Errorf("reading user: %w", err)
		}
		for i, m := range mutations {
			err := m(u)
			if err != nil {
				return fmt.Errorf("running %d-th mutation: %w", i, err)
			}
		}
		err = d.putUser(tx, u)
		if err != nil {
			return fmt.Errorf("putting user: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("updating user: %w", err)
	}
	return nil
}

func (d *DB) DeleteUser(tx db.Tx, id pacta.UserID) error {
	err := d.RunOrContinueTransaction(tx, func(db.Tx) error {
		// TODO(grady) add entity deletions here
		err := d.exec(tx, `DELETE FROM pacta_user WHERE id = $1;`, id)
		if err != nil {
			return fmt.Errorf("deleting user: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("performing initiative deletion: %w", err)
	}
	return nil
}

func (d *DB) putUser(tx db.Tx, u *pacta.User) error {
	var lang pgtype.Text
	if u.PreferredLanguage != "" {
		lang.Valid = true
		lang.String = string(u.PreferredLanguage)
	}
	err := d.exec(tx, `
		UPDATE pacta_user SET
			admin = $2,
			super_admin = $3,
			name = $4,
			preferred_language = $5
		WHERE id = $1;
		`, u.ID, u.Admin, u.SuperAdmin, u.Name, lang)
	if err != nil {
		return fmt.Errorf("updating pacta_user writable fields: %w", err)
	}
	return nil
}

func validateUserForCreation(u *pacta.User) error {
	if u.ID != "" {
		return fmt.Errorf("user ID must be empty")
	}
	if u.AuthnID == "" {
		return fmt.Errorf("user AuthnID must not be empty")
	}
	if u.AuthnMechanism == "" {
		return fmt.Errorf("user AuthnMechanism must not be empty")
	}
	if u.EnteredEmail == "" {
		return fmt.Errorf("user EnteredEmail must not be empty")
	}
	if u.CanonicalEmail == "" {
		return fmt.Errorf("user CanonicalEmail must not be empty")
	}
	return nil
}

func rowToUser(row rowScanner) (*pacta.User, error) {
	var (
		lang  pgtype.Text
		authm string
	)

	u := &pacta.User{}
	err := row.Scan(
		&u.ID,
		&authm,
		&u.AuthnID,
		&u.EnteredEmail,
		&u.CanonicalEmail,
		&u.Admin,
		&u.SuperAdmin,
		&u.Name,
		&lang,
		&u.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning into user: %w", err)
	}
	a, err := pacta.ParseAuthnMechanism(authm)
	if err != nil {
		return nil, fmt.Errorf("parsing authn_mechanism: %w", err)
	}
	u.AuthnMechanism = a
	if lang.Valid {
		l, err := pacta.ParseLanguage(lang.String)
		if err != nil {
			return nil, fmt.Errorf("parsing user preffered_language: %w", err)
		}
		u.PreferredLanguage = l
	}
	return u, nil
}

func rowsToUsers(rows pgx.Rows) ([]*pacta.User, error) {
	return mapRows("user", rows, rowToUser)
}

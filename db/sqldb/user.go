package sqldb

import (
	"fmt"
	"strings"

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

func (d *DB) GetOrCreateUserByAuthn(tx db.Tx, authnMechanism pacta.AuthnMechanism, authnID, enteredEmail, canonicalEmail string) (*pacta.User, error) {
	var user *pacta.User
	err := d.RunOrContinueTransaction(tx, func(tx db.Tx) error {
		u, err := d.UserByAuthn(tx, authnMechanism, authnID)
		if err == nil {
			user = u
			return nil
		}
		if !db.IsNotFound(err) {
			return fmt.Errorf("looking up user by authn: %w", err)
		}
		uID, err := d.createUser(tx, &pacta.User{
			CanonicalEmail: canonicalEmail,
			EnteredEmail:   enteredEmail,
			AuthnMechanism: authnMechanism,
			AuthnID:        authnID,
		})
		if err != nil {
			return fmt.Errorf("creating user: %w", err)
		}
		u, err = d.User(tx, uID)
		if err != nil {
			return fmt.Errorf("reading back created user: %w", err)
		}
		user = u
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("running get_or_create_user txn: %w", err)
	}
	return user, nil
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

func (d *DB) QueryUsers(tx db.Tx, q *db.UserQuery) ([]*pacta.User, *db.PageInfo, error) {
	if q.Limit <= 0 {
		return nil, nil, fmt.Errorf("limit must be greater than 0, was %d", q.Limit)
	}
	offset, err := offsetFromCursor(q.Cursor)
	if err != nil {
		return nil, nil, fmt.Errorf("converting cursor to offset: %w", err)
	}
	sql, args, err := userQuery(q)
	if err != nil {
		return nil, nil, fmt.Errorf("building user query: %w", err)
	}
	rows, err := d.query(tx, sql, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("executing user query: %w", err)
	}
	us, err := rowsToUsers(rows)
	if err != nil {
		return nil, nil, fmt.Errorf("getting users from rows: %w", err)
	}
	// This will incorrectly say "yes there are more results" if we happen to hit the actual limit, but
	// that's a pretty small performance loss.
	hasNextPage := len(us) == q.Limit
	cursor := offsetToCursor(offset + len(us))
	return us, &db.PageInfo{HasNextPage: hasNextPage, Cursor: db.Cursor(cursor)}, nil
}

func userQuery(q *db.UserQuery) (string, []any, error) {
	args := &queryArgs{}
	selectFrom := `SELECT ` + userSelectColumns + ` FROM pacta_user`
	where := userQueryWheresToSQL(q.Wheres, args)
	order := userQuerySortsToSQL(q.Sorts)
	limit := fmt.Sprintf("LIMIT %d", q.Limit)
	offset := ""
	if q.Cursor != "" {
		o, err := offsetFromCursor(q.Cursor)
		if err != nil {
			return "", nil, fmt.Errorf("extracting offset from cursor in audit-log query: %w", err)
		}
		offset = fmt.Sprintf("OFFSET %d", o)
	}
	sql := fmt.Sprintf("%s %s %s %s %s;", selectFrom, where, order, limit, offset)
	return sql, args.values, nil
}

func userQuerySortsToSQL(ss []*db.UserQuerySort) string {
	sorts := []string{}
	for _, s := range ss {
		v := " DESC"
		if s.Ascending {
			v = " ASC"
		}
		sorts = append(sorts, fmt.Sprintf("pacta_user.%s %s", s.By, v))
	}
	// Forces a deterministic sort for pagination.
	sorts = append(sorts, "pacta_user.id ASC")
	return "ORDER BY " + strings.Join(sorts, ", ")
}

func userQueryWheresToSQL(qs []*db.UserQueryWhere, args *queryArgs) string {
	wheres := []string{}
	for _, q := range qs {
		if q.NameOrEmailLike != "" {
			wheres = append(wheres,
				fmt.Sprintf(
					`name ILIKE ('%%' || %[1]s || '%%')
					OR
					canonical_email ILIKE ('%%' || %[1]s || '%%')`,
					args.add(q.NameOrEmailLike)))
		}
	}
	if len(wheres) == 0 {
		return ""
	}
	return "WHERE " + strings.Join(wheres, " AND ")
}

func (d *DB) createUser(tx db.Tx, u *pacta.User) (pacta.UserID, error) {
	if err := validateUserForCreation(u); err != nil {
		return "", fmt.Errorf("validating user for creation: %w", err)
	}
	var pl pgtype.Text
	if u.PreferredLanguage != "" {
		pl.Valid = true
		pl.String = string(u.PreferredLanguage)
	}
	id := pacta.UserID(d.randomID(userIDNamespace))
	err := d.RunOrContinueTransaction(tx, func(db.Tx) error {
		err := d.exec(tx, `
			INSERT INTO pacta_user 
				(id, authn_mechanism, authn_id, entered_email, canonical_email, admin, super_admin, name, preferred_language)
				VALUES
				($1, $2, $3, $4, $5, $6, $7, $8, $9);
			`, id, u.AuthnMechanism, u.AuthnID, u.EnteredEmail, u.CanonicalEmail, false, false, u.Name, pl)
		if err != nil {
			return fmt.Errorf("creating pacta_user row for %q: %w", id, err)
		}
		_, err = d.createOwner(tx, &pacta.Owner{User: &pacta.User{ID: id}})
		if err != nil {
			return fmt.Errorf("creating owner: %w", err)
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("creating user: %w", err)
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

func (d *DB) DeleteUser(tx db.Tx, id pacta.UserID) ([]pacta.BlobURI, error) {
	buris := []pacta.BlobURI{}
	err := d.RunOrContinueTransaction(tx, func(db.Tx) error {
		userOwnerID, err := d.GetOwnerForUser(tx, id)
		if err != nil {
			if !db.IsNotFound(err) {
				return fmt.Errorf("getting owner for user: %w", err)
			}
		} else {
			newBuris, err := d.DeleteOwner(tx, userOwnerID)
			if err != nil {
				return fmt.Errorf("deleting owner: %w", err)
			}
			buris = append(buris, newBuris...)
		}
		err = d.exec(tx, `DELETE FROM initiative_invitation WHERE used_by_user_id = $1;`, id)
		if err != nil {
			return fmt.Errorf("deleting initiative_invitation rows: %w", err)
		}
		err = d.exec(tx, `UPDATE portfolio_initiative_membership SET added_by_user_id = NULL WHERE added_by_user_id = $1;`, id)
		if err != nil {
			return fmt.Errorf("clearing portfolio_initiative_membership.added_by_user_id: %w", err)
		}
		err = d.exec(tx, `DELETE FROM pacta_user WHERE id = $1;`, id)
		if err != nil {
			return fmt.Errorf("deleting actual user: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("performing user deletion: %w", err)
	}
	return buris, nil
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

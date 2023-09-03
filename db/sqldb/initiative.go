package sqldb

import (
	"fmt"
	"regexp"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
	"github.com/jackc/pgx/v5"
)

const initiativeSelectColumns = `
	initiative.id,
	initiative.name,
	initiative.affiliation,
	initiative.public_description,
	initiative.internal_description,
	initiative.requires_invitation_to_join,
	initiative.is_accepting_new_members,
	initiative.is_accepting_new_portfolios,
	initiative.pacta_version_id,
	initiative.language,
	initiative.created_at`

func (d *DB) Initiative(tx db.Tx, id pacta.InitiativeID) (*pacta.Initiative, error) {
	rows, err := d.query(tx, `
		SELECT `+initiativeSelectColumns+`
		FROM initiative 
		WHERE id = $1;`, id)
	if err != nil {
		return nil, fmt.Errorf("querying initiative: %w", err)
	}
	pvs, err := rowsToInitiatives(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to initiatives: %w", err)
	}
	return exactlyOne("initiative", id, pvs)
}

func (d *DB) Initiatives(tx db.Tx, ids []pacta.InitiativeID) (map[pacta.InitiativeID]*pacta.Initiative, error) {
	rows, err := d.query(tx, `
		SELECT `+initiativeSelectColumns+`
		FROM initiative 
		WHERE id IN `+createWhereInFmt(len(ids))+`;`, idsToInterface(ids)...)
	if err != nil {
		return nil, fmt.Errorf("querying initiatives: %w", err)
	}
	pvs, err := rowsToInitiatives(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to initiatives: %w", err)
	}
	result := make(map[pacta.InitiativeID]*pacta.Initiative, len(pvs))
	for _, pv := range pvs {
		result[pv.ID] = pv
	}
	return result, nil
}

func (d *DB) AllInitiatives(tx db.Tx) ([]*pacta.Initiative, error) {
	rows, err := d.query(tx, `
		SELECT `+initiativeSelectColumns+`
		FROM initiative;`)
	if err != nil {
		return nil, fmt.Errorf("querying initiatives: %w", err)
	}
	pvs, err := rowsToInitiatives(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to initiatives: %w", err)
	}
	return pvs, nil
}

func (d *DB) CreateInitiative(tx db.Tx, i *pacta.Initiative) error {
	if err := validateInitiativeForCreation(i); err != nil {
		return fmt.Errorf("validating initiative for creation: %w", err)
	}
	err := d.exec(tx, `
		INSERT INTO initiative 
			(id, name, affiliation, public_description, internal_description, requires_invitation_to_join, is_accepting_new_members, is_accepting_new_portfolios, pacta_version_id, language)
			VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`,
		i.ID, i.Name, i.Affiliation, i.PublicDescription, i.InternalDescription, i.RequiresInvitationToJoin, i.IsAcceptingNewMembers, i.IsAcceptingNewPortfolios, i.PACTAVersion.ID, i.Language)
	if err != nil {
		return fmt.Errorf("creating initiative: %w", err)
	}
	return nil
}

func (d *DB) UpdateInitiative(tx db.Tx, id pacta.InitiativeID, mutations ...db.UpdateInitiativeFn) error {
	err := d.RunOrContinueTransaction(tx, func(db.Tx) error {
		init, err := d.Initiative(tx, id)
		if err != nil {
			return fmt.Errorf("reading initiative: %w", err)
		}
		for i, m := range mutations {
			err := m(init)
			if err != nil {
				return fmt.Errorf("running %d-th mutation: %w", i, err)
			}
		}
		err = d.putInitiative(tx, init)
		if err != nil {
			return fmt.Errorf("putting initiative: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("updating initiative: %w", err)
	}
	return nil
}

func (d *DB) DeleteInitiative(tx db.Tx, id pacta.InitiativeID) error {
	err := d.RunOrContinueTransaction(tx, func(db.Tx) error {
		// TODO(grady) add owner deletions here - where the initiative is the owner of the asset/blob
		// TODO(grady) do snapshot deletions here
		err := d.exec(tx, `DELETE FROM initiative_invitation WHERE initiative_id = $1;`, id)
		if err != nil {
			return fmt.Errorf("deleting initiative_invitations: %w", err)
		}
		err = d.exec(tx, `DELETE FROM initiative_user_relationship WHERE initiative_id = $1;`, id)
		if err != nil {
			return fmt.Errorf("deleting initiative_user_relationships: %w", err)
		}
		err = d.exec(tx, `DELETE FROM portfolio_initiative_membership WHERE initiative_id = $1;`, id)
		if err != nil {
			return fmt.Errorf("deleting portfolio_initiative_memberships: %w", err)
		}
		err = d.exec(tx, `DELETE FROM initiative WHERE id = $1;`, id)
		if err != nil {
			return fmt.Errorf("deleting initiative: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("performing initiative deletion: %w", err)
	}
	return nil
}

func rowToInitiative(row rowScanner) (*pacta.Initiative, error) {
	var (
		pvid pacta.PACTAVersionID
		lang string
	)

	i := &pacta.Initiative{}
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Affiliation,
		&i.PublicDescription,
		&i.InternalDescription,
		&i.RequiresInvitationToJoin,
		&i.IsAcceptingNewMembers,
		&i.IsAcceptingNewPortfolios,
		&pvid,
		&lang,
		&i.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning into pacta_version: %w", err)
	}
	if pvid != "" {
		i.PACTAVersion = &pacta.PACTAVersion{ID: pvid}
	}
	l, err := pacta.ParseLanguage(lang)
	if err != nil {
		return nil, fmt.Errorf("parsing pacta_version language: %w", err)
	}
	i.Language = l
	return i, nil
}

func rowsToInitiatives(rows pgx.Rows) ([]*pacta.Initiative, error) {
	return allRows("initiative", rows, rowToInitiative)
}

func (db *DB) putInitiative(tx db.Tx, i *pacta.Initiative) error {
	err := db.exec(tx, `
		UPDATE initiative SET
			name = $2,
			affiliation = $3, 
			public_description = $4,
			internal_description = $5,
			requires_invitation_to_join = $6,
			is_accepting_new_members = $7,
			is_accepting_new_portfolios = $8,
			pacta_version_id = $9,
			language = $10
		WHERE id = $1;
		`, i.ID, i.Name, i.Affiliation, i.PublicDescription, i.InternalDescription, i.RequiresInvitationToJoin, i.IsAcceptingNewMembers, i.IsAcceptingNewPortfolios, i.PACTAVersion.ID, i.Language)
	if err != nil {
		return fmt.Errorf("updating initiative writable fields: %w", err)
	}
	return nil
}

// TODO(grady) move this to the Resolver-equivalent layer when that exists.
// Unlike other mechanisms for creating IDs, initiative IDs are user-specified.
var initiativeIDRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func validateInitiativeForCreation(i *pacta.Initiative) error {
	if i.ID == "" {
		return fmt.Errorf("initiative id must be present")
	}
	if !initiativeIDRegex.MatchString(string(i.ID)) {
		return fmt.Errorf("initiative id must match %q", initiativeIDRegex.String())
	}
	if !i.CreatedAt.IsZero() {
		return fmt.Errorf("initiative created_at must be zero")
	}
	if i.Language == "" {
		return fmt.Errorf("initiative language must be present")
	}
	if i.Name == "" {
		return fmt.Errorf("initiative name must be present")
	}
	if i.PACTAVersion == nil || i.PACTAVersion.ID == "" {
		return fmt.Errorf("initiative pacta_version must be nil")
	}
	return nil
}

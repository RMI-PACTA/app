package sqldb

import (
	"errors"
	"fmt"
	"time"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (d *DB) PortfolioGroup(tx db.Tx, id pacta.PortfolioGroupID) (*pacta.PortfolioGroup, error) {
	pgs, err := d.PortfolioGroups(tx, []pacta.PortfolioGroupID{id})
	if err != nil {
		return nil, fmt.Errorf("querying portfolio_groups: %w", err)
	}
	return exactlyOneFromMap("portfolioGroup", id, pgs)
}

const portfolioGroupQueryStanza = `
	SELECT
		portfolio_group.id, 
		portfolio_group.owner_id, 
		portfolio_group.name, 
		portfolio_group.description, 
		portfolio_group.created_at,
		portfolio_group_membership.portfolio_id,
		portfolio_group_membership.created_at
	FROM portfolio_group
	LEFT JOIN portfolio_group_membership 
	ON portfolio_group_membership.portfolio_group_id = portfolio_group.id
`

func (d *DB) PortfolioGroups(tx db.Tx, ids []pacta.PortfolioGroupID) (map[pacta.PortfolioGroupID]*pacta.PortfolioGroup, error) {
	if len(ids) == 0 {
		return make(map[pacta.PortfolioGroupID]*pacta.PortfolioGroup), nil
	}
	ids = dedupeIDs(ids)
	rows, err := d.query(tx, portfolioGroupQueryStanza+`
		WHERE portfolio_group.id IN `+createWhereInFmt(len(ids))+`;`, idsToInterface(ids)...)
	if err != nil {
		return nil, fmt.Errorf("querying portfolio_groups: %w", err)
	}
	return rowsToPortfolioGroups(rows)
}

func (d *DB) PortfolioGroupsByOwner(tx db.Tx, ownerID pacta.OwnerID) ([]*pacta.PortfolioGroup, error) {
	rows, err := d.query(tx, portfolioGroupQueryStanza+` WHERE portfolio_group.owner_id = $1;`, ownerID)
	if err != nil {
		return nil, fmt.Errorf("querying portfolio_groups: %w", err)
	}
	// Note the map interface here is ~required in the deserialization process to track multiple memberships,
	// so we're not just converting to a map and back.
	asMap, err := rowsToPortfolioGroups(rows)
	if err != nil {
		return nil, fmt.Errorf("converting rows to portfolio groups: %w", err)
	}
	return valuesFromMap(asMap), nil
}

func (d *DB) CreatePortfolioGroup(tx db.Tx, p *pacta.PortfolioGroup) (pacta.PortfolioGroupID, error) {
	if err := validatePortfolioGroupForCreation(p); err != nil {
		return "", fmt.Errorf("validating portfolio_group for creation: %w", err)
	}
	p.ID = pacta.PortfolioGroupID(d.randomID("pfgp"))
	err := d.exec(tx, `
		INSERT INTO portfolio_group
			(id, owner_id, name, description)
			VALUES
			($1, $2, $3, $4);`,
		p.ID, p.Owner.ID, p.Name, p.Description)
	if err != nil {
		return "", fmt.Errorf("creating portfolio_group: %w", err)
	}
	return p.ID, nil
}

func (d *DB) UpdatePortfolioGroup(tx db.Tx, id pacta.PortfolioGroupID, mutations ...db.UpdatePortfolioGroupFn) error {
	err := d.RunOrContinueTransaction(tx, func(db.Tx) error {
		init, err := d.PortfolioGroup(tx, id)
		if err != nil {
			return fmt.Errorf("reading portfolioGroup: %w", err)
		}
		for i, m := range mutations {
			err := m(init)
			if err != nil {
				return fmt.Errorf("running %d-th mutation: %w", i, err)
			}
		}
		err = d.putPortfolioGroup(tx, init)
		if err != nil {
			return fmt.Errorf("putting portfolioGroup: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("updating portfolioGroup: %w", err)
	}
	return nil
}

func (d *DB) DeletePortfolioGroup(tx db.Tx, id pacta.PortfolioGroupID) error {
	return d.RunOrContinueTransaction(tx, func(tx db.Tx) error {
		err := d.exec(tx, `DELETE FROM portfolio_group_membership WHERE portfolio_group_id = $1;`, id)
		if err != nil {
			return fmt.Errorf("deleting portfolio_group_memberships: %w", err)
		}
		err = d.exec(tx, `DELETE FROM portfolio_group WHERE id = $1;`, id)
		if err != nil {
			return fmt.Errorf("deleting portfolio_group: %w", err)
		}
		return nil
	})
}

type portfolioGroupRow struct {
	PortfolioGroup                      *pacta.PortfolioGroup
	PortfolioGroupMembershipPortfolioID pacta.PortfolioID
	PortfolioGroupMembershipCreatedAt   time.Time
}

func rowToPortfolioGroupRow(row rowScanner) (*portfolioGroupRow, error) {
	p := &portfolioGroupRow{PortfolioGroup: &pacta.PortfolioGroup{Owner: &pacta.Owner{}}}
	mi := pgtype.Text{}
	ca := pgtype.Timestamptz{}
	err := row.Scan(
		&p.PortfolioGroup.ID,
		&p.PortfolioGroup.Owner.ID,
		&p.PortfolioGroup.Name,
		&p.PortfolioGroup.Description,
		&p.PortfolioGroup.CreatedAt,
		&mi,
		&ca,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning into portfolio_group row: %w", err)
	}
	if mi.Valid {
		p.PortfolioGroupMembershipPortfolioID = pacta.PortfolioID(mi.String)
		p.PortfolioGroupMembershipCreatedAt = ca.Time
	}
	return p, nil
}

func rowsToPortfolioGroups(rows pgx.Rows) (map[pacta.PortfolioGroupID]*pacta.PortfolioGroup, error) {
	pgRows, err := mapRows("portfolioGroup", rows, rowToPortfolioGroupRow)
	if err != nil {
		return nil, fmt.Errorf("translating rows to portfolio_groups: %w", err)
	}
	result := make(map[pacta.PortfolioGroupID]*pacta.PortfolioGroup)
	for _, row := range pgRows {
		id := row.PortfolioGroup.ID
		if _, ok := result[id]; !ok {
			result[id] = &pacta.PortfolioGroup{
				ID:          row.PortfolioGroup.ID,
				Owner:       row.PortfolioGroup.Owner,
				Name:        row.PortfolioGroup.Name,
				Description: row.PortfolioGroup.Description,
				CreatedAt:   row.PortfolioGroup.CreatedAt,
			}
		}
		if row.PortfolioGroupMembershipPortfolioID != "" {
			result[id].Members = append(result[id].Members, &pacta.PortfolioGroupMembership{
				Portfolio: &pacta.Portfolio{
					ID: row.PortfolioGroupMembershipPortfolioID,
				},
				CreatedAt: row.PortfolioGroupMembershipCreatedAt,
			})
		}
	}
	return result, nil
}

func (d *DB) putPortfolioGroup(tx db.Tx, p *pacta.PortfolioGroup) error {
	err := d.exec(tx, `
		UPDATE portfolio_group SET
			owner_id = $2,
			name = $3, 
			description = $4
		WHERE id = $1;
		`, p.ID, p.Owner.ID, p.Name, p.Description)
	if err != nil {
		return fmt.Errorf("updating portfolioGroup writable fields: %w", err)
	}
	return nil
}

func (d *DB) CreatePortfolioGroupMembership(tx db.Tx, pgID pacta.PortfolioGroupID, pID pacta.PortfolioID) error {
	err := d.exec(tx, `
		INSERT INTO portfolio_group_membership
			(portfolio_group_id, portfolio_id)
			VALUES
			($1, $2)
		ON CONFLICT (portfolio_group_id, portfolio_id) DO NOTHING;`,
		pgID, pID)
	if err != nil {
		return fmt.Errorf("creating portfolio_group_membership: %w", err)
	}
	return nil
}

func (d *DB) DeletePortfolioGroupMembership(tx db.Tx, pgID pacta.PortfolioGroupID, pID pacta.PortfolioID) error {
	err := d.exec(tx, `
		DELETE FROM portfolio_group_membership
		WHERE portfolio_group_id = $1 AND portfolio_id = $2;`,
		pgID, pID)
	if err != nil {
		return fmt.Errorf("deleting portfolio_group_membership: %w", err)
	}
	return nil
}

func validatePortfolioGroupForCreation(p *pacta.PortfolioGroup) error {
	if p.ID != "" {
		return errors.New("portfolio_group id must be empty")
	}
	if p.Owner == nil || p.Owner.ID == "" {
		return errors.New("portfolio_group must contain a non-nil owner with a present ID")
	}
	if p.Name == "" {
		return errors.New("portfolio_group name must be present")
	}
	if !p.CreatedAt.IsZero() {
		return fmt.Errorf("portfolio_group created_at must be zero")
	}
	return nil
}

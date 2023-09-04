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

func (d *DB) PortfolioGroups(tx db.Tx, ids []pacta.PortfolioGroupID) (map[pacta.PortfolioGroupID]*pacta.PortfolioGroup, error) {
	rows, err := d.query(tx, `
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
		WHERE portfolio_goup.id IN `+createWhereInFmt(len(ids))+`;`, idsToInterface(ids)...)
	if err != nil {
		return nil, fmt.Errorf("querying portfolio_groups: %w", err)
	}
	return rowsToPortfolioGroups(rows)
}

func (d *DB) CreatePortfolioGroup(tx db.Tx, p *pacta.PortfolioGroup) error {
	if err := validatePortfolioGroupForCreation(p); err != nil {
		return fmt.Errorf("validating portfolio_group for creation: %w", err)
	}
	p.ID = pacta.PortfolioGroupID(d.randomID("pfgp"))
	err := d.exec(tx, `
		INSERT INTO portfolio_group
			(id, owner_id, name, description)
			VALUES
			($1, $2, $3, $4);`,
		p.ID, p.Owner.ID, p.Name, p.Description)
	if err != nil {
		return fmt.Errorf("creating portfolio_group: %w", err)
	}
	return nil
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
	err := d.exec(tx, `DELETE FROM portfolio_group_membership WHERE portfolio_group_id = $1;`, id)
	if err != nil {
		return fmt.Errorf("deleting portfolio_group_memberships: %w", err)
	}
	return nil
}

type portfolioGroupRow struct {
	ID                       pacta.PortfolioGroupID
	Owner                    *pacta.Owner
	Name                     string
	Description              string
	CreatedAt                time.Time
	PortfolioMemberID        pacta.PortfolioID
	PortfolioMemberCreatedAt time.Time
}

func rowToPortfolioGroupRow(row rowScanner) (*portfolioGroupRow, error) {
	p := &portfolioGroupRow{Owner: &pacta.Owner{}}
	mi := pgtype.Text{}
	ca := pgtype.Timestamptz{}
	err := row.Scan(
		&p.ID,
		&p.Owner.ID,
		&p.Name,
		&p.Description,
		&p.CreatedAt,
		&mi,
		&ca,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning into portfolio_group row: %w", err)
	}
	if mi.Valid {
		p.PortfolioMemberID = pacta.PortfolioID(mi.String)
		p.PortfolioMemberCreatedAt = ca.Time
	}
	return p, nil
}

func rowsToPortfolioGroups(rows pgx.Rows) (map[pacta.PortfolioGroupID]*pacta.PortfolioGroup, error) {
	pgrs, err := allRows("portfolioGroup", rows, rowToPortfolioGroupRow)
	if err != nil {
		return nil, fmt.Errorf("translating rows to portfolio_groups: %w", err)
	}
	result := make(map[pacta.PortfolioGroupID]*pacta.PortfolioGroup)
	for _, pgr := range pgrs {
		if _, ok := result[pgr.ID]; !ok {
			result[pgr.ID] = &pacta.PortfolioGroup{
				ID:          pgr.ID,
				Owner:       pgr.Owner,
				Name:        pgr.Name,
				Description: pgr.Description,
				CreatedAt:   pgr.CreatedAt,
			}
		}
		if pgr.PortfolioMemberID != "" {
			result[pgr.ID].Members = append(result[pgr.ID].Members, &pacta.PortfolioGroupMembership{
				Portfolio: &pacta.Portfolio{
					ID: pgr.PortfolioMemberID,
				},
				CreatedAt: pgr.PortfolioMemberCreatedAt,
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
			description = $4,
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
			($1, $2);`,
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

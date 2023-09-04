package sqldb

import (
	"fmt"
	"time"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
	"github.com/jackc/pgx/v5"
)

func (d *DB) PortfolioInitiativeMembershipsByPortfolio(tx db.Tx, pid pacta.PortfolioID) ([]*pacta.PortfolioInitiativeMembership, error) {
	rows, err := d.query(tx, `
		SELECT portfolio_id, initiative_id, created_at
		FROM portfolio_initiative_membership 
		WHERE portfolio_id = $1;`, pid)
	if err != nil {
		return nil, fmt.Errorf("querying portfolio_initiative_membership: %w", err)
	}
	pims, err := rowsToPortfolioInitiativeMemberships(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to portfolio_initiative_membership: %w", err)
	}
	return pims, nil
}

func (d *DB) PortfolioInitiativeMembershipsByInitiative(tx db.Tx, iid pacta.InitiativeID) ([]*pacta.PortfolioInitiativeMembership, error) {
	rows, err := d.query(tx, `
		SELECT portfolio_id, initiative_id, created_at
		FROM portfolio_initiative_membership 
		WHERE initiative_id = $1;`, iid)
	if err != nil {
		return nil, fmt.Errorf("querying portfolio_initiative_membership: %w", err)
	}
	pims, err := rowsToPortfolioInitiativeMemberships(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to portfolio_initiative_membership: %w", err)
	}
	return pims, nil
}

func (d *DB) CreatePortfolioInitiativeMembership(tx db.Tx, pim *pacta.PortfolioInitiativeMembership) error {
	if err := validatePortfolioInitiativeMembershipForCreate(pim); err != nil {
		return fmt.Errorf("validating portfolio_initiative_membership for creation: %w", err)
	}
	createdAt := time.Now()
	err := d.exec(tx, `
		INSERT INTO portfolio_initiative_membership	
			(portfolio_id, initiative_id, added_by_user_id)
			VALUES
			($1, $2, $3)
		ON CONFLICT DO NOTHING;`,
		pim.Portfolio.ID, pim.Initiative.ID, createdAt, pim.AddedBy.ID)
	if err != nil {
		return fmt.Errorf("creating portfolio_initiative_membership: %w", err)
	}
	return nil
}

func (d *DB) DeletePortfolioInitiativeMembership(tx db.Tx, pid pacta.PortfolioID, iid pacta.InitiativeID) error {
	err := d.exec(tx, `
		DELETE FROM portfolio_initiative_membership	
		WHERE portfolio_id = $1 AND initiative_id = $2;`, pid, iid)
	if err != nil {
		return fmt.Errorf("deleting portfolio_initiative_membership: %w", err)
	}
	return nil
}

func rowToPortfolioInitiativeMembership(row rowScanner) (*pacta.PortfolioInitiativeMembership, error) {
	var addedByUserID *string
	m := &pacta.PortfolioInitiativeMembership{
		Portfolio:  &pacta.Portfolio{},
		Initiative: &pacta.Initiative{},
	}
	err := row.Scan(
		&m.Portfolio.ID,
		&m.Initiative.ID,
		&m.CreatedAt,
		&addedByUserID,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning into portfolio_initiative_membership: %w", err)
	}
	if addedByUserID != nil {
		m.AddedBy = &pacta.User{ID: pacta.UserID(*addedByUserID)}
	}
	return m, nil
}

func rowsToPortfolioInitiativeMemberships(rows pgx.Rows) ([]*pacta.PortfolioInitiativeMembership, error) {
	return allRows("portfolio_initiaitve_membership", rows, rowToPortfolioInitiativeMembership)
}

func validatePortfolioInitiativeMembershipForCreate(pim *pacta.PortfolioInitiativeMembership) error {
	if pim.AddedBy == nil || pim.AddedBy.ID == "" {
		return fmt.Errorf("portfolio initiative membership must be instantiated with a creating user")
	}
	if pim.Portfolio == nil || pim.Portfolio.ID == "" {
		return fmt.Errorf("portfolio initiative membership must be instantiated with a portfolio")
	}
	if pim.Initiative == nil || pim.Initiative.ID == "" {
		return fmt.Errorf("portfolio initiative membership must be instantiated with an initiative")
	}
	return nil
}
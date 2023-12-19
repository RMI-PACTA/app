package sqldb

import (
	"errors"
	"fmt"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func portfolioQueryStanza(where string) string {
	return fmt.Sprintf(`
	SELECT
		portfolio.id,
		portfolio.owner_id,
		portfolio.name,
		portfolio.description,
		portfolio.created_at,
		portfolio.holdings_date,
		portfolio.blob_id,
		portfolio.admin_debug_enabled,
		portfolio.number_of_rows,
		ARRAY_AGG(portfolio_group_membership.portfolio_group_id),
		ARRAY_AGG(portfolio_group_membership.created_at)
	FROM portfolio
	LEFT JOIN portfolio_group_membership 
	ON portfolio_group_membership.portfolio_id = portfolio.id
	%s
	GROUP BY portfolio.id;`, where)
}

func (d *DB) Portfolio(tx db.Tx, id pacta.PortfolioID) (*pacta.Portfolio, error) {
	rows, err := d.query(tx, portfolioQueryStanza(`WHERE id = $1`), id)
	if err != nil {
		return nil, fmt.Errorf("querying portfolio: %w", err)
	}
	pvs, err := rowsToPortfolios(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to portfolios: %w", err)
	}
	return exactlyOne("portfolio", id, pvs)
}

func (d *DB) Portfolios(tx db.Tx, ids []pacta.PortfolioID) (map[pacta.PortfolioID]*pacta.Portfolio, error) {
	if len(ids) == 0 {
		return make(map[pacta.PortfolioID]*pacta.Portfolio), nil
	}
	ids = dedupeIDs(ids)
	rows, err := d.query(tx, portfolioQueryStanza(`WHERE id IN `+createWhereInFmt(len(ids))), idsToInterface(ids)...)
	if err != nil {
		return nil, fmt.Errorf("querying portfolios: %w", err)
	}
	pvs, err := rowsToPortfolios(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to portfolios: %w", err)
	}
	result := make(map[pacta.PortfolioID]*pacta.Portfolio, len(pvs))
	for _, pv := range pvs {
		result[pv.ID] = pv
	}
	return result, nil
}

func (d *DB) PortfoliosByOwner(tx db.Tx, ownerID pacta.OwnerID) ([]*pacta.Portfolio, error) {
	rows, err := d.query(tx, portfolioQueryStanza(`WHERE owner_id = $1`), ownerID)
	if err != nil {
		return nil, fmt.Errorf("querying portfolios: %w", err)
	}
	pvs, err := rowsToPortfolios(rows)
	if err != nil {
		return nil, fmt.Errorf("translating rows to portfolios: %w", err)
	}
	return pvs, nil
}

func (d *DB) CreatePortfolio(tx db.Tx, p *pacta.Portfolio) (pacta.PortfolioID, error) {
	if err := validatePortfolioForCreation(p); err != nil {
		return "", fmt.Errorf("validating portfolio for creation: %w", err)
	}
	hd, err := encodeHoldingsDate(p.HoldingsDate)
	if err != nil {
		return "", fmt.Errorf("validating holdings date: %w", err)
	}
	p.ID = pacta.PortfolioID(d.randomID("pflo"))
	err = d.exec(tx, `
		INSERT INTO portfolio 
			(id, owner_id, name, description, holdings_date, blob_id, admin_debug_enabled, number_of_rows)
			VALUES
			($1, $2, $3, $4, $5, $6, $7, $8);`,
		p.ID, p.Owner.ID, p.Name, p.Description, hd, p.Blob.ID, p.AdminDebugEnabled, p.NumberOfRows)
	if err != nil {
		return "", fmt.Errorf("creating portfolio: %w", err)
	}
	return p.ID, nil
}

func (d *DB) UpdatePortfolio(tx db.Tx, id pacta.PortfolioID, mutations ...db.UpdatePortfolioFn) error {
	err := d.RunOrContinueTransaction(tx, func(db.Tx) error {
		init, err := d.Portfolio(tx, id)
		if err != nil {
			return fmt.Errorf("reading portfolio: %w", err)
		}
		for i, m := range mutations {
			err := m(init)
			if err != nil {
				return fmt.Errorf("running %d-th mutation: %w", i, err)
			}
		}
		err = d.putPortfolio(tx, init)
		if err != nil {
			return fmt.Errorf("putting portfolio: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("updating portfolio: %w", err)
	}
	return nil
}

func (d *DB) DeletePortfolio(tx db.Tx, id pacta.PortfolioID) ([]pacta.BlobURI, error) {
	buris := []pacta.BlobURI{}
	err := d.RunOrContinueTransaction(tx, func(db.Tx) error {
		p, err := d.Portfolio(tx, id)
		if err != nil {
			return fmt.Errorf("reading portfolio: %w", err)
		}
		err = d.exec(tx, `DELETE FROM portfolio_group_membership WHERE portfolio_id = $1;`, id)
		if err != nil {
			return fmt.Errorf("deleting portfolio_group_memberships: %w", err)
		}
		err = d.exec(tx, `DELETE FROM portfolio_initiative_membership WHERE portfolio_id = $1;`, id)
		if err != nil {
			return fmt.Errorf("deleting portfolio_initiative_memberships: %w", err)
		}
		err = d.exec(tx, `DELETE FROM portfolio_snapshot WHERE portfolio_id = $1;`, id)
		if err != nil {
			return fmt.Errorf("deleting portfolio_snapshots: %w", err)
		}
		err = d.exec(tx, `DELETE FROM portfolio WHERE id = $1;`, id)
		if err != nil {
			return fmt.Errorf("deleting portfolio: %w", err)
		}
		buri, err := d.DeleteBlob(tx, p.Blob.ID)
		if err != nil {
			return fmt.Errorf("deleting blob: %w", err)
		}
		buris = append(buris, buri)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("performing portfolio deletion: %w", err)
	}
	return buris, nil
}

func rowToPortfolio(row rowScanner) (*pacta.Portfolio, error) {
	p := &pacta.Portfolio{Owner: &pacta.Owner{}, Blob: &pacta.Blob{}}
	hd := pgtype.Timestamptz{}
	mid := []pgtype.Text{}
	mca := []pgtype.Timestamptz{}
	err := row.Scan(
		&p.ID,
		&p.Owner.ID,
		&p.Name,
		&p.Description,
		&p.CreatedAt,
		&hd,
		&p.Blob.ID,
		&p.AdminDebugEnabled,
		&p.NumberOfRows,
		&mid,
		&mca,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning into portfolio row: %w", err)
	}
	p.HoldingsDate, err = decodeHoldingsDate(hd)
	if err != nil {
		return nil, fmt.Errorf("decoding holdings date: %w", err)
	}
	if len(mid) != len(mca) {
		return nil, fmt.Errorf("portfolio group membership ids and created ats must be the same length")
	}
	for i := range mid {
		if !mid[i].Valid && !mca[i].Valid {
			continue // skip nulls
		}
		if !mid[i].Valid {
			return nil, fmt.Errorf("portfolio group membership ids must be non-null")
		}
		if !mca[i].Valid {
			return nil, fmt.Errorf("portfolio group membership createdAt must be non-null")
		}
		p.MemberOf = append(p.MemberOf, &pacta.PortfolioGroupMembership{
			PortfolioGroup: &pacta.PortfolioGroup{
				ID: pacta.PortfolioGroupID(mid[i].String),
			},
			CreatedAt: mca[i].Time,
		})
	}
	return p, nil
}

func rowsToPortfolios(rows pgx.Rows) ([]*pacta.Portfolio, error) {
	return mapRows("portfolio", rows, rowToPortfolio)
}

func (db *DB) putPortfolio(tx db.Tx, p *pacta.Portfolio) error {
	hd, err := encodeHoldingsDate(p.HoldingsDate)
	if err != nil {
		return fmt.Errorf("validating holdings date: %w", err)
	}
	err = db.exec(tx, `
		UPDATE portfolio SET
			owner_id = $2,
			name = $3, 
			description = $4,
			holdings_date = $5,
			admin_debug_enabled = $6,
			number_of_rows = $7
		WHERE id = $1;
		`, p.ID, p.Owner.ID, p.Name, p.Description, hd, p.AdminDebugEnabled, p.NumberOfRows)
	if err != nil {
		return fmt.Errorf("updating portfolio writable fields: %w", err)
	}
	return nil
}

func validatePortfolioForCreation(p *pacta.Portfolio) error {
	if p.ID != "" {
		return errors.New("portfolio id must be empty")
	}
	if p.Owner == nil || p.Owner.ID == "" {
		return errors.New("portfolio must contain a non-nil owner with a present ID")
	}
	if p.Name == "" {
		return errors.New("portfolio name must be present")
	}
	if !p.CreatedAt.IsZero() {
		return fmt.Errorf("portfolio created_at must be zero")
	}
	if p.Blob == nil || p.Blob.ID == "" {
		return errors.New("portfolio must contain a non-nil blob with a present ID")
	}
	if p.NumberOfRows < 0 {
		return fmt.Errorf("portfolio number_of_rows must be non-negative")
	}
	if p.HoldingsDate == nil || p.HoldingsDate.Time.IsZero() {
		return fmt.Errorf("portfolio holdings_date must be non-nil and non-zero")
	}
	return nil
}

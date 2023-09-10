package sqldb

import (
	"fmt"
	"sort"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
	"github.com/jackc/pgx/v5/pgtype"
)

func (d *DB) CreateSnapshotOfPortfolio(tx db.Tx, pID pacta.PortfolioID) (pacta.PortfolioSnapshotID, error) {
	return d.createSnapshot(tx, pID, "", "", []pacta.PortfolioID{pID})
}

func (d *DB) CreateSnapshotOfPortfolioGroup(tx db.Tx, pgID pacta.PortfolioGroupID) (pacta.PortfolioSnapshotID, error) {
	pg, err := d.PortfolioGroup(tx, pgID)
	if err != nil {
		return "", fmt.Errorf("reading portfolio group: %w", err)
	}
	ids := []pacta.PortfolioID{}
	for _, p := range pg.Members {
		ids = append(ids, p.Portfolio.ID)
	}
	return d.createSnapshot(tx, "", pgID, "", ids)
}

func (d *DB) CreateSnapshotOfInitiative(tx db.Tx, iID pacta.InitiativeID) (pacta.PortfolioSnapshotID, error) {
	pims, err := d.PortfolioInitiativeMembershipsByInitiative(tx, iID)
	if err != nil {
		return "", fmt.Errorf("reading portfolio initiative memberships: %w", err)
	}
	ids := []pacta.PortfolioID{}
	for _, pim := range pims {
		ids = append(ids, pim.Portfolio.ID)
	}
	return d.createSnapshot(tx, "", "", iID, ids)
}

func (d *DB) PortfolioSnapshot(tx db.Tx, psID pacta.PortfolioSnapshotID) (*pacta.PortfolioSnapshot, error) {
	var pID, pgID, iID pgtype.Text
	var portfolioIDs []string
	err := d.queryRow(tx, `
		SELECT portfolio_id, portfolio_group_id, initiative_id, portfolio_ids
		FROM portfolio_snapshot
		WHERE id = $1;`,
		psID).Scan(&pID, &pgID, &iID, &portfolioIDs)
	if err != nil {
		return nil, fmt.Errorf("reading portfolio snapshot: %w", err)
	}
	ps := &pacta.PortfolioSnapshot{
		PortfolioIDs: stringsToIDs[pacta.PortfolioID](portfolioIDs),
	}
	if pID.Valid {
		ps.PortfolioID = pacta.PortfolioID(pID.String)
	}
	if pgID.Valid {
		ps.PortfolioGroupID = pacta.PortfolioGroupID(pgID.String)
	}
	if iID.Valid {
		ps.InitiatiativeID = pacta.InitiativeID(iID.String)
	}
	return ps, nil
}

func (d *DB) createSnapshot(tx db.Tx, pID pacta.PortfolioID, pgID pacta.PortfolioGroupID, iID pacta.InitiativeID, portfolioIDs []pacta.PortfolioID) (pacta.PortfolioSnapshotID, error) {
	included := make(map[string]bool)
	canonical := []string{}
	for _, id := range portfolioIDs {
		s := string(id)
		if !included[s] {
			canonical = append(canonical, s)
		}
	}
	sort.Strings(canonical)
	snapshotID := pacta.PortfolioSnapshotID(d.randomID("pfsn"))
	err := d.exec(tx, `
		INSERT INTO portfolio_snapshot
			(id, portfolio_id, portfolio_group_id, initiative_id, portfolio_ids)
			VALUES
			($1, $2, $3, $4, $5);`,
		snapshotID, strToNilable(pID), strToNilable(pgID), strToNilable(iID), canonical)
	if err != nil {
		return "", fmt.Errorf("creating portfolio_snapshot: %w", err)
	}
	return snapshotID, nil
}

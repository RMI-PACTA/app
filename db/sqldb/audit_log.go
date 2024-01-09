package sqldb

import (
	"errors"
	"fmt"
	"strings"

	"github.com/RMI/pacta/db"
	"github.com/RMI/pacta/pacta"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const auditLogIDNamespace = "al"
const auditLogSelectColumns = `
	audit_log.id,
	audit_log.action,
	audit_log.actor_type,
	audit_log.actor_id,
	audit_log.actor_owner_id,
	audit_log.primary_target_type,
	audit_log.primary_target_id,
	audit_log.primary_target_owner_id,
	audit_log.secondary_target_type,
	audit_log.secondary_target_id,
	audit_log.secondary_target_owner_id,
	audit_log.created_at
`

func (d *DB) AuditLogs(tx db.Tx, q *db.AuditLogQuery) ([]*pacta.AuditLog, *db.PageInfo, error) {
	if q.Limit <= 0 {
		return nil, nil, fmt.Errorf("limit must be greater than 0, was %d", q.Limit)
	}
	offset, err := offsetFromCursor(q.Cursor)
	if err != nil {
		return nil, nil, fmt.Errorf("converting cursor to offset: %w", err)
	}
	q, err = d.expandAuditLogQueryToAccountForMerges(tx, q)
	if err != nil {
		return nil, nil, fmt.Errorf("expanding audit_log query to account for merges: %w", err)
	}
	sql, args, err := auditLogQuery(q)
	if err != nil {
		return nil, nil, fmt.Errorf("building audit_log query: %w", err)
	}
	rows, err := d.query(tx, sql, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("executing audit_log query: %w", err)
	}
	als, err := rowsToAuditLogs(rows)
	if err != nil {
		return nil, nil, fmt.Errorf("getting audit_logs from rows: %w", err)
	}
	// This will incorrectly say "yes there are more results" if we happen to hit the actual limit, but
	// that's a pretty small performance loss.
	hasNextPage := len(als) == q.Limit
	cursor := offsetToCursor(offset + len(als))
	return als, &db.PageInfo{HasNextPage: hasNextPage, Cursor: db.Cursor(cursor)}, nil
}

func (d *DB) CreateAuditLog(tx db.Tx, a *pacta.AuditLog) (pacta.AuditLogID, error) {
	sql, args, id, err := d.buildCreateAuditLogQuery(tx, a)
	if err != nil {
		return "", err
	}
	err = d.exec(tx, sql, args...)
	if err != nil {
		return "", fmt.Errorf("creating audit_log row: %w", err)
	}
	return id, nil
}

func (d *DB) CreateAuditLogs(tx db.Tx, als []*pacta.AuditLog) error {
	if len(als) == 0 {
		return nil
	}
	if len(als) == 1 {
		_, err := d.CreateAuditLog(tx, als[0])
		return err
	}
	batch := &pgx.Batch{}
	for _, al := range als {
		sql, args, _, err := d.buildCreateAuditLogQuery(tx, al)
		if err != nil {
			return fmt.Errorf("building batch audit_log updates: %w", err)
		}
		batch.Queue(sql, args...)
	}
	if err := d.ExecBatch(tx, batch); err != nil {
		return fmt.Errorf("batch creating audit_logs: %w", err)
	}
	return nil
}

func (d *DB) buildCreateAuditLogQuery(tx db.Tx, a *pacta.AuditLog) (string, []interface{}, pacta.AuditLogID, error) {
	if err := validateAuditLogForCreation(a); err != nil {
		return "", nil, "", fmt.Errorf("validating audit_log for creation: %w", err)
	}
	id := pacta.AuditLogID(d.randomID(auditLogIDNamespace))
	ownerFn := func(o *pacta.Owner) pgtype.Text {
		if o == nil {
			return pgtype.Text{}
		}
		return pgtype.Text{String: string(o.ID), Valid: true}
	}
	var stt pgtype.Text
	if a.SecondaryTargetType != "" {
		stt.Valid = true
		stt.String = string(a.SecondaryTargetType)
	}
	sql := `
		INSERT INTO audit_log 
			(
				id, action, actor_type, actor_id, actor_owner_id,
				primary_target_type, primary_target_id, primary_target_owner_id,
				secondary_target_type, secondary_target_id, secondary_target_owner_id
			)
			VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
	`
	args := []interface{}{
		id, a.Action, a.ActorType, a.ActorID, ownerFn(a.ActorOwner),
		a.PrimaryTargetType, a.PrimaryTargetID, ownerFn(a.PrimaryTargetOwner),
		stt, a.SecondaryTargetID, ownerFn(a.SecondaryTargetOwner),
	}
	return sql, args, id, nil
}

func rowsToAuditLogs(rows pgx.Rows) ([]*pacta.AuditLog, error) {
	return mapRows("auditLog", rows, rowToAuditLog)
}

func rowToAuditLog(row rowScanner) (*pacta.AuditLog, error) {
	a := &pacta.AuditLog{}
	var actorType, primaryType string
	var actorOwner, primaryOwner pacta.OwnerID
	var secondaryType, secondaryOwner pgtype.Text
	err := row.Scan(
		&a.ID, &a.Action, &actorType, &a.ActorID, &actorOwner,
		&primaryType, &a.PrimaryTargetID, &primaryOwner,
		&secondaryType, &a.SecondaryTargetID, &secondaryOwner,
		&a.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scanning into audit_log: %w", err)
	}
	if a.ActorType, err = pacta.ParseAuditLogActorType(actorType); err != nil {
		return nil, fmt.Errorf("parsing audit_log actor_type: %w", err)
	}
	if a.PrimaryTargetType, err = pacta.ParseAuditLogTargetType(primaryType); err != nil {
		return nil, fmt.Errorf("parsing audit_log primary_target_type: %w", err)
	}
	if secondaryType.Valid {
		if a.SecondaryTargetType, err = pacta.ParseAuditLogTargetType(secondaryType.String); err != nil {
			return nil, fmt.Errorf("parsing audit_log secondary_target_type: %w", err)
		}
	}
	if secondaryOwner.Valid {
		a.SecondaryTargetOwner = &pacta.Owner{ID: pacta.OwnerID(secondaryOwner.String)}
	}
	if actorOwner != "" {
		a.ActorOwner = &pacta.Owner{ID: actorOwner}
	}
	if primaryOwner != "" {
		a.PrimaryTargetOwner = &pacta.Owner{ID: primaryOwner}
	}
	return a, nil
}

func validateAuditLogForCreation(a *pacta.AuditLog) error {
	if a.ID != "" {
		return fmt.Errorf("audit log already has an ID")
	}
	if !a.CreatedAt.IsZero() {
		return fmt.Errorf("audit log already has a CreatedAt")
	}
	if a.Action == "" {
		return fmt.Errorf("audit log missing required action")
	}
	if a.ActorType == "" {
		return fmt.Errorf("audit log missing ActorType")
	}
	if a.ActorID == "" {
		return fmt.Errorf("audit log missing ActorID")
	}
	if a.ActorOwner == nil {
		return fmt.Errorf("audit log ActorOwner is nil")
	}
	if a.ActorOwner.ID == "" {
		return fmt.Errorf("audit log ActorOwnerID is empty")
	}
	if a.PrimaryTargetType == "" {
		return fmt.Errorf("audit log missing PrimaryTargetType")
	}
	if a.PrimaryTargetID == "" {
		return fmt.Errorf("audit log missing PrimaryTargetID")
	}
	if a.PrimaryTargetOwner == nil {
		return fmt.Errorf("audit log PrimaryTargetOwner is nil")
	}
	if a.PrimaryTargetOwner.ID == "" {
		return fmt.Errorf("audit log PrimaryTargetOwnerID is empty")
	}
	return nil
}

func auditLogQuery(q *db.AuditLogQuery) (string, []any, error) {
	args := &queryArgs{}
	selectFrom := `SELECT ` + auditLogSelectColumns + ` FROM audit_log `
	where := auditLogQueryWheresToSQL(q.Wheres, args)
	if where == "" {
		return "", nil, errors.New("where clause cannot be empty in audit_log query")
	}
	order := auditLogQuerySortsToSQL(q.Sorts)
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

func auditLogQuerySortsToSQL(ss []*db.AuditLogQuerySort) string {
	sorts := []string{}
	for _, s := range ss {
		v := " DESC"
		if s.Ascending {
			v = " ASC"
		}
		sorts = append(sorts, fmt.Sprintf("audit_log.%s %s", s.By, v))
	}
	// Forces a deterministic sort for pagination.
	sorts = append(sorts, "audit_log.id ASC")
	return "ORDER BY " + strings.Join(sorts, ", ")
}

func auditLogQueryWheresToSQL(qs []*db.AuditLogQueryWhere, args *queryArgs) string {
	wheres := []string{}
	for _, q := range qs {
		if len(q.InID) > 0 {
			wheres = append(wheres, eqOrIn("audit_log.id", q.InID, args))
		}
		if len(q.InAction) > 0 {
			wheres = append(wheres, eqOrIn("audit_log.action", q.InAction, args))
		}
		if !q.MinCreatedAt.IsZero() {
			wheres = append(wheres, "audit_log.created_at >= "+args.add(q.MinCreatedAt))
		}
		if !q.MaxCreatedAt.IsZero() {
			wheres = append(wheres, "audit_log.created_at <= "+args.add(q.MaxCreatedAt))
		}
		if len(q.InActorType) > 0 {
			wheres = append(wheres, eqOrIn("audit_log.actor_type", q.InActorType, args))
		}
		if len(q.InActorID) > 0 {
			wheres = append(wheres, eqOrIn("audit_log.actor_id", q.InActorID, args))
		}
		if len(q.InActorOwnerID) > 0 {
			wheres = append(wheres, eqOrIn("audit_log.actor_owner_id", q.InActorOwnerID, args))
		}
		if len(q.InTargetType) > 0 {
			or := fmt.Sprintf("(%s OR %s)",
				eqOrIn("audit_log.primary_target_type", q.InTargetType, args),
				eqOrIn("audit_log.secondary_target_type", q.InTargetType, args),
			)
			wheres = append(wheres, or)
		}
		if len(q.InTargetID) > 0 {
			or := fmt.Sprintf("(%s OR %s)",
				eqOrIn("audit_log.primary_target_id", q.InTargetID, args),
				eqOrIn("audit_log.secondary_target_id", q.InTargetID, args),
			)
			wheres = append(wheres, or)
		}
		if len(q.InTargetOwnerID) > 0 {
			or := fmt.Sprintf("(%s OR %s)",
				eqOrIn("audit_log.primary_target_owner_id", q.InTargetOwnerID, args),
				eqOrIn("audit_log.secondary_target_owner_id", q.InTargetOwnerID, args),
			)
			wheres = append(wheres, or)
		}
	}
	if len(wheres) == 0 {
		return ""
	}
	return "WHERE " + strings.Join(wheres, " AND ")
}

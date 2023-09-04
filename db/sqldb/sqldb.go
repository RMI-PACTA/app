// Package sqldb provides a database handle backed by a PostgreSQL database.
package sqldb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/RMI/pacta/db"
	"github.com/Silicon-Ally/cryptorand"
	"github.com/Silicon-Ally/idgen"
	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type DB struct {
	db          SQL
	idGenerator *idgen.Generator
}

type SQL interface {
	DBConn
	Begin(context.Context) (pgx.Tx, error)
}

func New(sqlDB SQL) (*DB, error) {
	r := cryptorand.New()
	idg, err := idgen.New(r, idgen.WithDefaultLength(20), idgen.WithCharSet([]rune("abcdef0123456789")))
	if err != nil {
		return nil, fmt.Errorf("initializing idgen: %w", err)
	}
	return &DB{
		db:          sqlDB,
		idGenerator: idg,
	}, nil
}

type ctxtx struct {
	err error
	tx  pgx.Tx
	ctx context.Context
}

func (db *DB) Begin(ctx context.Context) (db.Tx, error) {
	tx, err := db.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	o := &ctxtx{
		tx:  tx,
		ctx: ctx,
	}
	return o, nil
}

func (db *DB) NoTxn(ctx context.Context) db.Tx {
	return &ctxtx{
		ctx: ctx,
	}
}

func (o *ctxtx) Commit() error {
	if o.tx == nil {
		return errors.New("cannot commit an operation that didn't originate from 'Begin'.")
	}
	return o.tx.Commit(o.ctx)
}

func (o *ctxtx) Rollback() error {
	if o.tx == nil {
		return errors.New("cannot rollback an operation that didn't originate from 'Begin'.")
	}
	return o.tx.Rollback(o.ctx)
}

type errRow struct {
	err error
}

func (e *errRow) Scan(_ ...interface{}) error {
	return e.err
}

type DBConn interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

func (d *DB) withConn(tx db.Tx, fn func(*ctxtx, DBConn) error) error {
	if tx == nil {
		tx = &ctxtx{ctx: context.Background()}
	}
	c, ok := tx.(*ctxtx)
	if !ok {
		return fmt.Errorf("unexpected type for transaction: %T", tx)
	}
	if c.err != nil {
		return fmt.Errorf("when attempting to do work: %w", c.err)
	}
	var dbc DBConn
	if c.tx != nil {
		dbc = c.tx
	} else {
		dbc = d.db
	}
	return fn(c, dbc)
}

func (d *DB) query(tx db.Tx, sql string, args ...interface{}) (rows pgx.Rows, err error) {
	err = d.withConn(tx, func(c *ctxtx, dbc DBConn) error {
		r, e := dbc.Query(c.ctx, sql, args...)
		rows = r
		return e
	})
	return
}

func (d *DB) queryRow(tx db.Tx, sql string, args ...interface{}) rowScanner {
	var row rowScanner
	err := d.withConn(tx, func(c *ctxtx, dbc DBConn) error {
		row = dbc.QueryRow(c.ctx, sql, args...)
		return nil
	})
	if err != nil {
		return &errRow{err: err}
	}
	return row
}

func (d *DB) exec(tx db.Tx, sql string, args ...interface{}) error {
	err := d.withConn(tx, func(c *ctxtx, dbc DBConn) error {
		_, err := dbc.Exec(c.ctx, sql, args...)
		return err
	})
	return err
}

type rowScanner interface {
	Scan(...interface{}) error
}

func (db *DB) Transactional(ctx context.Context, fn func(tx db.Tx) error) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	if err := fn(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to perform operation: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit txn: %w", err)
	}
	return nil
}

func (db *DB) RunOrContinueTransaction(in db.Tx, fn func(db.Tx) error) error {
	tx, err, commitFn, rollbackFn := db.startOrContinueTx(in)
	if err != nil {
		return fmt.Errorf("failed to start or continue txn: %w", err)
	}
	err = fn(tx)
	if err == nil {
		err = commitFn()
	}
	if err == nil {
		return nil
	}
	rollbackErr := rollbackFn()
	if rollbackErr != nil {
		err = multierror.Append(err, rollbackErr)
	}
	return fmt.Errorf("err in txn: %w", err)
}

func (db *DB) startOrContinueTx(in db.Tx) (db.Tx, error, func() error, func() error) {
	nilFn := func() error { return nil }
	var c *ctxtx
	ctx := context.Background()
	if in != nil {
		cc, ok := in.(*ctxtx)
		if !ok {
			return nil, fmt.Errorf("unexpected type for transaction: %T", in), nilFn, nilFn
		}
		c = cc
		ctx = c.ctx
	}
	if c != nil && c.tx != nil {
		return c, nil, nilFn, nilFn
	}
	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin txn: %w", err), nilFn, nilFn
	}
	commitFn := func() error {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit txn: %w", err)
		}
		return nil
	}
	rollbackFn := func() error {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("failed to rollback txn: %w", err)
		}
		return nil
	}
	return tx, nil, commitFn, rollbackFn
}

type idNamespace string

const idNamespaceIDSeparator = "."

func (db *DB) randomID(ns idNamespace) string {
	return fmt.Sprintf("%s%s%s", ns, idNamespaceIDSeparator, db.idGenerator.NewID())
}

func timeToNilable(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func strToNilable[T ~string](id T) *string {
	if id == "" {
		return nil
	}
	s := string(id)
	return &s
}

func allRows[T any](name string, rows pgx.Rows, fn func(rowScanner) (T, error)) ([]T, error) {
	defer rows.Close()
	var ts []T
	for rows.Next() {
		t, err := fn(rows)
		if err != nil {
			return nil, fmt.Errorf("converting row to %s: %w", name, err)
		}
		ts = append(ts, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("while processing %s rows: %w", name, err)
	}
	return ts, nil
}

func exactlyOne[T any, I ~string](name string, id I, ts []T) (T, error) {
	var zeroValue T
	if len(ts) == 0 {
		return zeroValue, db.NotFound(id, name)
	} else if len(ts) == 1 {
		return ts[0], nil
	} else {
		return zeroValue, fmt.Errorf("expected exactly one %s in result but got %d", name, len(ts))
	}
}

func exactlyOneFromMap[V any, K ~string](name string, id K, m map[K]V) (V, error) {
	var zeroValue V
	if len(m) > 1 {
		return zeroValue, fmt.Errorf("expected exactly one %s in result but got %d", name, len(m))
	}
	v, ok := m[id]
	if !ok {
		return zeroValue, db.NotFound(id, name)
	}
	return v, nil
}

func createWhereInFmt(n int) string {
	dollaz := make([]string, n)
	for i := 0; i < n; i++ {
		dollaz[i] = fmt.Sprintf("$%d", i+1)
	}
	return "(" + strings.Join(dollaz, " , ") + ")"
}

func idsToInterface[T ~string](in []T) []interface{} {
	out := make([]interface{}, len(in))
	for i, e := range in {
		out[i] = e
	}
	return out
}

func dedupeIDs[T ~string](in []T) []T {
	if len(in) < 2 {
		return in
	}
	result := []T{}
	seen := make(map[T]bool)
	for _, t := range in {
		if seen[t] {
			continue
		}
		result = append(result, t)
		seen[t] = true
	}
	return result
}

func validateHoldingsDate(t time.Time) (*time.Time, error) {
	// TODO: validate the properties of the holdings date (i.e. aligned to window)
	return timeToNilable(t), nil
}

func decodeHoldingsDate(t pgtype.Timestamptz) (time.Time, error) {
	// TODO: validate the properties of the holdings date (i.e. aligned to window)
	if !t.Valid {
		return time.Time{}, nil
	}
	return t.Time, nil
}

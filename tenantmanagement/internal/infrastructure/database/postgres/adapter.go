package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/infrastructure/database"
)

// Adapter wraps a Postgres connection pool and implements database.DB
type Adapter struct {
	pool *pgxpool.Pool
}

func NewAdapter(pool *pgxpool.Pool) database.DB {
	return &Adapter{pool: pool}
}

// -------------------- Exec --------------------

// Exec executes a query that doesn't return rows (INSERT, UPDATE, DELETE)
func (a *Adapter) Exec(ctx context.Context, sql string, args ...any) (database.CommandTag, error) {
	tag, err := a.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &CommandTagAdapter{tag: tag}, nil
}

// CommandTagAdapter wraps pgconn.CommandTag to implement database.CommandTag
type CommandTagAdapter struct {
	tag pgconn.CommandTag
}

func (c *CommandTagAdapter) RowsAffected() int64 {
	return c.tag.RowsAffected()
}

// -------------------- QueryRow --------------------

// RowAdapter wraps pgx.Row to implement database.Row
type RowAdapter struct {
	row pgx.Row
}

// QueryRow executes a query that returns a single row
func (a *Adapter) QueryRow(ctx context.Context, sql string, args ...any) database.Row {
	row := a.pool.QueryRow(ctx, sql, args...)
	return &RowAdapter{row: row}
}

func (r *RowAdapter) Scan(dest ...any) error {
	return r.row.Scan(dest...)
}

// -------------------- Query --------------------

// RowsAdapter wraps pgx.Rows to implement database.Rows
type RowsAdapter struct {
	rows pgx.Rows
}

// Query executes a query that returns multiple rows
func (a *Adapter) Query(ctx context.Context, sql string, args ...any) (database.Rows, error) {
	rows, err := a.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &RowsAdapter{rows: rows}, nil
}

func (r *RowsAdapter) Next() bool {
	return r.rows.Next()
}

func (r *RowsAdapter) Scan(dest ...any) error {
	return r.rows.Scan(dest...)
}

func (r *RowsAdapter) Close() {
	r.rows.Close()
}

func (r *RowsAdapter) Err() error {
	return r.rows.Err()
}

// --------------------- TRANSACTIONS --------------------

// Tx starts a new transaction and returns a database.Tx wrapped in our adapter.
func (a *Adapter) Tx(ctx context.Context) (database.Tx, error) {
    tx, err := a.pool.Begin(ctx)
    if err != nil {
        return nil, err
    }
    return &TxAdapter{tx: tx}, nil
}

// TxAdapter wraps pgx.Tx to implement database.Tx AND database.DB
type TxAdapter struct {
    tx pgx.Tx
}

func (t *TxAdapter) Commit(ctx context.Context) error {
    return t.tx.Commit(ctx)
}

func (t *TxAdapter) Rollback(ctx context.Context) error {
    return t.tx.Rollback(ctx)
}

// Tx must be implemented to satisfy the database.DB interface 
// that is often embedded in database.Tx
func (t *TxAdapter) Tx(ctx context.Context) (database.Tx, error) {
    return nil, fmt.Errorf("nested transactions not implemented")
}

// -------------------- DB Implementation for Tx --------------------

func (t *TxAdapter) Exec(ctx context.Context, sql string, args ...any) (database.CommandTag, error) {
    tag, err := t.tx.Exec(ctx, sql, args...)
    if err != nil {
        return nil, err
    }
    return &CommandTagAdapter{tag: tag}, nil
}

func (t *TxAdapter) QueryRow(ctx context.Context, sql string, args ...any) database.Row {
    row := t.tx.QueryRow(ctx, sql, args...)
    return &RowAdapter{row: row}
}

func (t *TxAdapter) Query(ctx context.Context, sql string, args ...any) (database.Rows, error) {
    rows, err := t.tx.Query(ctx, sql, args...)
    if err != nil {
        return nil, err
    }
    return &RowsAdapter{rows: rows}, nil
}
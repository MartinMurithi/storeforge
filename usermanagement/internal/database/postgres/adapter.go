package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/database"
)

// Adapter wraps a Postgres connection pool and implements database.DB
type Adapter struct {
	pool *pgxpool.Pool
}

// NewAdapter creates a new Adapter from a pgxpool.Pool
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

// QueryRow executes a query that returns a single row
func (a *Adapter) QueryRow(ctx context.Context, sql string, args ...any) database.Row {
	row := a.pool.QueryRow(ctx, sql, args...)
	return &RowAdapter{row: row}
}

// RowAdapter wraps pgx.Row to implement database.Row
type RowAdapter struct {
	row pgx.Row
}

func (r *RowAdapter) Scan(dest ...any) error {
	return r.row.Scan(dest...)
}

// -------------------- Query --------------------

// Query executes a query that returns multiple rows
func (a *Adapter) Query(ctx context.Context, sql string, args ...any) (database.Rows, error) {
	rows, err := a.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &RowsAdapter{rows: rows}, nil
}

// RowsAdapter wraps pgx.Rows to implement database.Rows
type RowsAdapter struct {
	rows pgx.Rows
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

package database

import "context"

// Row represents a single database row that can be scanned into variables.
type Row interface {
	Scan(dest ...any) error
}

// CommandTag represents the result of a DB command and provides the number of rows affected.
// CommandTag abstracts pgx.CommandTag
type CommandTag interface {
	RowsAffected() int64
}

// DB defines the minimal database operations the application needs,
// independent of any specific database (Postgres, MySQL, etc.).
type DB interface {
    QueryRow(ctx context.Context, sql string, args ...any) Row
    Query(ctx context.Context, sql string, args ...any) (Rows, error)
    Exec(ctx context.Context, sql string, args ...any) (CommandTag, error)
    Tx(ctx context.Context) (Tx, error)
}

type Tx interface{
    DB
    Commit(ctx context.Context) error
    Rollback(ctx context.Context) error
}


// Rows represents an iterator over multiple query results,
// modeled after common Go database libraries.
type Rows interface {
    Next() bool
    Scan(dest ...any) error
    Close()
    Err() error
}

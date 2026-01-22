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

// DB defines the database capabilities needed by the application.
// It can fetch a single row (QueryRow) or execute a command that changes rows (Exec).
type DB interface {
    QueryRow(ctx context.Context, sql string, args ...any) Row
	Query(ctx context.Context, sql string, args ...any)( Rows, error)
    Exec(ctx context.Context, sql string, args ...any) (CommandTag, error)
}

type Rows interface {
    Next() bool
    Scan(dest ...any) error
    Close()
    Err() error
}
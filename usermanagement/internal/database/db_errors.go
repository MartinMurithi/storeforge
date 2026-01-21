// Package dberrors provides centralized PostgreSQL error mapping for domain-friendly errors.
//
// # Usage
//
// In repositories:
//
// 	err := repo.CreateUser(ctx, user)
// 	if err != nil {
// 		return dberrors.MapPostgresError(err)
// 	}
//
// In services:
//
// 	if errors.Is(err, dberrors.ErrUserAlreadyExists) {
// 		return fmt.Errorf("email already registered: %w", err)
// 	}

package dberrors

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// ----------------------------
// Domain-level errors
// ----------------------------

// Constraint / integrity errors
var (
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrMissingRequiredField = errors.New("missing required field")
	ErrForeignKeyViolation  = errors.New("foreign key violation")
	ErrCheckViolation       = errors.New("check constraint violation")
	ErrUniqueViolation      = errors.New("unique constraint violation")
	ErrExclusionViolation   = errors.New("exclusion constraint violation")
)

// Transaction errors
var (
	ErrSerializationFailure = errors.New("transaction serialization failure")
	ErrDeadlockDetected     = errors.New("deadlock detected")
	ErrQueryCanceled        = errors.New("query canceled")
	ErrNoActiveTransaction  = errors.New("no active transaction")
	ErrReadOnlyTransaction  = errors.New("read-only transaction")
)

// Connection / system errors
var (
	ErrDatabaseConnection = errors.New("database connection failed")
	ErrDatabase           = errors.New("database error")
	ErrIOError            = errors.New("database I/O error")
)

// Syntax / permission errors
var (
	ErrSyntaxError           = errors.New("syntax or access rule violation")
	ErrInsufficientPrivilege = errors.New("insufficient privilege")
	ErrInvalidSQLState       = errors.New("invalid SQL statement")
)

// Fallback error
var ErrUnknownDatabaseError = errors.New("unknown database error")

// ----------------------------
// MapPostgresError maps pgconn.PgError to domain-friendly errors.
// Non-Postgres errors are returned as-is.
// ----------------------------
func MapPostgresError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {

		switch pgErr.Code {

		// Class 23 — Integrity Constraint Violation
		case "23505": // unique_violation
			return ErrUserAlreadyExists
		case "23502": // not_null_violation
			return ErrMissingRequiredField
		case "23503": // foreign_key_violation
			return ErrForeignKeyViolation
		case "23514": // check_violation
			return ErrCheckViolation
		case "23P01": // exclusion_violation
			return ErrExclusionViolation

		// Class 40 — Transaction Rollback
		case "40001": // serialization_failure
			return ErrSerializationFailure
		case "40P01": // deadlock_detected
			return ErrDeadlockDetected
		case "57014": // query_canceled
			return ErrQueryCanceled
		case "25P01": // no_active_sql_transaction
			return ErrNoActiveTransaction
		case "25006": // read_only_sql_transaction
			return ErrReadOnlyTransaction

		// Class 08 — Connection Exception
		case "08003", "08006", "08001", "08004": // connection issues
			return ErrDatabaseConnection

		// Class 53 / 54 — Resource / Program Limits
		case "53100", "53200", "53300": // out of memory, too many connections
			return ErrDatabase

		// Class 42 — Syntax / Access Rule Violation
		case "42601", "42000": // syntax_error / access violation
			return ErrSyntaxError
		case "42501": // insufficient privilege
			return ErrInsufficientPrivilege
		case "2D000", "2F000": // invalid transaction termination / routine exception
			return ErrInvalidSQLState

		// Class 58 / F0 — System / I/O Errors
		case "58000", "58030": // system_error, io_error
			return ErrIOError

		// Default — any other PG errors
		default:
			return ErrUnknownDatabaseError
		}
	}

	// Non-Postgres errors are returned unchanged
	return err
}

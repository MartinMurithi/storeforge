package database

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ----------------------------
// Infra-level errors
// ----------------------------

var (
	ErrNotFound        = errors.New("row not found")
	ErrUniqueViolation = errors.New("unique constraint violation")
	ErrForeignKey      = errors.New("foreign key violation")
	ErrNotNull         = errors.New("not null violation")
	ErrSerialization   = errors.New("serialization failure")
	ErrConnection      = errors.New("database connection error")
	ErrUnknown         = errors.New("unknown database error")
)

func MapPostgresError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return ErrUniqueViolation
		case "23503":
			return ErrForeignKey
		case "23502":
			return ErrNotNull
		case "40001":
			return ErrSerialization
		case "08003", "08006":
			return ErrConnection
		default:
			return ErrUnknown
		}
	}

	return err
}

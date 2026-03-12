package value_object

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type TenantID struct{ value pgtype.UUID }
type ThemeID struct{ value pgtype.UUID }

// ---------------------------------------------------------
// Constructors
// ---------------------------------------------------------

func NewTenantID(id string) (TenantID, error) {
	var u pgtype.UUID
	err := u.Scan(id) // pgtype.UUID handles string parsing natively
	if err != nil {
		return TenantID{}, fmt.Errorf("invalid tenant id: %w", err)
	}
	return TenantID{value: u}, nil
}

func NewThemeID(id string) (ThemeID, error) {
	var u pgtype.UUID
	err := u.Scan(id)
	if err != nil {
		return ThemeID{}, fmt.Errorf("invalid theme id: %w", err)
	}
	return ThemeID{value: u}, nil
}

func NewTenantIDFromUUID(u uuid.UUID) TenantID {
	return TenantID{value: pgtype.UUID{Bytes: u, Valid: true}}
}
func NewThemeIDFromUUID(u uuid.UUID) ThemeID {
	return ThemeID{value: pgtype.UUID{Bytes: u, Valid: true}}
}

// ---------------------------------------------------------
// Accessors
// ---------------------------------------------------------

func (t TenantID) Raw() pgtype.UUID { return t.value }
func (t TenantID) String() string {
	if !t.value.Valid {
		return ""
	}
	// Converts [16]byte to string format
	u := uuid.UUID(t.value.Bytes)
	return u.String()
}

// ---------------------------------------------------------
// Database Interfaces (pgx compatibility)
// ---------------------------------------------------------

// Scan implements sql.Scanner for database/sql compatibility
func (t *TenantID) Scan(value interface{}) error {
	return t.value.Scan(value)
}

// Value implements driver.Valuer for database/sql compatibility
func (t TenantID) Value() (driver.Value, error) {
	if !t.value.Valid {
		return nil, nil
	}
	return t.value.Value()
}

func (t *ThemeID) Scan(value interface{}) error { return t.value.Scan(value) }
func (t ThemeID) Value() (driver.Value, error)  { return t.value.Value() }

func (t ThemeID) Raw() pgtype.UUID { return t.value }

func (t ThemeID) String() string {
    if !t.value.Valid {
        return ""
    }
    u := uuid.UUID(t.value.Bytes)
    return u.String()
}
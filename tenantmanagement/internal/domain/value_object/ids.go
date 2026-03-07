package value_object

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

type TenantID struct{ value uuid.UUID }
type ThemeID struct{ value uuid.UUID }

// ---------------------------------------------------------
// Constructors (Keep these for Service Layer)
// ---------------------------------------------------------

func NewTenantID(id string) (TenantID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return TenantID{}, fmt.Errorf("invalid tenant id")
	}
	return TenantID{value: parsed}, nil
}

func NewThemeID(id string) (ThemeID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return ThemeID{}, fmt.Errorf("invalid theme id")
	}
	return ThemeID{value: parsed}, nil
}

func NewTenantIDFromUUID(u uuid.UUID) TenantID { return TenantID{value: u} }
func NewThemeIDFromUUID(u uuid.UUID) ThemeID   { return ThemeID{value: u} }

// ---------------------------------------------------------
// Accessors
// ---------------------------------------------------------

func (t TenantID) UUID() uuid.UUID { return t.value }
func (t TenantID) String() string  { return t.value.String() }

func (t ThemeID) UUID() uuid.UUID { return t.value }
func (t ThemeID) String() string  { return t.value.String() }

// ---------------------------------------------------------
// Database Interfaces (Required for pgx/sql Scanning)
// ---------------------------------------------------------

// Scan implements sql.Scanner to read from DB (UUID -> ValueObject)
func (t *TenantID) Scan(value interface{}) error {
	return scanUUID(&t.value, value)
}

func (t *ThemeID) Scan(value interface{}) error {
	return scanUUID(&t.value, value)
}

// Value implements driver.Valuer to write to DB (ValueObject -> UUID)
func (t TenantID) Value() (driver.Value, error) {
	return t.value.String(), nil
}

func (t ThemeID) Value() (driver.Value, error) {
	return t.value.String(), nil
}

// Internal helper to avoid code duplication
func scanUUID(target *uuid.UUID, value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case []byte:
		u, err := uuid.FromBytes(v)
		if err != nil {
			return err
		}
		*target = u
	case string:
		u, err := uuid.Parse(v)
		if err != nil {
			return err
		}
		*target = u
	default:
		return fmt.Errorf("unsupported type for UUID scan: %T", value)
	}
	return nil
}

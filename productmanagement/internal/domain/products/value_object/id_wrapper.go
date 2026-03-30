package value_object

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProductID struct{ value pgtype.UUID }
type TenantID struct { value pgtype.UUID }

// ---------------------------------------------------------
// Constructors
// ---------------------------------------------------------

func NewProductID(id string) (ProductID, error) {
	var u pgtype.UUID
	err := u.Scan(id) // pgtype.UUID handles string parsing natively
	if err != nil {
		return ProductID{}, fmt.Errorf("invalid product id: %w", err)
	}
	return ProductID{value: u}, nil
}

func NewProductIDFromUUID(u uuid.UUID) ProductID {
	return ProductID{value: pgtype.UUID{Bytes: u, Valid: true}}
}


func NewTenantID(id string) (TenantID, error) {
	var u pgtype.UUID
	err := u.Scan(id) // pgtype.UUID handles string parsing natively
	if err != nil {
		return TenantID{}, fmt.Errorf("invalid product id: %w", err)
	}
	return TenantID{value: u}, nil
}

func NewTenantIDFromUUID(u uuid.UUID) TenantID {
	return TenantID{value: pgtype.UUID{Bytes: u, Valid: true}}
}

// ---------------------------------------------------------
// Accessors
// ---------------------------------------------------------

func (t ProductID) Raw() pgtype.UUID { return t.value }
func (t ProductID) String() string {
	if !t.value.Valid {
		return ""
	}
	// Converts [16]byte to string format
	u := uuid.UUID(t.value.Bytes)
	return u.String()
}

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

func (t *TenantID) Scan(value interface{}) error { return t.value.Scan(value) }
func (t TenantID) Value() (driver.Value, error)  { return t.value.Value() }

package value_object

import (
	"github.com/google/uuid"
	apperrors "github.com/MartinMurithi/storeforge/pkg/errors"
)

type TenantID struct {
	value uuid.UUID
}

type ThemeID struct {
	value uuid.UUID
}

func NewTenantID(id string) (TenantID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return TenantID{}, apperrors.ErrInvalidID
	}
	return TenantID{value: parsed}, nil
}

func (t TenantID) UUID() uuid.UUID { return t.value }
func (t TenantID) String() string  { return t.value.String() }

func NewThemeID(id string) (ThemeID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return ThemeID{}, apperrors.ErrInvalidID
	}
	return ThemeID{value: parsed}, nil
}

func (t ThemeID) UUID() uuid.UUID { return t.value }
func (t ThemeID) String() string  { return t.value.String() }

// NewTenantIDFromUUID allows creating the Value Object directly from a generated UUID.
func NewTenantIDFromUUID(u uuid.UUID) TenantID {
	return TenantID{value: u}
}
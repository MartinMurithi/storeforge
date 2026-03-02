package entity

import (
	"time"

	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/value_object"
)

// ThemeConfig is our "BSON" equivalent.
// It allows for infinite flexibility in store design.
type ThemeConfig map[string]any

type Settings struct {
	TenantID  value_object.TenantID
	ThemeID   value_object.ThemeID
	Config    ThemeConfig // Updated to use the type alias
	Version   int         // For Optimistic Concurrency
	UpdatedAt time.Time
}


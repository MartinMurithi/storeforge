package entity

import (
	"time"

	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/value_object"
)

// ThemeConfig is our "BSON" equivalent.
// It allows for infinite flexibility in store design.
type ThemeConfig map[string]any

type Settings struct {
	ThemeID   value_object.ThemeID
	TenantID  value_object.TenantID
	Config    ThemeConfig
	Version   int         // For Optimistic Concurrency
	UpdatedAt time.Time
}


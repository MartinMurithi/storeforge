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
	Version   int // For Optimistic Concurrency
	UpdatedAt time.Time
}

// SystemSettings defines core operational rules for a tenant store.
// These settings affect business logic but NOT UI structure or design.
type SystemSettings struct {
	Currency      string `json:"currency"`       // KES, USD
	Timezone      string `json:"timezone"`       // Africa/Nairobi
	Locale        string `json:"locale"`         // en, fr
	InventoryMode string `json:"inventory_mode"` // track | untracked
	TaxEnabled    bool   `json:"tax_enabled"`
}

// Settings control runtime behavior of the StoreForge system but DO NOT
// define UI structure or visual appearance.
//
// Examples include:
//   - currency rules
//   - feature flags
//   - system-level business logic toggles
type Settings1 struct {
	Role     string         `json:"role"` // current user role context (owner/admin/staff)
	System   SystemSettings `json:"system"`
	Features FeatureFlags   `json:"features"`
}

// The theme system is based on a BASE + OVERRIDES architecture:
//
//	Base Theme (immutable, shared across tenants)
//	        +
//	Tenant Overrides (partial diff)
//	        =
//	Resolved Theme (runtime computed)
//
// This allows safe customization without duplicating full theme objects.
type ThemeConfig1 struct {
	ThemeID          string         `json:"theme_id"` // reference to base theme
	BaseThemeVersion string         `json:"base_theme_version"`
	Status           string         `json:"status"` // draft | published
	Overrides        ThemeOverrides `json:"overrides"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// ThemeOverrides represents partial modifications to a base theme.
//
// Only fields explicitly set here will override base theme values.
// Nil fields indicate inheritance from base theme.
type ThemeOverrides struct {
	Colors *ColorOverrides `json:"colors,omitempty"`
	// Typography *TypographyOverrides `json:"typography,omitempty"`
	// Layout *LayoutOverrides `json:"layout,omitempty"`
}

type ColorOverrides struct {
	Primary *string `json:"primary,omitempty"`
	Secondary *string `json:"secondary,omitempty"`
	Accent *string `json:"accent,omitempty"`
	Background *string `json:"background,omitempty"`
}

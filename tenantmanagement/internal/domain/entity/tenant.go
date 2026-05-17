package entity

import (
	"time"

	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/value_object"
)

type TenantStatus string

const (
	TenantProvisioning TenantStatus = "provisioning"
	TenantActive       TenantStatus = "active"
	TenantSuspended    TenantStatus = "suspended"
	TenantDeleted      TenantStatus = "deleted"
)

// Tenant represents a single store owner entity in the StoreForge platform.
//
// It contains only identity and routing information and MUST NOT include:
//   - UI configuration
//   - Theme settings
//   - Business logic or workflows
//
// Tenant acts as the root identity boundary for multi-tenant isolation.
type Tenant struct {
	ID           value_object.TenantID
	StoreName    string //generates slug for domain
	BusinessType string
	Slug         string
	SubDomain    string
	Domain       string
	Status       string //provisioning, active, suspended, pending deletion, deleted
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	DeletedAt    *time.Time //for soft deletes

	Settings *Settings
}

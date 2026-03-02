package entity

import (
	"time"

	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/value_object"
)

type Tenant struct {
	ID           value_object.TenantID
	StoreName    string //generates slug for domain
	BusinessType string
	Slug         string
	SubDomain    string
	Status       string //provisioning, active, suspended, pending deletion, deleted
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	DeletedAt    *time.Time //for soft deletes

	Settings *Settings
}

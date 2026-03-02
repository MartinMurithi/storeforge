package entity

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID           uuid.UUID `json:"id"`
	StoreName    string    `json:"storeName"`    //generates slug for domain
	BusinessType string    `json:"businessType"` //help select default theme
	Slug         string    `json:"slug"`
	SubDomain    string    `json:"subDomain"`
	Status       string    `json:"status"` //provisioning, active, suspended, pending deletion, deleted
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"` //for soft deletes
}

package entity

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID        uuid.UUID `json:"id"`
	StoreName string    `json:"storeName"`
	Slug      string    `json:"slug"`
	SubDomain string    `json:"subDomain"`
	Status    string    `json:"status"` //provisioning, active, suspended, pending deletion, deleted

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"` //for soft deletes

	Users []User `json:"-"` // all users in this tenant
}

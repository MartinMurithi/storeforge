package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID    `json:"id"`
	FullName    string       `json:"fullName"`
	Email        string       `json:"email"`
	Phone        string       `json:"phone"`
	PasswordHash string       `json:"password"`
	BusinessType string       `json:"businessType"` //help select default theme
	BusinessName string       `json:"businessName"` //generates slug for domain
	IsVerified   bool         `json:"isVerified"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    *time.Time    `json:"updated_at,omitempty"`
	DeletedAt    *time.Time   `json:"deleted_at,omitempty"` //for soft deletes

	Roles []Role `json:"-"`
	Tenants []Tenant `json:"-"`
}

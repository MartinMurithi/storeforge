package models

type User struct {
	ID           string `json:"id"`
	FullNames    string `json:"fullNames"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	PasswordHash string `json:"password"`
	BusinessType string `json:"businessType"` //help select default theme
	BusinessName string `json:"businessName"` //generates slug for domain
	IsVerified   bool   `json:"isVerified"`
	Roles        []Role `json:"-"`

	Tenants []Tenant // all tenants(shops) user belongs to
}

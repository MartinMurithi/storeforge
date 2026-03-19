package dto

import "time"

type CreateTenantResponseDTO struct {
	Message  string      `json:"message`
	Tenant   TenantDTO   `json:"tenant"`
	Theme    ThemeDTO    `json:"theme"`
	Settings SettingsDTO `json:"settings"`
	Token    TokenDTO    `json:"token"`
}

type TenantDTO struct {
	ID           string    `json:"id"`
	StoreName    string    `json:"store_name"`
	Slug         string    `json:"slug"`
	SubDomain    string    `json:"sub_domain"`
	Domain       string    `json:"domain"`
	BusinessType string    `json:"business_type"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type ThemeDTO struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Config      interface{} `json:"theme_config"`
	IsActive    bool        `json:"is_active"`
}

type SettingsDTO struct {
	TenantID    string      `json:"tenant_id"`
	ThemeID     string      `json:"theme_id"`
	ThemeConfig interface{} `json:"theme_config"`
	Version     int32       `json:"version"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

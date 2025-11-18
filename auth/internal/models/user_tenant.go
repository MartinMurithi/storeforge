package models

type UserTenant struct {
	UserID   string
	TenantID string
	RoleID   string

	User   User
	Tenant Tenant
	Role   Role
}

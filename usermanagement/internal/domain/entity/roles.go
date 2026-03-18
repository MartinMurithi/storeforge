package entity

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Role struct {
	ID          pgtype.UUID
	Name        string
	Slug        string
	Description string
	IsSystem    bool // if its a store owner
	Permissions []*Permission
}

const (
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleMember = "member"
	RoleViewer = "viewer"
)

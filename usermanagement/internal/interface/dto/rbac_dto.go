package dto

import (
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateRoleRequestDTO struct {
	Name          string        `json:"name"`
	Slug          string        `json:"slug"`
	Description   string        `json:"description"`
	PermissionIDs []pgtype.UUID `json:"permission_ids"`
}

func (d *CreateRoleRequestDTO) Normalize() {
	d.Name = strings.TrimSpace(d.Name)
	d.Slug = strings.ToLower(strings.TrimSpace(d.Slug))
}

type UpdateRoleRequestDTO struct {
	RoleID        pgtype.UUID        `json:"role_id" binding:"required"`
	Name          string        `json:"name" binding:"required"`
	Description   string        `json:"description"`
	PermissionIDs []pgtype.UUID `json:"permission_ids"`
}

// Normalize cleans up input to make it safe for processing.
// func (u *UpdateRoleRequestDTO) Normalize() {
// 	// 1. Trim whitespace from Name and Description
// 	u.Name = strings.TrimSpace(u.Name)
// 	u.Description = strings.TrimSpace(u.Description)

// 	// 2. Deduplicate PermissionIDs
// 	if len(u.PermissionIDs) > 0 {
// 		seen := make(map[string]struct{}, len(u.PermissionIDs))
// 		unique := make([]string, 0, len(u.PermissionIDs))

// 		for _, pid := range u.PermissionIDs {
// 			pid = strings.TrimSpace(pid)
// 			if pid == "" {
// 				continue // skip empty strings
// 			}
// 			if _, exists := seen[pid]; !exists {
// 				seen[pid] = struct{}{}
// 				unique = append(unique, pid)
// 			}
// 		}
// 		u.PermissionIDs = unique
// 	} else {
// 		u.PermissionIDs = []string{} // ensure slice is not nil
// 	}
// }

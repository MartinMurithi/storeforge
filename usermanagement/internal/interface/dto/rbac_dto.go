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

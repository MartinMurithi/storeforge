package entity

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type RolesPermissions struct {
	RoleId       pgtype.UUID
	PermissionId pgtype.UUID
	CreatedAt    time.Time
}

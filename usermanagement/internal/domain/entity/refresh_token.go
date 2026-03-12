package entity

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type RefreshToken struct {
	Id           pgtype.UUID
	UserId       pgtype.UUID //ID of logged in user
	TokenHash    string
	ExpiresAt    time.Time
	Revoked      bool
	LastRole     string
	LastTenantId pgtype.UUID
	CreatedAt    time.Time
	RevokedAt    *time.Time
}

package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserTenant struct {
	UserID   uuid.UUID `json:"userId"`
	TenantID uuid.UUID `json:"tenantId"`
	RoleID   uuid.UUID `json:"roleId"`

	JoinedAt time.Time `json:"joinedAt"`
}

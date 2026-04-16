package entity

import (
	"time"

	"github.com/google/uuid"
)

// Membership represents a user's role in a tenant (many-to-many relationship)
//owner(can do everything), admin(manage orders, settings, products, etc), member(manage products only), viewer(read only)
type UsersTenants struct {
	UserID    uuid.UUID // FK to users.id
	TenantID  uuid.UUID // FK to tenants.id
	Role      uuid.UUID    // owner, admin, editor
	CreatedAt time.Time
}

package entity

import "github.com/jackc/pgx/v5/pgtype"

// Permission represents a granular action (e.g., "orders:read")
// owner(can do everything), admin(manage orders, settings, products, etc), member(manage products only), viewer(read only)
type Permission struct {
	Id          pgtype.UUID
	Slug        string //eg edit_products
	Category    string
	Description string
}

// Permissions will be seeded directly to DB
var DefaultPermissions = []Permission{
	{Slug: "products:read", Description: "View products"},
	{Slug: "products:write", Description: "Edit products"},
	{Slug: "orders:manage", Description: "Fulfill orders"},
}

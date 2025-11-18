package models

//owner(can do everything), admin(manage orders, settings, products, etc), member(manage products only), viewer(read only)

type Permission struct {
	Name string `json:"name"` //eg edit_products
}

type Role struct {
	Name        string       `json:"name"` // e.g., "admin", "editor", "member", "viewer"
	Description string       `json:"description"`
	Permissions []Permission `json:"-"`
}

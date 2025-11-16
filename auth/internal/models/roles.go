package models

import "gorm.io/gorm"

//owner(can do everything), admin(manage orders, settings, products, etc), member(manage products only), viewer(read only)

type Permission struct {
	gorm.Model
	Name string `gorm:"unique;not null" json:"name"` //eg edit_products
}

type Role struct {
	gorm.Model
	Name        string       `gorm:"unique;not null" json:"name"` // e.g., "admin", "editor", "member", "viewer"
	Description string       `gorm:"not null" json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions" json:"-"`
}

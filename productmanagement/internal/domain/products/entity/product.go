package entity

import (
	"time"

	"github.com/MartinMurithi/storeforge/productmanagement/internal/domain/products/value_object"
)

type ProductStatus string

const (
	Draft      ProductStatus = "draft"
	Active     ProductStatus = "active"
	Archived   ProductStatus = "archived"
	OutOfStock ProductStatus = "out_of_stock"
)

// ProductProperties is our "BSON" equivalent.
// It allows for infinite flexibility in product creation.
type ProductProperties map[string]any

type Product struct {
	ID          value_object.ProductID
	TenantID    value_object.TenantID
	Name        string
	Description string
	Price       int64
	SKU         string
	Stock       int64
	Status      *ProductStatus
	Properties  *ProductProperties
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time //for soft deletes

	ProductImages []ProductImage
}

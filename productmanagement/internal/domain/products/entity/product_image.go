package entity

import (
	"time"

	"github.com/MartinMurithi/storeforge/productmanagement/internal/domain/products/value_object"
)

type ProductImage struct{
	ID value_object.ProductImageID
	ProductID value_object.ProductID
	ImageUrl string
	IsPrimary bool
	SortOrder int
	CreatedAt time.Time
	DeletedAt *time.Time
}


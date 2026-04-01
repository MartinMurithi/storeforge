package dto

import "time"

type ProductImageDTO struct {
	ID        string     `json:"id"`
	ProductID string     `json:"product_id"`
	ImageURL  string     `json:"image_url"`
	IsPrimary bool       `json:"is_primary"`
	SortOrder int32      `json:"sort_order"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type ProductDTO struct {
	ID          string                 `json:"id"`
	TenantID    string                 `json:"tenant_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Price       int64                  `json:"price"`
	Currency    string                 `json:"currency"`
	SKU         string                 `json:"sku"`
	Stock       int64                  `json:"stock"`
	Status      string                 `json:"status"`
	Images      []ProductImageDTO      `json:"images"`
	Properties  map[string]interface{} `json:"properties"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	DeletedAt   *time.Time             `json:"deleted_at,omitempty"`
}

type CreateProductResponseDTO struct {
	Message string     `json:"message"`
	Product ProductDTO `json:"product"`
}
package product

import "github.com/MartinMurithi/storeforge/productmanagement/internal/domain/products/entity"

type CreateProductRequestDTO struct {
	TenantID    string                    `json:"tenant_id"`
	UserID      string                    `json:"user_id"`
	Name        string                    `json:"name"`
	Price       int64                     `json:"price"`
	Description string                    `json:"description"`
	SKU         string                    `json:"sku"`
	Stock       int64                     `json:"stock"`
	Status      entity.ProductStatus      `json:"status"`
	Properties  *entity.ProductProperties `json:"properties"`
}

// Product image input DTO
type ProductImageInputDTO struct {
	URL       string `json:"url"`                  // required: image URL
	SortOrder int    `json:"sort_order,omitempty"` // optional, backend assigns if 0
	IsPrimary bool   `json:"is_primary,omitempty"` // optional, backend sets first image as primary if missing
}

type CreateProductResponseDTO struct {
	Product *entity.Product `json:"product"`
	Message string          `json:"message"`
}

type GetTenantProductsRequestDTO struct {
	TenantID string `json:"tenant_id"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

type UpdateProductRequestDTO struct {
	ProductID string `json:"product_id"`
	TenantID  string `json:"tenant_id"`
	UserID    string `json:"user_id"`

	Name        *string               `json:"name,omitempty"`
	Description *string               `json:"description,omitempty"`
	Price       *int64                `json:"price,omitempty"`
	Currency    *string               `json:"currency,omitempty"`
	SKU         *string               `json:"sku,omitempty"`
	Stock       *int64                `json:"stock,omitempty"`
	Status      *entity.ProductStatus `json:"status,omitempty"`

	Properties *entity.ProductProperties `json:"properties,omitempty"`

	// Images []UpdateProductImageDTO `json:"images,omitempty"`
}

type AddProductImagesRequestDTO struct {
	ProductID string `json:"product_id"`
	TenantID  string `json:"tenant_id"`
	UserID    string `json:"user_id"`

	Images []AddProductImageInputDTO `json:"images"`
}

type AddProductImageInputDTO struct {
	URL       string `json:"url"`                  // required
	SortOrder int    `json:"sort_order,omitempty"` // optional
	IsPrimary bool   `json:"is_primary,omitempty"` // optional
}

type DeleteProductImagesRequestDTO struct {
	ProductID string `json:"product_id"`
	TenantID  string `json:"tenant_id"`
	UserID    string `json:"user_id"`

	ImageIDs []string `json:"image_ids"`
}

type DeleteProductRequestDTO struct {
	ProductID string `json:"product_id"`
	TenantID  string `json:"tenant_id"`
	UserID    string `json:"user_id"`
}

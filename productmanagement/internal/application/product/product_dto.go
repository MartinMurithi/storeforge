package product

import "github.com/MartinMurithi/storeforge/productmanagement/internal/domain/products/entity"

type CreateProductRequestDTO struct {
	TenantID    string                    `json:"tenant_id"`
	Name        string                    `json:"name"`
	Price       int64                     `json:"price"`
	Description string                    `json:"description"`
	SKU         string                    `json:"sku"`
	Stock       int64                     `json:"stock"`
	Status      entity.ProductStatus      `json:"status"`
	Properties  *entity.ProductProperties `json:"properties"`
	Images      []ProductImageInputDTO    `json:"images"`
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

package product

import "github.com/MartinMurithi/storeforge/productmanagement/internal/domain/products/entity"

type CreateProductRequestDTO struct {
	Name        string
	Description string
	SKU         string
	Stock       int64
	Status      *entity.ProductStatus
	Properties  *entity.ProductProperties
}

type CreateProductResponseDTO struct {
	Product *entity.Product
	Message string
}

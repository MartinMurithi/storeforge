package mapper

import (

	productv1 "github.com/MartinMurithi/storeforge/api/protos/productmanagement/product/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/dto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapCreateProductResponse(pbRes *productv1.CreateProductResponse) dto.CreateProductResponseDTO {

	productDTO := MapGetProductResponse(
		&productv1.GetProductByIDResponse{
			Product: pbRes.Product,
		},
	)

	return dto.CreateProductResponseDTO{
		Message: pbRes.Message,
		Product: productDTO,
	}
}

func MapGetProductResponse(pbRes *productv1.GetProductByIDResponse) dto.ProductDTO {

	images := make([]dto.ProductImageDTO, 0, len(pbRes.Product.Images))

	for _, img := range pbRes.Product.Images {

		images = append(images, dto.ProductImageDTO{
			ID:        img.Id,
			ProductID: img.ProductId,
			ImageURL:  img.ImageUrl,
			IsPrimary: img.IsPrimary,
			SortOrder: img.SortOrder,
			CreatedAt: img.CreatedAt.AsTime(),
		})
	}

	return dto.ProductDTO{
		ID:          pbRes.Product.Id,
		TenantID:    pbRes.Product.TenantId,
		Name:        pbRes.Product.Name,
		Description: pbRes.Product.Description,

		Price:       pbRes.Product.Price,
		Currency:    pbRes.Product.Currency,
		SKU:         pbRes.Product.Sku,
		Stock:       pbRes.Product.Stock,

		Status:      pbRes.Product.Status,

		Images:      images,

		Properties:  pbRes.Product.Properties.AsMap(),

		CreatedAt:   pbRes.Product.CreatedAt.AsTime(),
		UpdatedAt:   pbRes.Product.UpdatedAt.AsTime(),
	}
}

func MapGetTenantProductsResponse(pbRes *productv1.GetTenantProductsResponse) []dto.ProductDTO {

	products := make([]dto.ProductDTO, 0, len(pbRes.Products))

	for _, p := range pbRes.Products {

		images := make([]dto.ProductImageDTO, 0, len(p.Images))

		for _, img := range p.Images {
			images = append(images, dto.ProductImageDTO{
				ID:        img.Id,
				ProductID: img.ProductId,
				ImageURL:  img.ImageUrl,
				IsPrimary: img.IsPrimary,
				SortOrder: img.SortOrder,
				CreatedAt: img.CreatedAt.AsTime(),
			})
		}

		var updatedAt *timestamppb.Timestamp
		if p.UpdatedAt != nil {
			t := p.UpdatedAt
			updatedAt = t
		}

		products = append(products, dto.ProductDTO{
			ID:          p.Id,
			TenantID:    p.TenantId,
			Name:        p.Name,
			Description: p.Description,

			Price:       p.Price,
			Currency:    p.Currency,
			SKU:         p.Sku,
			Stock:       p.Stock,

			Status:      p.Status,

			Images:      images,

			Properties:  p.Properties.AsMap(),

			CreatedAt:   p.CreatedAt.AsTime(),
			UpdatedAt:   updatedAt.AsTime(),
		})
	}

	return products
}
package product

import (
	"github.com/MartinMurithi/storeforge/productmanagement/internal/application/product"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/domain/products/entity"
)

// ToProtoProduct maps the Product domain entity to its Protobuf representation.
// func ToProtoProduct(t *entity.Product) *Productv1.Product {
// 	if t == nil {
// 		return nil
// 	}
// 	return &Productv1.Product{

// 	}
// }

// // ToProtoCreateProductResponse composes multiple entity mappers into a single response DTO.
// func ToProtoCreateProductResponse(data *product.CreateProductResponseDTO) *Productv1.CreateProductResponse {
// 	if data == nil || data.Product == nil {
// 		return nil
// 	}

// 	resp := &Productv1.CreateProductResponse{
// 		Product:   ToProtoProduct(data.Product),
// 		Theme:    ToProtoTheme(data.Theme),
// 		Settings: ToProtoSettings(data.Product.Settings),
// 	}


// 	return resp
// }
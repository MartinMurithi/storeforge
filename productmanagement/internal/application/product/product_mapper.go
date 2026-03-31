package product

import (
	"time"

	productv1 "github.com/MartinMurithi/storeforge/api/protos/productmanagement/product/v1"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/domain/products/entity"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToProtoProduct maps the Product domain entity to its Protobuf representation.
func ToProtoProduct(p *entity.Product) *productv1.Product {
	if p == nil {
		return nil
	}

	// Map images
	var protoImages []*productv1.ProductImage
	for _, img := range p.ProductImages {
		protoImages = append(protoImages, &productv1.ProductImage{
			Id:        img.ID.String(),
			ProductId: img.ProductID.String(),
			ImageUrl:  img.ImageUrl,
			IsPrimary: img.IsPrimary,
			SortOrder: int32(img.SortOrder),
			CreatedAt: timestamppb.New(img.CreatedAt),
		})
	}

	// Map product properties (entity.ProductProperties -> google.protobuf.Struct)
	var props *structpb.Struct
	if p.Properties != nil {
		props, _ = structpb.NewStruct(*p.Properties)
	}

	return &productv1.Product{
		Id:          p.ID.String(),
		TenantId:    p.TenantID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Sku:         p.SKU,
		Stock:       p.Stock,
		Status:      string(p.Status),
		Properties:  props,
		Images:      protoImages,
		CreatedAt:   toProtoTimestamp(&p.CreatedAt),
		UpdatedAt:   toProtoTimestamp(p.UpdatedAt),
	}
}

// toProtoTimestamp converts an optional time.Time pointer into a protobuf Timestamp.
//
// A nil time value is preserved as nil, allowing optional timestamps
// (e.g. updated_at, deleted_at) to remain unset in the wire representation.
//
// This helper avoids leaking protobuf concerns into the domain layer.
func toProtoTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

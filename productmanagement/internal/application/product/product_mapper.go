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

	var props *structpb.Struct
	if p.Properties != nil {
		propsMap := map[string]any{
			"version": p.Properties.Version,
			"data":    p.Properties.Data,
		}

		props, _ = structpb.NewStruct(propsMap)
	}

	return &productv1.Product{
		Id:          p.ID.String(),
		TenantId:    p.TenantID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Currency:    p.Currency,
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

// ToProtoPaginationMeta converts internal pagination metadata to gRPC proto.
func ToProtoPaginationMeta(meta *PaginationMeta) *productv1.PaginationMeta {
	if meta == nil {
		return nil
	}

	return &productv1.PaginationMeta{
		Page:       int32(meta.Page),
		Limit:      int32(meta.Limit),
		Total:      int32(meta.Total),
		TotalPages: int32(meta.TotalPages),
		HasNext:    meta.HasNext,
		HasPrev:    meta.HasPrev,
	}
}

// ToProtoFetchTenantProductsResponse maps service-layer products + pagination → gRPC response
func ToProtoFetchTenantProductsResponse(products []*entity.Product, meta *PaginationMeta) *productv1.GetTenantProductsResponse {
	protoProducts := make([]*productv1.Product, len(products))

	for i, u := range products {
		protoProducts[i] = ToProtoProduct(u)
	}

	return &productv1.GetTenantProductsResponse{
		Products: protoProducts,
		Meta:     ToProtoPaginationMeta(meta),
	}
}

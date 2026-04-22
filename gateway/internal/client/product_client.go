package client

import (
	"context"
	"log"
	"time"

	productv1 "github.com/MartinMurithi/storeforge/api/protos/productmanagement/product/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductClient struct {
	Service productv1.ProductServiceClient
}

func NewProductClient(addr string) *ProductClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect to Product Service: %v", err)
	}

	return &ProductClient{
		Service: productv1.NewProductServiceClient(conn),
	}
}

func (c *ProductClient) CreateProduct(ctx context.Context, req *productv1.CreateProductRequest) (*productv1.CreateProductResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.Service.CreateProduct(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *ProductClient) GetTenantProducts(ctx context.Context, req *productv1.GetTenantProductsRequest) (*productv1.GetTenantProductsResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.Service.GetTenantProducts(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *ProductClient) GetProductByID(ctx context.Context, req *productv1.GetProductByIDRequest) (*productv1.GetProductByIDResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.Service.GetProductByID(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *ProductClient) UpdateProduct(ctx context.Context, req *productv1.UpdateProductRequest) (*productv1.UpdateProductResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.Service.UpdateProduct(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *ProductClient) SoftDeleteProduct(ctx context.Context, req *productv1.SoftDeleteProductRequest) (*productv1.SoftDeleteProductResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.Service.SoftDeleteProduct(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *ProductClient) AddProductImages(ctx context.Context, req *productv1.AddProductImagesRequest) (*productv1.AddProductImagesResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.Service.AddProductImages(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *ProductClient) DeleteProductImages(ctx context.Context, req *productv1.DeleteProductImagesRequest) (*productv1.DeleteProductImagesResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.Service.DeleteProductImages(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

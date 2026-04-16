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

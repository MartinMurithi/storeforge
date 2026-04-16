package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/MartinMurithi/storeforge/productmanagement/internal/application/product/services"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/config"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/domain/products/repository"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/infrastructure/database/postgres"
	grpcclient "github.com/MartinMurithi/storeforge/productmanagement/internal/infrastructure/grpc_client"
	grpc_trans "github.com/MartinMurithi/storeforge/productmanagement/internal/infrastructure/transport/grpc"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/infrastructure/transport/grpc/handlers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	DB             *postgres.Pool
	ProductRepo    repository.IProductRepository
	ProductService *services.ProductService
	Handler        *handlers.ProductGrpcHandler
	GRPCServer     *grpc_trans.Server
}

// Init initializes the application with full dependency injection
func Init(cfg *config.Config) (*App, error) {
	// -------------------------
	// Database setup
	// -------------------------
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := postgres.Connect(ctx, &cfg.DB)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	dbAdapter := postgres.NewAdapter(pool.Pool)

	// -------------------------
	// Tenant Service Client
	// -------------------------
	// Establish the connection to Tenant Service
	tenantConn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", cfg.GRPC.TenantServerPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service client: %w", err)
	}

	// 2. Initialize the Client Wrapper
	tenantSvcClient := grpcclient.NewTenantSvcClient(tenantConn)

	// -------------------------
	// Repository
	// --------------------------
	productRepo := repository.NewProductRepository(dbAdapter)

	// -------------------------
	// Application services
	// -------------------------
	productSvc := services.NewProductService(productRepo, tenantSvcClient)

	// -------------------------
	// gRPC Handler
	// -------------------------
	productHandler := handlers.NewProductGrpcHandler(productSvc)

	// -------------------------
	// gRPC Server
	// -------------------------
	grpcSrv, err := grpc_trans.NewGRPCServer(cfg.GRPC.Port, productSvc)
	if err != nil {
		return nil, fmt.Errorf("failed to start gRPC server: %w", err)
	}

	// Expose JWKs for envoy JWT verification
	// go jwks.ServeJWKS(publicKey)

	return &App{
		DB:             pool,
		ProductRepo:    productRepo,
		ProductService: productSvc,
		Handler:        productHandler,
		GRPCServer:     grpcSrv,
	}, nil
}

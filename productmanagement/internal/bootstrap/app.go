package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/MartinMurithi/storeforge/productmanagement/internal/application/product/services"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/config"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/domain/products/repository"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/infrastructure/database/postgres"
	grpc_trans "github.com/MartinMurithi/storeforge/productmanagement/internal/infrastructure/transport/grpc"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/infrastructure/transport/grpc/handlers"
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

	// if err := postgres.RunMigrations(cfg.DB.DSN); err != nil {
	// 	return nil, fmt.Errorf("failed to run migrations: %w", err)
	// }

	dbAdapter := postgres.NewAdapter(pool.Pool)

	// -------------------------
	// Repository
	// --------------------------
	productRepo := repository.NewProductRepository(dbAdapter)

	// -------------------------
	// Application services
	// -------------------------
	productSvc := services.NewProductService(productRepo)

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

package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/application/services/tenant"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/config"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/repository"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/infrastructure/database/postgres"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/infrastructure/transport/grpc"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/infrastructure/transport/grpc/handlers"
)

type App struct {
	DB            *postgres.Pool
	TenantRepo    repository.ITenantRepository
	TenantService *tenant.TenantService
	Handler       *handlers.TenantGrpcHandler
	GRPCServer    *grpc.Server
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
	tenantRepo := repository.NewTenantRepository(dbAdapter)
	themeRepo := repository.NewThemeRepository(dbAdapter)

	// -------------------------
	// Application services
	// -------------------------
	tenantSrv := tenant.NewTenantService(tenantRepo, themeRepo)

	// -------------------------
    // gRPC Handler
    // -------------------------
    tenantHandler := handlers.NewTenantGrpcHandler(tenantSrv)

	// -------------------------
	// gRPC Server
	// -------------------------
	grpcSrv, err := grpc.NewGRPCServer(cfg.GRPC.Port, tenantSrv)
	if err != nil {
		return nil, fmt.Errorf("failed to start gRPC server: %w", err)
	}

	return &App{
		DB:            pool,
		TenantRepo:    tenantRepo,
		TenantService: tenantSrv,
		Handler: tenantHandler,
		GRPCServer:    grpcSrv,
	}, nil
}

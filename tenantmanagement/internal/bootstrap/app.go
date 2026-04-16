package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/application/services/tenant"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/config"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/repository"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/infrastructure/clients"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/infrastructure/database/postgres"
	grpc_trans "github.com/MartinMurithi/storeforge/tenantmanagement/internal/infrastructure/transport/grpc"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/infrastructure/transport/grpc/handlers"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	DB            *postgres.Pool
	TenantRepo    repository.ITenantRepository
	TenantService *tenant.TenantService
	Handler       *handlers.TenantGrpcHandler
	GRPCServer    *grpc_trans.Server
	UserConn      *grpc.ClientConn
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
	// User Service Client
	// -------------------------
	// Establish the connection to User Service
	userConn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", cfg.GRPC.UserSvcAddr), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service client: %w", err)
	}

	// 2. Initialize the Client Wrapper
	userSvcClient := clients.NewUserServiceClient(userConn)

	// -------------------------
	// Repository
	// --------------------------
	tenantRepo := repository.NewTenantRepository(dbAdapter)
	themeRepo := repository.NewThemeRepository(dbAdapter)

	// -------------------------
	// Application services
	// -------------------------
	tenantSrv := tenant.NewTenantService(tenantRepo, themeRepo, userSvcClient)

	// -------------------------
	// gRPC Handler
	// -------------------------
	tenantHandler := handlers.NewTenantGrpcHandler(tenantSrv)

	// -------------------------
	// gRPC Server
	// -------------------------
	grpcSrv, err := grpc_trans.NewGRPCServer(cfg.GRPC.Port, tenantSrv)
	if err != nil {
		userConn.Close() // Cleanup if server fails
		return nil, fmt.Errorf("failed to start gRPC server: %w", err)
	}

	return &App{
		DB:            pool,
		TenantRepo:    tenantRepo,
		TenantService: tenantSrv,
		Handler:       tenantHandler,
		GRPCServer:    grpcSrv,
	}, nil
}

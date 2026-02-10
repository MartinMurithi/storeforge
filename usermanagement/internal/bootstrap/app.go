package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/application/auth"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/application/user"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/config"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/database/postgres"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/jwks"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/repository"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/token"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/transport/grpc"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/transport/http"
)

type App struct {
	DB          *postgres.Pool
	Repo        repository.IUserRepository
	UserService *user.UserService
	AuthService *auth.AuthService
	Handler     *http.UserHandler
	JWTMaker    *token.JWTMaker
	GRPCServer  *grpc.Server
}

// Init initializes the application with full dependency injection
func Init(cfg *config.Config) (*App, error) {
	// -------------------------
	// JWT setup
	// -------------------------
	privateKey, err := token.LoadPrivateKey(cfg.JWT.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error loading JWT private key: %w", err)
	}

	publicKey, err := token.LoadPublicKey(cfg.JWT.PublicKeyPath)

	if err != nil {
		return nil, fmt.Errorf("error loading JWT public key: %w", err)
	}

	jwtMaker, err := token.NewJWTMaker(privateKey)
	if err != nil {
		return nil, fmt.Errorf("error initializing JWT maker: %w", err)
	}

	// -------------------------
	// Database setup
	// -------------------------
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := postgres.Connect(ctx, &cfg.DB)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	if err := postgres.RunMigrations(cfg.DB.DSN); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	dbAdapter := postgres.NewAdapter(pool.Pool)
	repo := repository.NewUserRepository(dbAdapter)

	// -------------------------
	// Application services
	// -------------------------
	userSrv := user.NewUserService(repo)
	authSrv := auth.NewAuthService(repo, jwtMaker)
	handler := http.NewUserHandler(userSrv, authSrv)

	// -------------------------
	// gRPC Server
	// -------------------------
	grpcSrv, err := grpc.NewGRPCServer(cfg.GRPC.Port, userSrv, authSrv)
	if err != nil {
		return nil, fmt.Errorf("failed to start gRPC server: %w", err)
	}

	// Expose JWKs for envoy JWT verification
	go jwks.ServeJWKS(publicKey)

	return &App{
		DB:          pool,
		Repo:        repo,
		UserService: userSrv,
		AuthService: authSrv,
		Handler:     handler,
		JWTMaker:    jwtMaker,
		GRPCServer:  grpcSrv,
	}, nil
}

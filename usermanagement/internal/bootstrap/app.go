package bootstrap

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/application"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/database/postgres"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/handler"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/middleware"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/repository"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/routes"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/token"
)

type App struct {
	DB       *postgres.Pool
	Repo     repository.IUserRepository
	Service  *application.UserService
	Handler  *handler.UserHandler
	Router   *gin.Engine
	JWTMaker *token.JWTMaker
}

func Init() (*App, error) {

	// --------- LOAD ENV CONFIGS ---------
	err := godotenv.Load(".env")

	if err != nil {
		log.Printf("Warning: .env not loaded, relying on system env: %v", err)
	}

	jwtPrivateKeyPath := os.Getenv("JWT_PRIVATE_KEY_PATH")

	if jwtPrivateKeyPath == "" {
		return nil, fmt.Errorf("JWT_PRIVATE_KEY_PATH is not set")
	}

	// --------- LOAD PRIVATE KEY ---------
	privateKey, err := token.LoadPrivateKey(jwtPrivateKeyPath)

	if err != nil {
		return nil, fmt.Errorf("error loading JWT_PRIVATE_KEY %w", err)
	}

	jwtMaker, err := token.NewJWTMaker(privateKey)

	if err != nil {
		return nil, fmt.Errorf("error initializing JWT Maker %w", err)
	}

	// --------- LOAD PUBLIC KEY ---------
	jwtPublicKeyPath := os.Getenv("JWT_PUBLIC_KEY_PATH")

	if jwtPublicKeyPath == "" {
		return nil, fmt.Errorf("JWT_PUBLIC_KEY_PATH is not set")
	}

	publicKey, err := token.LoadPublicKey(jwtPublicKeyPath)

	if err != nil {
		return nil, fmt.Errorf("error loading JWT_PUBLIC_KEY %w", err)
	}

	// --------- AUTH MIDDLEWARE ---------
	authMiddleware := middleware.AuthMiddleware(jwtMaker, publicKey, "storeforge-api", "auth.storeforge")

	// --------- INITIALIZE DB ---------
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = postgres.InitDB(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to init db: %s", err)
	}

	// --------- RUN DB MIGRATIONS ---------

	err = postgres.RunMigrations(os.Getenv("DATABASE_URL"))

	if err != nil {
		return nil, fmt.Errorf("failed to run migrations %w", err)
	}

	// --------- DOMAIN ---------
	pool := postgres.Get()                   // *postgres.Pool
	pgxPool := pool.Pool                     // *pgxpool.Pool
	db := postgres.NewAdapter(pgxPool)       // database.DB interface
	repo := repository.NewUserRepository(db) // IUserRepository
	srv := application.NewUserService(repo, jwtMaker)
	handler := handler.NewUserHandler(srv)
	// -------- ROUTER ---------
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	if err := router.SetTrustedProxies([]string{}); err != nil {
		log.Fatalf("failed to set trusted proxies%v", err)
	}

	router.GET("/", func(c *gin.Context) {
		fmt.Printf("ClientIP: %s\n", c.ClientIP())
	})

	routes.RegisterUserRoutes(router, handler, authMiddleware)

	return &App{
		DB:       pool,
		Repo:     repo,
		Service:  srv,
		Handler:  handler,
		Router:   router,
		JWTMaker: jwtMaker,
	}, err
}

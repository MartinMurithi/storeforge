package bootstrap

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/MartinMurithi/storeforge/auth/internal/database/config"
	"github.com/MartinMurithi/storeforge/auth/internal/handler"
	"github.com/MartinMurithi/storeforge/auth/internal/repository"
	"github.com/MartinMurithi/storeforge/auth/internal/services"
	"github.com/MartinMurithi/storeforge/auth/internal/token"
)

type App struct {
	DB       *config.Pool
	Repo     *repository.UserRepository
	Service  *services.UserService
	Handler  *handler.UserHandler
	JWTMaker *token.JWTMaker
}

func Init() (*App, error) {

	err := godotenv.Load(".env")

	if err != nil {
		log.Printf("Warning: .env not loaded, relying on system env: %v", err)
	}

	jwtPrivateKeyPath := os.Getenv("JWT_PRIVATE_KEY_PATH")

	if jwtPrivateKeyPath == "" {
		return nil, fmt.Errorf("JWT_PRIVATE_KEY_PATH is not set")
	}

	privateKey, err := token.LoadPrivateKey(jwtPrivateKeyPath)

	if err != nil {
		return nil, fmt.Errorf("error loading JWT_PRIVATE_KEY %w", err)
	}

	jwtMaker, err := token.NewJWTMaker(privateKey)

	if err != nil {
		return nil, fmt.Errorf("error initializing JWT Maker %w", err)
	}

	// Load Public Key for Token Verification
	jwtPublicKeyPath := os.Getenv("JWT_PUBLIC_KEY_PATH")

	if jwtPublicKeyPath == "" {
		return nil, fmt.Errorf("JWT_PUBLIC_KEY_PATH is not set")
	}

	_, err = token.LoadPublicKey(jwtPublicKeyPath)

	if err != nil {
		return nil, fmt.Errorf("error loading JWT_PUBLIC_KEY %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = config.InitDB(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to init db: %s", err)
	}

	db := config.Get()

	err = config.RunMigrations(os.Getenv("DATABASE_URL"))

	if err != nil {
		return nil, fmt.Errorf("failed to run migrations %w", err)
	}

	repo := repository.NewUserRepository(db)
	srv := services.NewUserService(repo, jwtMaker)
	handler := handler.NewUserHandler(srv)

	return &App{
		DB:       db,
		Repo:     repo,
		Service:  srv,
		Handler:  handler,
		JWTMaker: jwtMaker,
	}, err
}

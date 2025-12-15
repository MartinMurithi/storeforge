package bootstrap

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/MartinMurithi/storeforge/auth/internal/database/config"
	"github.com/MartinMurithi/storeforge/auth/internal/handler"
	"github.com/MartinMurithi/storeforge/auth/internal/repository"
	"github.com/MartinMurithi/storeforge/auth/internal/services"
)

type App struct {
	DB      *config.Pool
	Repo    *repository.UserRepository
	Service *services.UserService
	Handler *handler.UserHandler
	Router  *gin.Engine
}

func Init() (*App, error) {

	err := godotenv.Load(".env")

	if err != nil {
		log.Printf("Warning: .env not loaded, relying on system env: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
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
	srv := services.NewUserService(repo)
	handler := handler.NewUserHandler(srv)

	router := gin.Default()

	return &App{
		DB:      db,
		Repo:    repo,
		Service: srv,
		Handler: handler,
		Router:  router,
	}, err
}

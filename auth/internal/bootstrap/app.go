package bootstrap

import (
	"context"
	"time"

	"github.com/joho/godotenv"

	"github.com/MartinMurithi/storeforge.io/internal/database/config"
)

type App struct {
	DB *config.Pool
}

func Init() (*App, error) {


	err := godotenv.Load(".env")

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = config.InitDB(ctx)

	if err != nil {
		return nil, err
	}

	db := config.Get()

	return &App{
		DB: db,
	}, err
}

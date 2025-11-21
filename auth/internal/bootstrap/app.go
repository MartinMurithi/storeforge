package bootstrap

import (
	"context"
	"time"

	"github.com/joho/godotenv"

	"github.com/MartinMurithi/storeforge.io/internal/database"
)

type App struct {
	DB *database.Pool
}

func Init() (*App, error) {


	err := godotenv.Load(".env")

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = database.InitDB(ctx)

	if err != nil {
		return nil, err
	}

	db := database.Get()

	return &App{
		DB: db,
	}, err
}

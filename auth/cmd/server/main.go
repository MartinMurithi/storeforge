package main

import (
	"fmt"
	"os"

	"github.com/MartinMurithi/storeforge.io/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Printf("[AUTH SERVICE ENABLED]: Ready to launch in 3 ..... 2 ...... 1.......\n")

	// connect to DB for testing purposes
	err := godotenv.Load()

	if err != nil {
		fmt.Printf("failed to load env \n")
	}

	cfg := database.Config{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Port:     os.Getenv("DB_PORT"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		Database: os.Getenv("DB_DATABASE"),
	}

	db, err := database.InitDB(cfg)

	if err != nil {
		fmt.Printf("failed to connect to database %s \n", err.Error())
	}

	fmt.Printf("connected to database successfully %s \n", db.Name())
}

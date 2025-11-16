package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/MartinMurithi/storeforge.io/internal/models"
)

type Config struct {
	Host     string
	User     string
	Password string
	Port     string
	SSLMode  string
	Database string
}

var DB *gorm.DB

func InitDB(cfg Config) (*gorm.DB, error) {

	// Retrieve DSN(data source name) from environment variable with fallback
	// dsn := os.Getenv("DATABASE_URL")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Database, cfg.Port, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database : %s /n", err.Error())
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Tenant{},
		&models.UserTenant{},
		&models.Role{},
		&models.Permission{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database : %s /n", err.Error())
	}

	fmt.Printf("migrated database /n")

	DB = db

	return DB, nil

}

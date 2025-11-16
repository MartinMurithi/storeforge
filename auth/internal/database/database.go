package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/MartinMurithi/storeforge.io/internal/models"
)

// postgres://postgres:martin321!@localhost:5432/storeforge

type Config struct {
	Host     string
	User     string
	Password string
	Port     string
	SSLMode  string
	DBName   string
}

var DB *gorm.DB

func InitDB(cfg Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s SSLMode=%s dbName=%s \n", cfg.Host, cfg.User, cfg.Password, cfg.Port, cfg.SSLMode, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database : " + err.Error())
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Tenant{},
		&models.UserTenant{},
		&models.Role{},
		&models.Permission{},
		&models.Claims{},
	); err != nil {
		panic("failed to migrate database : " + err.Error())
	}

	fmt.Printf("migrated database")

	DB = db

}

package postgres

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/config"
	
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var mu sync.Mutex

func InitDB(ctx context.Context, cfg *config.DBConfig) error {
	mu.Lock()
	defer mu.Unlock()

	if dbPool != nil {
		return nil
	}

	pool, err := Connect(ctx, &config.DBConfig{})

	if err != nil {
		return err
	}

	dbPool = pool
	return nil

}

// Get returns the singleton pool
func Get() *Pool {
	return dbPool
}

// Close safely closes the pool
func Close() {
	if dbPool != nil && dbPool.Pool != nil {
		dbPool.Close()
		dbPool = nil
		log.Println("[DATABASE] Connection closed")
	}
}

// Reset is only for tests
func Reset() {
	dbPool = nil
}

func RunMigrations(databaseUrl string) error {
	m, err := migrate.New("file://migrations", databaseUrl)

	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("No new migrations to apply")
	} else {
		log.Println("Migrations applied successfully")
	}

	log.Println("Migrations applied (if any)")

	return nil
}

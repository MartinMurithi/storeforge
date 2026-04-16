package postgres

import (
	"context"
	"log"
	"sync"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/config"

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
package config

import (
	"context"
	"log"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var mu sync.Mutex

func InitDB(ctx context.Context) error {
	mu.Lock()
	defer mu.Unlock()

	if db != nil {
		return nil
	}

	pool, err := NewPool(ctx)

	if err != nil {
		return err
	}

	db = pool
	return nil

}

func Get() *Pool {
	return db
}

func Close() {
	if db != nil && db.Pool != nil {
		db.Close()
		db = nil
		log.Println("Database connection closed!")
	}
}

// for testing purpose
func Reset() {
	db = nil
}

func RunMigrations(databaseUrl string) error {
	m, err := migrate.New("file://internal/database/migrations", databaseUrl)

	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		return err
	}

	log.Println("Migrations applied (if any)")

	return nil
}

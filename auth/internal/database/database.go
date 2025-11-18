package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	// "github.com/MartinMurithi/storeforge.io/internal/models"
)

// Create a postgres connection pool for the app
type Pool struct {
	*pgxpool.Pool
}

func NewPool() (*Pool, error) {
	dsn := os.Getenv("DB_URL")

	if dsn == "" {
		dsn = "ostgres://postgres:martin321!@localhost:5432/storeforge?sslmode=disable"
		log.Printf("[DB CONFIG] : Warning, using default DSN; ensure DB URL is set for production environment.")
	}

	// Parse and configure a new connection pool
	config, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		return nil, fmt.Errorf("invalid dsn %w", err)
	}

	// customize pool settings
	config.MaxConns = 6
	config.MinIdleConns = 3
	config.MaxConnLifetime = 60 * 60 * time.Second //1 hour

	// create a new pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		return nil, fmt.Errorf("an error occurred when creating a new connection pool %w", err)
	}

	db := &Pool{pool}

	if err := db.Ping(context.Background()); err != nil{
		return nil, fmt.Errorf("database ping failed %w", err)
	}

	fmt.Printf("database connection successful!")

	return db, nil
}


// migrate creates the users table if it does not exist.
//=> to Migrations
func (d *Pool) migrate() error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL
		);
	`
	_, err := d.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	return nil
}


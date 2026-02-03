package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/config"
)

func TestConnect_Success(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://postgres:martin321!@localhost:5432/storeforge")

	ctx := context.Background()

	pool, err := Connect(ctx, &config.DBConfig{})

	if err != nil {
		t.Fatalf("expected connection to succeed, got error: %v", err)
	}

	// Ping to ensure the DB is reachable
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("expected ping to succeed, got error: %v", err)
	}

	// Clean up
	pool.Close()
}

func TestConnect_MissingDSN(t *testing.T) {
	os.Unsetenv("DATABASE_URL")

	ctx := context.Background()

	pool, err := Connect(ctx, &config.DBConfig{})

	if err == nil {
		t.Fatal("expected error due to missing DATABASE_URL, got nil")
	}
	if pool != nil {
		t.Fatal("expected pool to be nil when connection fails")
	}
}

func TestConnect_InvalidDSN(t *testing.T) {
	os.Setenv("DATABASE_URL", "invalid_dsn_here")

	ctx := context.Background()

	pool, err := Connect(ctx, &config.DBConfig{})

	if err == nil {
		t.Fatal("expected error due to invalid DSN, got nil")
	}
	if pool != nil {
		t.Fatal("expected pool to be nil when DSN is invalid")
	}
}

func TestConnect_Timeout(t *testing.T) {
	// Use a port where nothing is listening
	os.Setenv("DATABASE_URL", "postgres://postgres:martin321!@localhost:5431 /storeforge")

	ctx := context.Background()

	pool, err := Connect(ctx, &config.DBConfig{})

	if err == nil {
		t.Fatal("expected error due to unreachable DB, got nil")
	}
	if pool != nil {
		t.Fatal("expected pool to be nil when DB unreachable")
	}
}

package bootstrap

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MartinMurithi/storeforge/internal/database/config"
)

// TestInit_BadDSN_FailsImmediately ensures that when InitDB is given an invalid DSN:
// 1. It returns an error immediately instead of hanging or succeeding silently.
// 2. The global DB instance remains nil, preventing use of a partially initialized DB.
//
// The test resets the global DB state to stay isolated from other tests.
func TestInit_BadDSN_FailsImmediately(t *testing.T) {
	config.Reset()

	t.Setenv("DATABASE_URL", "postgres://user:pwd@localhost:5432/nonexistent")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := config.InitDB(ctx)

	if err == nil {
		t.Fatal("expected error with dsn(data source name), returned nil")
	}

	db := config.Get()

	if db != nil {
		t.Fatal("expected nil db after InitDB failed, returned non-nil")
	}
}

// TestInit_RespectsContextTimeout ensures InitDB fails quickly when the DB is unreachable.
// It verifies that InitDB returns an error without hanging, respects the context, and
// completes within a reasonable time (fail-fast behavior).
func TestInit_RespectsContextTimeout(t *testing.T) {
	// Point to a port with nothing listening → pgxpool would hang forever without timeout
	t.Setenv("DATABASE_URL", "postgres://user:pwd@localhost:8/nonexistent")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	start := time.Now()
	err := config.InitDB(ctx)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatalf("expected context timeout error, got nil")
	}

	if errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context.DeadlineExceeded, got: %v", err)
	}

	if elapsed > 500*time.Millisecond {
		t.Fatalf("did NOT fail fast — took %v (want <500ms)", elapsed)
	}

	t.Logf("Correctly failed in %v ← perfect", elapsed)

}

func TestInit_RepeatedInitialization(t *testing.T) {
	config.Reset()

	t.Setenv("DATABASE_URL", "postgres://postgres:martin321!@localhost:5432/storeforge")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//1st DB initialization
	err := config.InitDB(ctx)

	if err != nil {
		t.Fatalf("first init db failed %v", err)
	}

	db1 := config.Get()

	if db1 == nil {
		t.Fatalf("expected db to be non-nil after first initialization, got nil")
	}

	//2nd DB initialization
	err = config.InitDB(ctx)

	if err != nil {
		t.Fatalf("second init db failed %v", err)
	}

	db2 := config.Get()

	if db2 == nil {
		t.Fatalf("expected db to be non-nil after second initialization, got nil")
	}

	if db1 != db2 {
		t.Fatalf("expected repeated InitDB calls to return the same DB instance")
	}

}

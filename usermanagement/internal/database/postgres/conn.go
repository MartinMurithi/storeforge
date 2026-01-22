package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Create a postgres connection pool for the app
type Pool struct {
	*pgxpool.Pool
}

var db *Pool

func Connect(ctx context.Context) (*Pool, error) {

	const maxConnections = 20
	const minConnections = 4
	const minIdleConnections = 3
	const maxConnectionsLifeTime = 1 * time.Hour
	const MaxConnIdleTime = 30 * time.Minute
	const healthCheckPeriod = 1 * time.Minute
	const connectionTimeout = 10 * time.Second //fail first on bad network

	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		return nil, fmt.Errorf("Database_URL is required")
	}

	// Parse and configure a new connection pool
	config, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		return nil, fmt.Errorf("invalid DATABASE_URL %s", err)
	}

	config.MaxConns = maxConnections
	config.MinConns = minConnections
	config.MinIdleConns = minIdleConnections
	config.MaxConnLifetime = maxConnectionsLifeTime
	config.MaxConnIdleTime = MaxConnIdleTime
	config.HealthCheckPeriod = healthCheckPeriod
	config.ConnConfig.ConnectTimeout = connectionTimeout

	// Enforce TLS in production, revisit this later
	// if os.Getenv("ENV") != "development" && os.Getenv("GO_ENV") != "development" {
	//     if config.ConnConfig.TLSConfig == nil {
	//         config.ConnConfig.TLSConfig = &tls.Config{
	//             // NEVER disable verification in prod
	//             InsecureSkipVerify: false,
	//             MinVersion:         tls.VersionTLS12,
	//         }
	//     }
	// }

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// create a new pool
	pool, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		return nil, fmt.Errorf("unable to create pool %w", err)
	}

	db := &Pool{pool}

	if err := db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("initial database ping failed(DB unreachable) %w", err)
	}

	fmt.Println("[DATABASE] : database connection successful maxConns=%w, minConns=%w", config.MaxConns, config.MinConns)

	return db, nil
}


package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Create a postgres connection pool for the app
type Pool struct {
	*pgxpool.Pool
}

var dbPool *Pool

func Connect(ctx context.Context, cfg *config.DBConfig) (*Pool, error) {

	// Parse and configure a new connection pool
	config, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("invalid DATABASE_URL: %w", err)
	}

	config.MaxConns = cfg.MaxConns
	config.MinConns = cfg.MinConns
	config.MinIdleConns = cfg.MinIdleConns
	config.MaxConnLifetime = cfg.MaxConnLifetime
	config.MaxConnIdleTime = cfg.MaxConnIdleTime
	config.HealthCheckPeriod = cfg.HealthCheckPeriod
	config.ConnConfig.ConnectTimeout = cfg.ConnectTimeout

	// const maxConnections = &cfg.MaxConns
	// const minConnections = 4
	// const minIdleConnections = 3
	// const maxConnectionsLifeTime = 1 * time.Hour
	// const MaxConnIdleTime = 30 * time.Minute
	// const healthCheckPeriod = 1 * time.Minute
	// const connectionTimeout = 10 * time.Second //fail first on bad network

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

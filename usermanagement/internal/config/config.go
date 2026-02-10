package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type DBConfig struct {
	DSN               string        // DATABASE_URL
	MaxConns          int32         // maximum connections in the pool
	MinConns          int32         // minimum connections in the pool
	MinIdleConns      int32         // minimum idle connections
	MaxConnLifetime   time.Duration // max lifetime of connections
	MaxConnIdleTime   time.Duration // max idle time
	HealthCheckPeriod time.Duration // how often to ping
	ConnectTimeout    time.Duration // fail fast on bad network
}

type GRPCConfig struct {
	Port int
}

type JWTConfig struct {
	PrivateKeyPath string
	PublicKeyPath  string
}

type Config struct {
	DB   DBConfig
	GRPC GRPCConfig
	JWT  JWTConfig
	Env  string
}

// getEnv returns the value of the environment variable `key`
// or `fallback` if the variable is not set.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func Load() (*Config, error) {
	maxConns, _ := strconv.Atoi(getEnv("DB_MAX_CONNS", "20"))
	minConns, _ := strconv.Atoi(getEnv("DB_MIN_CONNS", "4"))
	minIdle, _ := strconv.Atoi(getEnv("DB_MIN_IDLE", "3"))
	maxLife, _ := time.ParseDuration(getEnv("DB_MAX_LIFETIME", "1h"))
	maxIdle, _ := time.ParseDuration(getEnv("DB_MAX_IDLE_TIME", "30m"))
	healthPeriod, _ := time.ParseDuration(getEnv("DB_HEALTH_PERIOD", "1m"))
	connectTimeout, _ := time.ParseDuration(getEnv("DB_CONNECT_TIMEOUT", "10s"))

	grpcPort, _ := strconv.Atoi(getEnv("GRPC_PORT", "50051"))

	cfg := &Config{
		DB: DBConfig{
			DSN:               getEnv("DATABASE_URL", "postgres://postgres:martin321!@localhost:5432/storeforge?sslmode=disable"),
			MaxConns:          int32(maxConns),
			MinConns:          int32(minConns),
			MinIdleConns:      int32(minIdle),
			MaxConnLifetime:   maxLife,
			MaxConnIdleTime:   maxIdle,
			HealthCheckPeriod: healthPeriod,
			ConnectTimeout:    connectTimeout,
		},
		GRPC: GRPCConfig{
			Port: grpcPort,
		},
		JWT: JWTConfig{
			PrivateKeyPath: getEnv("JWT_PRIVATE_KEY_PATH", "/home/martin-wachira/Martin/storeforge/usermanagement/internal/keys/jwt_private.pem"),
			PublicKeyPath:  getEnv("JWT_PUBLIC_KEY_PATH", "/home/martin-wachira/Martin/storeforge/usermanagement/internal/keys/jwt_public.pem"),
		},
		Env: "prod",
	}

	// Validate required fields
	if cfg.DB.DSN == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.JWT.PrivateKeyPath == "" || cfg.JWT.PublicKeyPath == "" {
		return nil, fmt.Errorf("JWT_PRIVATE_KEY_PATH and JWT_PUBLIC_KEY_PATH are required")
	}

	return cfg, nil
}

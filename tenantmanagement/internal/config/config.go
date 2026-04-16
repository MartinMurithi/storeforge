package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/MartinMurithi/storeforge/pkg/env"
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
	Port        int
	UserSvcHost string
	UserSvcPort string
}

type Config struct {
	DB   DBConfig
	GRPC GRPCConfig
	Env  string
}

func Load() (*Config, error) {
	maxConns, _ := strconv.Atoi(env.GetEnv("DB_MAX_CONNS", "20"))
	minConns, _ := strconv.Atoi(env.GetEnv("DB_MIN_CONNS", "4"))
	minIdle, _ := strconv.Atoi(env.GetEnv("DB_MIN_IDLE", "3"))
	maxLife, _ := time.ParseDuration(env.GetEnv("DB_MAX_LIFETIME", "1h"))
	maxIdle, _ := time.ParseDuration(env.GetEnv("DB_MAX_IDLE_TIME", "30m"))
	healthPeriod, _ := time.ParseDuration(env.GetEnv("DB_HEALTH_PERIOD", "1m"))
	connectTimeout, _ := time.ParseDuration(env.GetEnv("DB_CONNECT_TIMEOUT", "10s"))

	grpcPort, _ := strconv.Atoi(env.GetEnv("GRPC_PORT", "50052"))
	userHost := env.GetEnv("USER_SVC_HOST", "user-svc")
	userPort := env.GetEnv("USER_SVC_GRPC_PORT", "50051")

	cfg := &Config{
		DB: DBConfig{
			DSN:               env.GetEnv("DATABASE_URL", "postgres://postgres:martin321!@localhost:5432/storeforge?sslmode=disable"),
			MaxConns:          int32(maxConns),
			MinConns:          int32(minConns),
			MinIdleConns:      int32(minIdle),
			MaxConnLifetime:   maxLife,
			MaxConnIdleTime:   maxIdle,
			HealthCheckPeriod: healthPeriod,
			ConnectTimeout:    connectTimeout,
		},
		GRPC: GRPCConfig{
			Port:        grpcPort,
			UserSvcHost: userHost,
			UserSvcPort: userPort,
		},
		Env: "prod",
	}

	// Validate required fields
	if cfg.DB.DSN == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

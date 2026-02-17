package config

import (
	"fmt"

	"github.com/MartinMurithi/storeforge/pkg/env"
)

type Config struct {
	PublicKeyPath string
	GrpcPort      string
}

func Load() (*Config, error) {

	grpcPort := env.GetEnv("GRPC_PORT", "50051")

	cfg := &Config{
		GrpcPort:      grpcPort,
		PublicKeyPath: env.GetEnv("JWT_PUBLIC_KEY_PATH", "/home/martin-wachira/Martin/storeforge/gateway/internal/certs/jwt_public.pem"),
	}

	// Validate required fields

	if cfg.GrpcPort == "" {
		return nil, fmt.Errorf("grpc port is required")
	}

	if cfg.PublicKeyPath == "" {
		return nil, fmt.Errorf("JWT_PUBLIC_KEY_PATH is required")
	}

	return cfg, nil
}

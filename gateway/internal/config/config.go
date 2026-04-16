package config

import (
	"fmt"

	"github.com/MartinMurithi/storeforge/pkg/env"
)

type Config struct {
	PublicKeyPath      string
	UserSvcGrpcPort    string
	TenantSvcGrpcPort  string
	ProductSvcGrpcPort string
}

func Load() (*Config, error) {

	userSvcGrpcPort := env.GetEnv("USER_SVC_GRPC_PORT", "50051")
	tenantSvcGrpcPort := env.GetEnv("TENANT_SVC_GRPC_PORT", "50052")
	productSvcGrpcPort := env.GetEnv("PRODUCT_SVC_GRPC_PORT", "50053")

	cfg := &Config{
		UserSvcGrpcPort:    userSvcGrpcPort,
		TenantSvcGrpcPort:  tenantSvcGrpcPort,
		ProductSvcGrpcPort: productSvcGrpcPort,
		PublicKeyPath:      env.GetEnv("JWT_PUBLIC_KEY_PATH", "/home/martin-wachira/Martin/storeforge/gateway/internal/certs/jwt_public.pem"),
	}

	if cfg.UserSvcGrpcPort == "" {
		return nil, fmt.Errorf("user service grpc port is required")
	}

	if cfg.TenantSvcGrpcPort == "" {
		return nil, fmt.Errorf("tenant service grpc port is required")
	}

	if cfg.ProductSvcGrpcPort == "" {
		return nil, fmt.Errorf("product service grpc port is required")
	}

	if cfg.PublicKeyPath == "" {
		return nil, fmt.Errorf("JWT_PUBLIC_KEY_PATH is required")
	}

	return cfg, nil
}

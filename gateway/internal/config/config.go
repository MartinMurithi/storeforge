package config

import (
	"fmt"

	"github.com/MartinMurithi/storeforge/pkg/env"
)

type Config struct {
	PublicKeyPath string

	UserSvcHost string
	UserSvcPort string

	TenantSvcPort string
	TenantSvcHost string

	ProductSvcPort string
	ProductSvcHost string

	GatewayPort string
}

func Load() (*Config, error) {

	userHost := env.GetEnv("USER_SVC_HOST", "user-svc")
	userPort := env.GetEnv("USER_SVC_GRPC_PORT", "50051")

	tenantSvcGrpcPort := env.GetEnv("TENANT_SVC_GRPC_PORT", "50052")
	tenantHost := env.GetEnv("TENANT_SVC_HOST", "tenant-svc")

	productSvcGrpcPort := env.GetEnv("PRODUCT_SVC_GRPC_PORT", "50053")
	productHost := env.GetEnv("PRODUCT_SVC_HOST", "product-svc")

	gatewayPort := env.GetEnv("GATEWAY_PORT", "9095")

	cfg := &Config{
		UserSvcHost: userHost,
		UserSvcPort: userPort,

		TenantSvcPort: tenantSvcGrpcPort,
		TenantSvcHost: tenantHost,

		ProductSvcPort: productSvcGrpcPort,
		ProductSvcHost: productHost,

		GatewayPort:   gatewayPort,
		PublicKeyPath: env.GetEnv("JWT_PUBLIC_KEY_PATH", "/home/martin-wachira/Martin/storeforge/certs/jwt_public.pem"),
	}

	if cfg.UserSvcPort == "" && cfg.UserSvcHost == "" {
		return nil, fmt.Errorf("user service grpc port is required")
	}

	if cfg.TenantSvcPort == "" {
		return nil, fmt.Errorf("tenant service grpc port is required")
	}

	if cfg.ProductSvcPort == "" {
		return nil, fmt.Errorf("product service grpc port is required")
	}

	if cfg.PublicKeyPath == "" {
		return nil, fmt.Errorf("JWT_PUBLIC_KEY_PATH is required")
	}

	if cfg.GatewayPort == "" {
		return nil, fmt.Errorf("gateway port is required")
	}

	return cfg, nil
}

package main

import (
	"fmt"
	"log"

	productv1 "github.com/MartinMurithi/storeforge/api/protos/productmanagement/product/v1"
	tenantv1 "github.com/MartinMurithi/storeforge/api/protos/tenantmanagement/tenant/v1"
	authv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/auth/v1"
	rbacv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/rbac/v1"
	userv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/user/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/config"
	"github.com/MartinMurithi/storeforge/gateway/internal/handlers"
	"github.com/MartinMurithi/storeforge/gateway/internal/jwt"
	"github.com/MartinMurithi/storeforge/gateway/internal/middleware"
	"github.com/MartinMurithi/storeforge/gateway/internal/router"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	cfg, err := config.Load()

	log.Printf("configs loaded: %v", cfg)

	if err != nil {
		log.Fatalf("an error occurred when lading config variables%v", err)
		return
	}

	// -------------- Usermanagement Client ------------

	userServerAddr := fmt.Sprintf("%s:%s", cfg.UserSvcHost, cfg.UserSvcPort)

	userConn, err := grpc.NewClient(userServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect to the usermanagement service %v", err)
		return
	}

	defer userConn.Close()

	// -------------- Tenant Management Client ------------
	tenantServerAddr := fmt.Sprintf("0.0.0.0:%s", cfg.TenantSvcGrpcPort)

	tenantConn, err := grpc.NewClient(tenantServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect to the tenant management service %v", err)
		return
	}

	// -------------- Product Management Client ------------
	productServerAddr := fmt.Sprintf("0.0.0.0:%s", cfg.ProductSvcGrpcPort)

	productConn, err := grpc.NewClient(productServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect to the product management service %v", err)
		return
	}

	defer productConn.Close()

	// Initialize gRPC clients USING the connection
	userClient := userv1.NewUserServiceClient(userConn)
	authClient := authv1.NewAuthServiceClient(userConn)
	tenantClient := tenantv1.NewTenantServiceClient(tenantConn)
	rbacClient := rbacv1.NewRbacServiceClient(userConn)
	productClient := productv1.NewProductServiceClient(productConn)

	// Initialize user & auth handlers
	userHandler := &handlers.UserHandler{UserClient: userClient}
	authHandler := &handlers.AuthHandler{AuthClient: authClient}
	tenantHandler := &handlers.TenantHandler{TenantClient: tenantClient}
	rbacHandler := &handlers.RbacHandler{RbacClient: rbacClient}
	productHandler := &handlers.ProductHandler{ProductClient: productClient}

	// Load the pub rsa key
	publicKey, err := jwt.LoadPublicKey(cfg.PublicKeyPath)

	if err != nil {
		log.Printf("error loading JWT public key: %v\n", err)
		return
	}
	authMiddleware := middleware.AuthMiddleware(publicKey, "storeforge-client", "storeforge-api")

	// setup router
	r := router.SetupRouter(userHandler, authHandler, tenantHandler, rbacHandler, productHandler, authMiddleware)

	log.Printf("Gateway running on port %s\n", cfg.GatewayPort)

	if err := r.Run(":" + cfg.GatewayPort); err != nil {
		log.Printf("an error when starting gateway server%v\n", err)
	}
}

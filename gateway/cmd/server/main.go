package main

import (
	"fmt"
	"log"

	tenantv1 "github.com/MartinMurithi/storeforge/api/protos/tenantmanagement/tenant/v1"
	authv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/auth/v1"
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

	serverAddr := fmt.Sprintf("0.0.0.0:%s", cfg.GrpcPort)

	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect to the usermanagement service %v", err)
		return
	}

	defer conn.Close()

	// Initialize gRPC clients USING the connection
	userClient := userv1.NewUserServiceClient(conn)
	authClient := authv1.NewAuthServiceClient(conn)
	tenantClient := tenantv1.NewTenantServiceClient(conn)

	// Initialize user & auth handlers
	userHandler := &handlers.UserHandler{UserClient: userClient}
	authHandler := &handlers.AuthHandler{AuthClient: authClient}
	tenantHandler := &handlers.TenantHandler{TenantClient: tenantClient}

	// Load the pub rsa key
	publicKey, err := jwt.LoadPublicKey(cfg.PublicKeyPath)

	if err != nil {
		log.Printf("error loading JWT public key: %v\n", err)
		return
	}
	authMiddleware := middleware.AuthMiddleware(publicKey, "storeforge-client", "storeforge-api")

	// setup router
	r := router.SetupRouter(userHandler, authHandler, tenantHandler, authMiddleware)

	log.Printf("Gateway running on port 9095\n")

	if err := r.Run(":9095"); err != nil {
		log.Printf("an error when starting gateway server%v\n", err)
	}
}

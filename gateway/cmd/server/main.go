package main

import (
	"log"

	authv1 "github.com/MartinMurithi/storeforge/api/protos/auth/v1"
	userv1 "github.com/MartinMurithi/storeforge/api/protos/user/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/handlers"
	"github.com/MartinMurithi/storeforge/gateway/internal/jwt"
	"github.com/MartinMurithi/storeforge/gateway/internal/middleware"
	"github.com/MartinMurithi/storeforge/gateway/internal/router"
	"github.com/MartinMurithi/storeforge/pkg/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("0.0.0.0:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect to the usermanagement service %v", err)
		return
	}

	defer conn.Close()

	// Initialize gRPC clients USING the connection
	userClient := userv1.NewUserServiceClient(conn)
	authClient := authv1.NewAuthServiceClient(conn)

	// Initialize user & auth handlers
	userHandler := &handlers.UserHandler{UserClient: userClient}
	authHandler := &handlers.AuthHandler{AuthClient: authClient}

	// Load the pub rsa key
	publicKeyPath := env.GetEnv("JWT_PUBLIC_KEY_PATH", "/home/martin-wachira/Martin/storeforge/gateway/internal/certs/jwt_public.pem")
	publicKey, err := jwt.LoadPublicKey(publicKeyPath)

	if err != nil {
		log.Printf("error loading JWT public key: %v\n", err)
		return
	}
	authMiddleware := middleware.AuthMiddleware(publicKey, "storeforge-client", "storeforge-api")

	// setup router
	r := router.SetupRouter(userHandler, authHandler, authMiddleware)

	log.Printf("Gateway running on port 9095\n")

	if err := r.Run(":9095"); err != nil {
		log.Printf("an error when starting gateway server%v\n", err)
	}
}

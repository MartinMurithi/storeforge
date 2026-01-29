package grpc

import (
	"fmt"
	"net"

	authapp "github.com/MartinMurithi/storeforge/usermanagement/internal/application/auth"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/application/user"
	authgrpc "github.com/MartinMurithi/storeforge/usermanagement/internal/transport/grpc/auth"
	usergrpc "github.com/MartinMurithi/storeforge/usermanagement/internal/transport/grpc/user"
	authv1 "github.com/MartinMurithi/storeforge/usermanagement/proto/auth/v1"
	userv1 "github.com/MartinMurithi/storeforge/usermanagement/proto/user/v1"

	"google.golang.org/grpc"
)

type Server struct {
	GRPCServer *grpc.Server
	Listener   net.Listener
}

// NewGRPCServer creates a gRPC server with all services and handlers registered.
func NewGRPCServer(port int, userSvc *user.UserService, authSvc *authapp.AuthService) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	// Create gRPC server
	// Add interceptors later
	grpcServer := grpc.NewServer(
		// grpc.ChainUnaryInterceptor(recoveryUnaryInterceptor(), loggingUnaryInterceptor(), authUnaryInterceptor()),
	)

	// Handlers
	authHandler := authgrpc.NewAuthGrpcHandler(authSvc)
	userHandler := usergrpc.NewUserGrpcHandler(userSvc)

	// Register services
	authv1.RegisterAuthServiceServer(grpcServer, authHandler)
	userv1.RegisterUserServiceServer(grpcServer, userHandler)

	return &Server{
		GRPCServer: grpcServer,
		Listener: lis,
	}, nil
}

// Start the gRPC server
func (s *Server) Start() error {
	fmt.Printf("gRPC server listening on %s\n", s.Listener.Addr())
	return s.GRPCServer.Serve(s.Listener)
}

// Stop gracefully stops the server
func (s *Server) Stop() {
	fmt.Println("Stopping gRPC server")
	s.GRPCServer.GracefulStop()
}


// func loggingInterceptor(
//     ctx context.Context,
//     req interface{},
//     info *grpc.UnaryServerInfo,
//     handler grpc.UnaryHandler,
// ) (interface{}, error) {
//     logger.Infof("[gRPC] %s called", info.FullMethod)
//     resp, err := handler(ctx, req)
//     if err != nil {
//         logger.Errorf("[gRPC] %s error: %v", info.FullMethod, err)
//     }
//     return resp, err
// }

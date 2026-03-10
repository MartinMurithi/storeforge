package grpc

import (
	"fmt"
	"net"

	authv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/auth/v1"
	membershipv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/membership/v1"
	userv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/user/v1"
	authapp "github.com/MartinMurithi/storeforge/usermanagement/internal/application/auth"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/application/membership"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/application/user"
	authgrpc "github.com/MartinMurithi/storeforge/usermanagement/internal/transport/grpc/auth"
	membershipgrpc "github.com/MartinMurithi/storeforge/usermanagement/internal/transport/grpc/membership"
	usergrpc "github.com/MartinMurithi/storeforge/usermanagement/internal/transport/grpc/user"

	"google.golang.org/grpc"
)

type Server struct {
	GRPCServer *grpc.Server
	Listener   net.Listener
}

// NewGRPCServer creates a gRPC server with all services and handlers registered.
func NewGRPCServer(port int, userSvc *user.UserService, authSvc *authapp.AuthService, membershipSvc *membership.MembershipService) (*Server, error) {
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
	membershipHandler := membershipgrpc.NewMembershipGrpcHandler(membershipSvc)

	// Register services
	authv1.RegisterAuthServiceServer(grpcServer, authHandler)
	userv1.RegisterUserServiceServer(grpcServer, userHandler)
	membershipv1.RegisterMembershipServiceServer(grpcServer, membershipHandler)

	return &Server{
		GRPCServer: grpcServer,
		Listener:   lis,
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

package grpc

import (
	"fmt"
	"net"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/application/user"

	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

func NewGRPCServer(port int, userSvc user.UserService) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, fmt.Errorf("an error occurred when creating a new gRPC server %s", err)
	}

	s := grpc.NewServer()

	// Register gRPC service

	// pb.RegisterUserServiceServer(s, NewUserServer(userSvc))

	return &Server{
		grpcServer: s,
		listener:   lis,
	}, nil
}

func (s *Server) Start() error {
	fmt.Printf("gRPC server listening on %s\n", s.listener.Addr())
	return s.grpcServer.Serve(s.listener)
}

func (s *Server) Stop() {
	fmt.Println("Stopping gRPC server")
	s.grpcServer.GracefulStop()
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

package grpc_trans

import (
	"fmt"
	"net"

	productv1 "github.com/MartinMurithi/storeforge/api/protos/productmanagement/product/v1"
	// "github.com/MartinMurithi/storeforge/pkg/grpcx"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/application/product/services"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/infrastructure/transport/grpc/handlers"
	"google.golang.org/grpc"
)

type Server struct {
	GRPCServer *grpc.Server
	Listener   net.Listener
}

// NewGRPCServer creates a gRPC server with all services and handlers registered.
func NewGRPCServer(port int, productSrv *services.ProductService) (*Server, error) {

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
	productHandler := handlers.NewProductGrpcHandler(productSrv)

	// Register services
	productv1.RegisterProductServiceServer(grpcServer, productHandler)

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

// // func loggingInterceptor(
// //     ctx context.Context,
// //     req interface{},
// //     info *grpc.UnaryServerInfo,
// //     handler grpc.UnaryHandler,
// // ) (interface{}, error) {
// //     logger.Infof("[gRPC] %s called", info.FullMethod)
// //     resp, err := handler(ctx, req)
// //     if err != nil {
// //         logger.Errorf("[gRPC] %s error: %v", info.FullMethod, err)
// //     }
// //     return resp, err
// // }



// func NewGRPCServer(
// 	port int,
// 	productSrv *services.ProductService,
// ) (*Server, error) {

// 	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

// 	if err != nil {
// 		return nil, fmt.Errorf(
// 			"failed to listen on port %d: %w",
// 			port,
// 			err,
// 		)
// 	}

// 	grpcServer := grpc.NewServer(

// 		grpc.ChainUnaryInterceptor(
// 			grpcx.RecoveryInterceptor(),
// 		),
// 	)

// 	productHandler := handlers.NewProductGrpcHandler(productSrv)

// 	productv1.RegisterProductServiceServer(
// 		grpcServer,
// 		productHandler,
// 	)

// 	return &Server{
// 		GRPCServer: grpcServer,
// 		Listener:   lis,
// 	}, nil
// }
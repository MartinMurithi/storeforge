package grpcx

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ForwardMetadataInterceptor() grpc.UnaryClientInterceptor {

	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {

		// extract incoming metadata
		if md, ok := metadata.FromIncomingContext(ctx); ok {

			// attach to outgoing request
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
package grpcx

import (
	"context"
	"errors"

	"google.golang.org/grpc/metadata"
)

// GetUserIDFromMetadata extracts the user-id from gRPC incoming context.
func GetUserIDFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata is not provided")
	}

	// gRPC metadata keys are always lowercase
	values := md.Get("user-id")
	if len(values) == 0 {
		return "", errors.New("user-id not found in metadata")
	}

	return values[0], nil
}

// GetTenantIDFromMetadata extracts the tenant-id from gRPC incoming context.
func GetTenantIDFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata is not provided")
	}

	// gRPC metadata keys are always lowercase
	values := md.Get("tenant-id")
	if len(values) == 0 {
		return "", errors.New("tenant-id not found in metadata")
	}

	return values[0], nil
}

// ForwardMetadata takes incoming metadata from the current context 
// and prepares it to be sent to the next microservice.
func ForwardMetadata(ctx context.Context) context.Context {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return metadata.NewOutgoingContext(ctx, md)
	}
	return ctx
}
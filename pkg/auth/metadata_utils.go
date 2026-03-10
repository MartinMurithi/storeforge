package auth

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


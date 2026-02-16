package errconv

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FromGrpcToHttp converts a gRPC error into an HTTP status and a standard slug.
func FromGrpcToHttp(err error) (int, string, string) {
    status, ok := status.FromError(err)
    if !ok {
        return http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred"
    }

    msg := status.Message()

    switch status.Code() {
    case codes.InvalidArgument:
        return http.StatusBadRequest, "INVALID_ARGUMENT", msg
    case codes.Unauthenticated:
        return http.StatusUnauthorized, "UNAUTHORIZED", msg
    case codes.PermissionDenied:
        return http.StatusForbidden, "FORBIDDEN", msg
    case codes.NotFound:
        return http.StatusNotFound, "NOT_FOUND", msg
    case codes.AlreadyExists:
        return http.StatusConflict, "ALREADY_EXISTS", msg
    case codes.DeadlineExceeded:
        return http.StatusGatewayTimeout, "TIMEOUT", "Service took too long to respond"
    default:
        return http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error"
    }
}
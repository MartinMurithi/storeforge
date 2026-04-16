# ==== BUILDER STAGE ==========
#  Use the official Go alpine image bcs it is lightweight
FROM golang:1.25-alpine AS builder

# Install git and build-base( required for fetching certain Go modules and CGO tasks)
#  --no-cache does not store the index locally.
RUN apk add --no-cache git build-base

# Set the working dir inside the container
WORKDIR /app

# Copy dependency files first
# This levarages Docker's layer caching. If go.mod hasn't changed, docker skips go mod download in future builds; hence saving time.
COPY go.mod go.sum ./
RUN go mod download

# COPY the REQUIRED shared directories first
# These provide the packages like pkg/auth, pkg/env, etc.
COPY pkg/ ./pkg/
COPY api/ ./api/

# Copy the rest of the source code
COPY gateway/ ./gateway/

# Compile the binary
# CGO_ENABLED=0 creates a statically linked binary(runs without external C libraries)
# GOOS=linux ensures the binary is built for a linux environment
# -ldflags="-s -w" strips debug information to reduce the binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o gateway-app ./gateway/cmd/server/main.go

# =========== RUNTIME ================================

# Start with a fresh, tiny alpine image (5mb) for prod
FROM alpine:3.19

# Install CA certificates so the service can make secure HTTPS/TLS requests
RUN apk add --no-cache ca-certificates

# Create a non-root user for security (best practice to avoid running as root)
RUN adduser -D -g '' storeforgeuser

WORKDIR /home/storeforgeuser/

# Copy ONLY the compiled binary from the 'builder' stage.
# This keeps the final image clean of source code and build tools.
COPY --from=builder /app/gateway-app ./gateway-app

# Set ownership of the binary to the non-root user
RUN chown storeforgeuser:storeforgeuser ./gateway-app

RUN chmod +x ./gateway-app

# Switch to the non-root user
USER storeforgeuser

# Inform Docker that the service listens on the http port
EXPOSE 9095

# Run the service
CMD ["./gateway-app"] 
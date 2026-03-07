package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/bootstrap"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/config"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize application
	app, err := bootstrap.Init(cfg)

	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	// Graceful shutdown context
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start gRPC server
	go func() {
		log.Printf("[gRPC] listening on %s", app.GRPCServer.Listener.Addr())
		if err := app.GRPCServer.Start(); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("shutdown signal received")

	// Graceful shutdown gRPC
	app.GRPCServer.Stop()

	// Close DB pool
	if app.DB != nil {
		app.DB.Close()
	}

	log.Println("[TENANT MANAGEMENT SERVICE]: Shutdown complete")
}

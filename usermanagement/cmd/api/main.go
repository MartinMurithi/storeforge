package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/bootstrap"
)

func main() {

	app, err := bootstrap.Init()

	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	srv := &http.Server{
		Addr:         ":8585",
		Handler:      app.Router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	defer stop()

	// Start HTTP server
	go func() {
		log.Println("[HTTP] listening on :8585")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Start gRPC server
	go func() {
		log.Printf("[gRPC] listening on %s", app.GRPCServer.Listener.Addr())
		if err := app.GRPCServer.Start(); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Graceful shutdown gRPC
	app.GRPCServer.Stop()

	srv.Shutdown(shutdownCtx)

	if app.DB != nil {
		app.DB.Close()
	}

	log.Println("[AUTH SERVICE]: Shutdown complete")

}

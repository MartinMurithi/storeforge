package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MartinMurithi/storeforge/auth/internal/bootstrap"
	"github.com/MartinMurithi/storeforge/auth/internal/routes"
)

func main() {

	app, err := bootstrap.Init()

	if err != nil{
		log.Fatalf("failed to initialize app: %v", err)
	}
	// Register routes with auth middleware
	routes.NewUserRouter(app.Router)

	// Setup HTTP server
	srv := &http.Server{
		Addr:    ":8585",
		Handler: app.Router,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen error: %v", err)
		}
	}()
	log.Println("[AUTH SERVICE]: Running on port 8585")

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("[AUTH SERVICE]: Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	if app.DB != nil {
		app.DB.Close()
	}

	log.Println("[AUTH SERVICE]: Shutdown complete")

}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/MartinMurithi/storeforge/auth/internal/bootstrap"
	"github.com/MartinMurithi/storeforge/auth/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	app, err := bootstrap.Init()

	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	if err := router.SetTrustedProxies([]string{}); err != nil {
		log.Fatalf("failed to set trusted proxies%v", err)
	}

	router.GET("/", func(c *gin.Context) {
		fmt.Printf("ClientIP: %s\n", c.ClientIP())
	})

	routes.RegisterUserRoutes(router, app.Handler)

	srv := &http.Server{
		Addr:         ":8585",
		Handler:      router,
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

	go func() {
		log.Println("[SERVER] listening on :8585")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	srv.Shutdown(shutdownCtx)

	if app.DB != nil {
		app.DB.Close()
	}

	log.Println("[AUTH SERVICE]: Shutdown complete")

}

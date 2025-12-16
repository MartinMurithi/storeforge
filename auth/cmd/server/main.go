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

	if err := router.SetTrustedProxies([]string{"192.168.1.2"}); err != nil {
		log.Fatalf("failed to set trusted proxies%v", err)
	}

	routes.NewUserRouter(router, app.Handler)

	router.GET("/", func(c *gin.Context) {
		// If the client is 192.168.1.2, use the X-Forwarded-For
		// header to deduce the original client IP from the trust-
		// worthy parts of that header.
		// Otherwise, simply return the direct client IP
		fmt.Printf("ClientIP: %s\n", c.ClientIP())
	})

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

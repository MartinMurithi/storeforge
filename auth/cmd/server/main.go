package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"


	"github.com/MartinMurithi/storeforge.io/internal/bootstrap"
)

func main() {
	log.Println("[AUTH SERVICE]: Hello auth.....")

	app, err := bootstrap.Init()

	if err != nil {
		log.Fatalf("failed to start application %s", err)
	}

	log.Println("[AUTH SERVICE]: Application started successfully")

	//To help with graceful shutdown
	stop := make(chan os.Signal, 1)

	//SIGINT, listens for CTRL + C
	//SIGTERM, listens for docker/kubernetes stop
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	log.Println("[AUTH SERVICE]: Running on port 8585. Press CTRL + C to stop")

	<-stop

	log.Println("[AUTH SERVICE]: Shutting down gracefully...")

	//give active requests time to finish
	// ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	// defer cancel()

	//Shutdown server
	// if err := httpServer.Shutdown(ctx); err != nil {
	//     log.Printf("HTTP server forced to shutdown: %v", err)
	// }

	if app.DB != nil {
		app.DB.Close()
	}

}

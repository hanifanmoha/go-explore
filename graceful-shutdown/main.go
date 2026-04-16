package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// Create a handler for / path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		for i := range 10 {
			slog.Info("Processing request", "step", i)
			time.Sleep(1 * time.Second)
		}

		w.Write([]byte("Hello World!"))
	})

	// Create server object
	srv := &http.Server{Addr: ":8080"}

	// Start the server in a goroutine
	go func() {
		slog.Info("Starting the server in port 8080...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
		}
	}()

	// Handle graceful shutdown

	// Create a channel to listen for OS signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	slog.Info("Waiting for signale ...")

	// Block until a signal is received
	<-sig

	slog.Info("Signal Received! Wait untill all processing is done (max 5 seconds)...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	slog.Info("Shutdown complete.")
}

// func main() {

// 	// Create a channel to listen for OS signals
// 	sig := make(chan os.Signal, 1)
// 	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
// 	slog.Info("Waiting for signale ...")

// 	// Block until a signal is received
// 	<-sig

// 	// Below code will be executed after a signal is received
// 	slog.Info("Signal Received! Wait 5 seconds...")
// 	time.Sleep(5 * time.Second)
// 	slog.Info("Shutdown complete.")
// }

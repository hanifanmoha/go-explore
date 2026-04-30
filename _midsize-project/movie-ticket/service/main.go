package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/clients"
	"github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/config"
	"github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/handlers"
	"github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/middlewares"
	booking_repository "github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/repositories/booking"
	movie_repository "github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/repositories/movie"
	seat_repository "github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/repositories/seat"
	booking_service "github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/services/booking"
	movie_service "github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/services/movie"
	seat_service "github.com/hanifanmoha/go-explore/_midsize-project/movie-ticket/service/internal/services/seat"
)

func main() {

	config.LoadEnv()

	// Initialize clients
	dbClient, err := clients.NewDatabaseClient()
	if err != nil {
		slog.Error("Failed to connect to the database", "error", err)
		return
	}
	defer dbClient.Close()

	movieRepository := movie_repository.NewMovieRepository(dbClient)
	seatRepository := seat_repository.NewSeatRepository(dbClient)
	bookingRepository := booking_repository.NewBookingRepository(dbClient)

	movieService := movie_service.NewMovieService(movieRepository)
	seatService := seat_service.NewSeatService(seatRepository)
	bookingService := booking_service.NewBookingService(bookingRepository, movieRepository, seatRepository, dbClient)

	handler := handlers.NewHandler(bookingService, movieService, seatService)

	// Set up routes
	router := http.NewServeMux()
	router.HandleFunc("GET /movies", handler.GetMovies)
	router.HandleFunc("GET /movies/{id}", handler.GetMovie)
	router.HandleFunc("GET /movies/{id}/seats", handler.GetSeats)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.GetPort()),
		Handler: middlewares.Middlewares(router),
	}

	// Start the server
	go func() {
		slog.Info(fmt.Sprintf("Starting the server in port %s...", config.GetPort()))
		if err := srv.ListenAndServe(); err != nil {
			slog.Error("Failed to start server", "error", err)
		}
	}()

	// Handle shutdown gracefully
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	slog.Info("Shutting down the server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)

	slog.Info("Server stopped!")
}

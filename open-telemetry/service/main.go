package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type User struct {
	ID           int
	Username     string
	Email        string
	Name         string
	Organization Organization
}

type Organization struct {
	ID          int
	Name        string
	CountryCode string
}

func getOrganization() Organization {
	return Organization{
		ID:          1,
		Name:        "Acme Corporation",
		CountryCode: "US",
	}
}

func getUser() (User, error) {

	// random 0 - 5 seconds
	sleepMillis := rand.IntN(5000)
	time.Sleep(time.Duration(sleepMillis) * time.Millisecond)

	// simulate error (10%)
	errChance := rand.IntN(100)
	if errChance < 10 {
		return User{}, fmt.Errorf("failed to get user")
	}

	return User{
		ID:           1,
		Username:     "johndoe",
		Email:        "johndoe@example.com",
		Name:         "John Doe",
		Organization: getOrganization(),
	}, nil
}

var (
	requests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total requests",
		},
		[]string{"path", "method", "status"},
	)

	errors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total errors",
		},
		[]string{"path", "method", "status"},
	)

	duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Request duration",
		},
		[]string{"path", "method", "status"},
	)
)

func main() {

	prometheus.MustRegister(requests)
	prometheus.MustRegister(errors)
	prometheus.MustRegister(duration)

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requests.WithLabelValues(r.URL.Path, r.Method, "200").Inc()

		user, err := getUser()
		if err != nil {
			errors.WithLabelValues(r.URL.Path, r.Method, "500").Inc()
			http.Error(w, err.Error(), http.StatusInternalServerError)

			duration.WithLabelValues(r.URL.Path, r.Method, "500").Observe(time.Since(start).Seconds())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)

		duration.WithLabelValues(r.URL.Path, r.Method, "200").Observe(time.Since(start).Seconds())
	})

	srv := &http.Server{Addr: ":8080"}

	go func() {
		slog.Info("Starting the server in port 8080...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	slog.Info("Waiting for signale ...")

	<-sig

	slog.Info("Signal Received! Wait untill all processing is done...")
	srv.Shutdown(context.Background())
	slog.Info("Shutdown complete.")
}

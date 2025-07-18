package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var app *newrelic.Application

// generateRandomUserID creates a random user ID for tracing
func generateRandomUserID() string {
	return fmt.Sprintf("user_%d", rand.Intn(10000))
}

func GetFromEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// APMTrace holds both Datadog and New Relic tracing objects
type APMTrace struct {
	DDSpan    tracer.Span
	NRSegment *newrelic.Segment
}

// StartTrace creates a new trace for both Datadog and New Relic with custom attributes
func StartTrace(ctx context.Context, operationName string, attributes map[string]interface{}) (*APMTrace, context.Context) {
	// Start Datadog span with proper parent context
	ddSpan, ddCtx := tracer.StartSpanFromContext(ctx, operationName)

	// Add custom attributes to Datadog span
	for key, value := range attributes {
		ddSpan.SetTag(key, value)
	}

	// Start New Relic segment
	var nrSegment *newrelic.Segment
	if txn := newrelic.FromContext(ctx); txn != nil {
		nrSegment = txn.StartSegment(operationName)
		// Add custom attributes to New Relic segment
		for key, value := range attributes {
			nrSegment.AddAttribute(key, value)
		}
	}

	trace := &APMTrace{
		DDSpan:    ddSpan,
		NRSegment: nrSegment,
	}

	return trace, ddCtx
}

// Finish ends both Datadog span and New Relic segment
func (t *APMTrace) Finish() {
	if t.DDSpan != nil {
		t.DDSpan.Finish()
	}
	if t.NRSegment != nil {
		t.NRSegment.End()
	}
}

func pingRepo1(ctx context.Context) {
	// Start dual APM tracing
	attributes := map[string]interface{}{
		"user_id": generateRandomUserID(),
		"repo":    "repo1",
	}
	trace, _ := StartTrace(ctx, "ping.repo1", attributes)
	defer trace.Finish()

	// sleep randomly to simulate work
	randomSleep := time.Duration(50+rand.Intn(100)) * time.Millisecond
	time.Sleep(randomSleep)
}

func pingRepo2(ctx context.Context) error {
	// Start dual APM tracing
	attributes := map[string]interface{}{
		"user_id": generateRandomUserID(),
		"repo":    "repo2",
	}
	trace, _ := StartTrace(ctx, "ping.repo2", attributes)
	defer trace.Finish()

	i := rand.Intn(10)
	if i < 3 {
		errString := fmt.Sprintf("simulated error type %d", i)
		return errors.New(errString)
	} else if i < 5 {
		errString := fmt.Sprintf("simulated error type %d", i)
		return errors.New(errString)
	}

	// sleep randomly to simulate work
	randomSleep := time.Duration(rand.Intn(100)) * time.Millisecond
	time.Sleep(randomSleep)

	return nil
}

func pingRepo3(ctx context.Context) {
	// Start dual APM tracing
	attributes := map[string]interface{}{
		"user_id": generateRandomUserID(),
		"repo":    "repo3",
	}
	trace, _ := StartTrace(ctx, "ping.repo3", attributes)
	defer trace.Finish()

	// sleep randomly to simulate work
	randomSleep := time.Duration(rand.Intn(100)) * time.Millisecond
	time.Sleep(randomSleep)
}

func pingService(ctx context.Context) {
	// Start dual APM tracing
	attributes := map[string]interface{}{
		"user_id": generateRandomUserID(),
		"service": "ping_service",
	}
	trace, ctx := StartTrace(ctx, "ping.service", attributes)
	defer trace.Finish()

	time.Sleep(100 * time.Millisecond)
	pingRepo1(ctx)
	pingRepo2(ctx)
	pingRepo3(ctx)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	// Start New Relic transaction for HTTP request
	txn := app.StartTransaction("ping.handler")
	defer txn.End()

	// Add request/response data to New Relic
	txn.SetWebRequestHTTP(r)
	w = txn.SetWebResponse(w)

	// Create context with New Relic transaction
	ctx := newrelic.NewContext(r.Context(), txn)

	// Start dual APM tracing
	attributes := map[string]interface{}{
		"user_id": generateRandomUserID(),
		"handler": "ping_handler",
		"method":  r.Method,
		"path":    r.URL.Path,
	}
	trace, ctx := StartTrace(ctx, "ping.handler", attributes)
	defer trace.Finish()

	time.Sleep(200 * time.Millisecond)

	pingService(ctx)

	fmt.Fprintln(w, "pong")
}

func main() {
	// Start Datadog tracer
	tracer.Start(
		tracer.WithServiceName("timed-exam-api"),
		tracer.WithEnv("local"),
	)
	defer tracer.Stop()

	// Initialize New Relic application
	var err error
	newRelicLicense := GetFromEnv("NEW_RELIC_LICENSE_KEY", "")
	app, err = newrelic.NewApplication(
		newrelic.ConfigAppName("timed-exam-api"),
		newrelic.ConfigLicense(newRelicLicense),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		fmt.Println("Error initializing New Relic:", err)
	}

	http.HandleFunc("/ping", pingHandler)

	fmt.Println("Listening on :8080")

	http.ListenAndServe(":8080", nil)
}

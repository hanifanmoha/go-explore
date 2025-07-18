package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var app *newrelic.Application

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

// StartTrace creates a new trace for both Datadog and New Relic
func StartTrace(ctx context.Context, operationName string) (*APMTrace, context.Context) {
	// Start Datadog span
	ddSpan := tracer.StartSpan(operationName)

	// Start New Relic segment
	var nrSegment *newrelic.Segment
	if txn := newrelic.FromContext(ctx); txn != nil {
		nrSegment = txn.StartSegment(operationName)
	}

	trace := &APMTrace{
		DDSpan:    ddSpan,
		NRSegment: nrSegment,
	}

	return trace, ctx
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
	trace, _ := StartTrace(ctx, "ping.repo1")
	defer trace.Finish()

	// sleep randomly to simulate work
	randomSleep := time.Duration(rand.Intn(100)) * time.Millisecond
	time.Sleep(randomSleep)
}

func pingRepo2(ctx context.Context) {
	// Start dual APM tracing
	trace, _ := StartTrace(ctx, "ping.repo2")
	defer trace.Finish()

	// sleep randomly to simulate work
	randomSleep := time.Duration(rand.Intn(100)) * time.Millisecond
	time.Sleep(randomSleep)
}

func pingService(ctx context.Context) {
	// Start dual APM tracing
	trace, _ := StartTrace(ctx, "ping.service")
	defer trace.Finish()

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		pingRepo1(ctx)
	}()
	go func() {
		defer wg.Done()
		pingRepo2(ctx)
	}()
	wg.Wait()
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
	trace, _ := StartTrace(ctx, "ping.handler")
	defer trace.Finish()

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

// Example 79
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Simple trace package to demonstrate concepts
type trace struct{}

// SpanContext holds trace and span identifiers
type SpanContext struct {
	TraceID string
	SpanID  string
}

func NewSpanContext() SpanContext {
	traceID := generateID()
	spanID := generateID()
	return SpanContext{
		TraceID: traceID,
		SpanID:  spanID,
	}
}

func (sc SpanContext) String() string {
	return sc.TraceID
}

// Span represents a unit of work
type Span struct {
	ctx      SpanContext
	name     string
	tags     map[string]string
	start    time.Time
	finished bool
}

func (s *Span) Context() SpanContext {
	return s.ctx
}

func (s *Span) Finish() {
	if !s.finished {
		s.finished = true
		duration := time.Since(s.start)
		log.Printf("Span %s finished in %v", s.name, duration)
	}
}

// SpanOption allows configuring spans when creating them
type SpanOption interface {
	apply(*Span)
}

// ChildOfOption implements SpanOption
type ChildOfOption struct {
	Parent SpanContext
}

// Apply implements SpanOption
func (c ChildOfOption) apply(s *Span) {
	s.ctx.TraceID = c.Parent.TraceID
	// Generate new span ID but keep the trace ID
	s.ctx.SpanID = generateID()
}

// ChildOf creates a child span from parent context
func ChildOf(parent SpanContext) SpanOption {
	return ChildOfOption{Parent: parent}
}

// Tags adds tags to a span
type Tags map[string]string

// Make Tags satisfy the SpanOption interface
func (t Tags) apply(s *Span) {
	for k, v := range t {
		s.tags[k] = v
	}
}

// Tracer creates and manages spans
type Tracer struct{}

func (t *Tracer) Extract(header http.Header) (SpanContext, error) {
	traceID := header.Get("X-Trace-ID")
	spanID := header.Get("X-Span-ID")

	if traceID == "" || spanID == "" {
		return SpanContext{}, fmt.Errorf("no trace context in headers")
	}

	return SpanContext{
		TraceID: traceID,
		SpanID:  spanID,
	}, nil
}

func (t *Tracer) StartSpan(name string, options ...SpanOption) *Span {
	span := &Span{
		ctx:   NewSpanContext(),
		name:  name,
		tags:  make(map[string]string),
		start: time.Now(),
	}

	for _, option := range options {
		option.apply(span)
	}

	log.Printf("Started span %s (trace: %s, span: %s)",
		span.name, span.ctx.TraceID, span.ctx.SpanID)

	return span
}

// Inject puts span context into HTTP headers
func (t *Tracer) Inject(ctx SpanContext, header http.Header) {
	header.Set("X-Trace-ID", ctx.TraceID)
	header.Set("X-Span-ID", ctx.SpanID)
}

// Helper function to generate random IDs
func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// ContextWithSpan adds span to context
func ContextWithSpan(ctx context.Context, span *Span) context.Context {
	return context.WithValue(ctx, "span", span)
}

// SpanFromContext extracts span from context
func SpanFromContext(ctx context.Context) (*Span, bool) {
	span, ok := ctx.Value("span").(*Span)
	return span, ok
}

// Logger interface for different logging implementations
type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

// Simple logger implementation
type SimpleLogger struct{}

func (l *SimpleLogger) Info(msg string, keysAndValues ...interface{}) {
	log.Printf("INFO: %s %v", msg, keysAndValues)
}

func (l *SimpleLogger) Error(msg string, keysAndValues ...interface{}) {
	log.Printf("ERROR: %s %v", msg, keysAndValues)
}

// MetricsRecorder interface for recording metrics
type MetricsRecorder interface {
	RecordRequest(method, path string, statusCode int, duration float64)
}

// Simple metrics implementation
type SimpleMetrics struct{}

func (m *SimpleMetrics) RecordRequest(method, path string, statusCode int, duration float64) {
	log.Printf("METRIC: %s %s %d %.2fms", method, path, statusCode, duration)
}

// TracingMiddleware handles distributed tracing
type TracingMiddleware struct {
	tracer  *Tracer
	metrics MetricsRecorder
	logger  Logger
}

func (m *TracingMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Extract trace context from headers
		spanCtx, err := m.tracer.Extract(r.Header)
		if err != nil {
			spanCtx = NewSpanContext()
		}

		// Create span
		span := m.tracer.StartSpan("http_request",
			ChildOf(spanCtx),
			Tags{
				"http.method": r.Method,
				"http.url":    r.URL.String(),
			},
		)
		defer span.Finish()

		// Add trace ID to response headers
		w.Header().Set("X-Trace-ID", span.Context().TraceID)
		w.Header().Set("X-Span-ID", span.Context().SpanID)

		// Continue with traced context
		next.ServeHTTP(w, r.WithContext(
			ContextWithSpan(ctx, span),
		))
	})
}

// Service A handler - calls Service B
func serviceAHandler(w http.ResponseWriter, r *http.Request) {
	span, ok := SpanFromContext(r.Context())
	if !ok {
		http.Error(w, "No span in context", http.StatusInternalServerError)
		return
	}

	// Log which service is handling the request
	log.Printf("Service A handling request (trace: %s, span: %s)",
		span.Context().TraceID, span.Context().SpanID)

	// Create a new client
	client := &http.Client{}

	// Create a request to Service B
	req, err := http.NewRequest("GET", "http://localhost:8081/api", nil)
	if err != nil {
		http.Error(w, "Failed to create request to Service B", http.StatusInternalServerError)
		return
	}

	// Inject trace context into the outgoing request
	tracer := &Tracer{}
	tracer.Inject(span.Context(), req.Header)

	// Call Service B
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to call Service B", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read response from Service B
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read Service B response", http.StatusInternalServerError)
		return
	}

	// Return combined response
	fmt.Fprintf(w, "Service A received: %s", body)
}

// Service B handler
func serviceBHandler(w http.ResponseWriter, r *http.Request) {
	span, ok := SpanFromContext(r.Context())
	if !ok {
		http.Error(w, "No span in context", http.StatusInternalServerError)
		return
	}

	// Log which service is handling the request
	log.Printf("Service B handling request (trace: %s, span: %s)",
		span.Context().TraceID, span.Context().SpanID)

	// Return response
	fmt.Fprint(w, "Hello from Service B!")
}

func main() {
	// Create dependencies
	tracer := &Tracer{}
	metrics := &SimpleMetrics{}
	logger := &SimpleLogger{}

	// Create middleware
	middleware := &TracingMiddleware{
		tracer:  tracer,
		metrics: metrics,
		logger:  logger,
	}

	// Start Service B in a separate goroutine
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/api", middleware.Wrap(http.HandlerFunc(serviceBHandler)))

		log.Println("Starting Service B on :8081")
		log.Fatal(http.ListenAndServe(":8081", mux))
	}()

	// Wait a moment for Service B to start
	time.Sleep(100 * time.Millisecond)

	// Start Service A
	mux := http.NewServeMux()
	mux.Handle("/", middleware.Wrap(http.HandlerFunc(serviceAHandler)))

	log.Println("Starting Service A on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
// Example 74
package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Logger interface for logging
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// SimpleLogger is a basic implementation of Logger
type SimpleLogger struct{}

func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	fmt.Printf("[INFO] "+msg+"\n", args...)
}

func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	fmt.Printf("[ERROR] "+msg+"\n", args...)
}

// MetricsRecorder interface for recording metrics
type MetricsRecorder interface {
	IncCounter(name string, labels ...string)
	ObserveLatency(name string, duration time.Duration, labels ...string)
}

// SimpleMetricsRecorder is a basic implementation of MetricsRecorder
type SimpleMetricsRecorder struct{}

func (m *SimpleMetricsRecorder) IncCounter(name string, labels ...string) {
	fmt.Printf("[METRIC] Incrementing counter %s with labels %v\n", name, labels)
}

func (m *SimpleMetricsRecorder) ObserveLatency(name string, duration time.Duration, labels ...string) {
	fmt.Printf("[METRIC] Observed latency for %s: %v with labels %v\n", name, duration, labels)
}

// CircuitBreaker for preventing cascading failures
type CircuitBreaker struct {
	name             string
	failureCount     int
	failureThreshold int
	state            string
	lastFailure      time.Time
	resetTimeout     time.Duration
	mutex            sync.Mutex
	logger           Logger
}

func NewCircuitBreaker(name string) *CircuitBreaker {
	return &CircuitBreaker{
		name:             name,
		failureThreshold: 3,
		state:            "CLOSED", // CLOSED = healthy, OPEN = tripping
		resetTimeout:     5 * time.Second,
		logger:           &SimpleLogger{},
	}
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mutex.Lock()
	if cb.state == "OPEN" {
		if time.Since(cb.lastFailure) > cb.resetTimeout {
			fmt.Printf("[CIRCUIT] %s: Half-open, allowing test request\n", cb.name)
			cb.state = "HALF-OPEN"
		} else {
			cb.mutex.Unlock()
			fmt.Printf("[CIRCUIT] %s: Open (tripped), fast-failing request\n", cb.name)
			return fmt.Errorf("circuit breaker open")
		}
	} else {
		fmt.Printf("[CIRCUIT] %s: Closed, executing request\n", cb.name)
	}
	cb.mutex.Unlock()

	err := fn()

	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if err != nil {
		cb.failureCount++
		fmt.Printf("[CIRCUIT] %s: Request failed, failure count: %d/%d\n",
			cb.name, cb.failureCount, cb.failureThreshold)

		if cb.failureCount >= cb.failureThreshold || cb.state == "HALF-OPEN" {
			cb.state = "OPEN"
			cb.lastFailure = time.Now()
			fmt.Printf("[CIRCUIT] %s: Circuit tripped open\n", cb.name)
		}
		return err
	}

	if cb.state == "HALF-OPEN" {
		fmt.Printf("[CIRCUIT] %s: Success in half-open state, resetting circuit\n", cb.name)
		cb.state = "CLOSED"
	}

	cb.failureCount = 0
	return nil
}

// Retrier for automatic retry on failures
type Retrier struct {
	maxRetries int
	logger     Logger
}

func NewRetrier(maxRetries int, logger Logger) Retrier {
	return Retrier{
		maxRetries: maxRetries,
		logger:     logger,
	}
}

func (r Retrier) Do(ctx context.Context, fn func() (*http.Response, error)) (*http.Response, error) {
	var lastErr error
	for attempt := 0; attempt <= r.maxRetries; attempt++ {
		if attempt > 0 {
			r.logger.Info("Retrying request, attempt %d of %d", attempt, r.maxRetries)
		}

		resp, err := fn()
		if err == nil && resp.StatusCode < 500 {
			// Only consider successful if status code is less than 500
			return resp, nil
		}

		if err == nil {
			// We got a response, but status indicates server error
			r.logger.Error("Request returned error status: %s", resp.Status)
			lastErr = fmt.Errorf("server error: %s", resp.Status)
		} else {
			// Network or other error
			lastErr = err
			r.logger.Error("Request failed: %v", err)
		}

		// Simple backoff
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Duration(attempt*100) * time.Millisecond):
			// Continue with retry
		}
	}

	return nil, fmt.Errorf("max retries reached: %w", lastErr)
}

// Option for configuring ServiceClient
type Option func(*ServiceClient)

// WithRetrier configures the retrier
func WithRetrier(retrier Retrier) Option {
	return func(client *ServiceClient) {
		client.retrier = retrier
	}
}

// WithCircuitBreaker configures the circuit breaker
func WithCircuitBreaker(cb *CircuitBreaker) Option {
	return func(client *ServiceClient) {
		client.circuitBreaker = cb
	}
}

// WithMetrics configures metrics
func WithMetrics(metrics MetricsRecorder) Option {
	return func(client *ServiceClient) {
		client.metrics = metrics
	}
}

// WithLogger configures logger
func WithLogger(logger Logger) Option {
	return func(client *ServiceClient) {
		client.logger = logger
	}
}

// Simple trace implementation for the example
type trace struct{}

func FromContext(ctx context.Context) trace {
	return trace{}
}

func (t trace) SpanContext() struct{ TraceID traceID } {
	return struct{ TraceID traceID }{TraceID: traceID("trace-123456789")}
}

type traceID string

func (t traceID) String() string {
	return string(t)
}

// ServiceClient from the original example
type ServiceClient struct {
	baseURL        string
	httpClient     *http.Client
	retrier        Retrier
	circuitBreaker *CircuitBreaker
	metrics        MetricsRecorder
	logger         Logger
}

func NewServiceClient(baseURL string, opts ...Option) *ServiceClient {
	client := &ServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:    100,
				MaxConnsPerHost: 100,
				IdleConnTimeout: 90 * time.Second,
			},
		},
		metrics: &SimpleMetricsRecorder{},
		logger:  &SimpleLogger{},
	}

	for _, opt := range opts {
		opt(client)
	}

	// Set defaults if not configured through options
	if client.retrier.maxRetries == 0 {
		client.retrier = NewRetrier(3, client.logger)
	}

	if client.circuitBreaker == nil {
		client.circuitBreaker = NewCircuitBreaker("default")
	}

	return client
}

func (c *ServiceClient) DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	start := time.Now()
	defer func() {
		c.metrics.ObserveLatency("service_request", time.Since(start),
			"method", req.Method,
			"path", req.URL.Path,
		)
	}()

	// Add tracing headers
	traceID := FromContext(ctx).SpanContext().TraceID.String()
	req.Header.Set("X-Trace-ID", traceID)

	// Execute with circuit breaker and retry
	var resp *http.Response
	err := c.circuitBreaker.Execute(func() error {
		var err error
		resp, err = c.retrier.Do(ctx, func() (*http.Response, error) {
			// Clone the request to prevent reuse of a closed request
			reqClone := req.Clone(req.Context())
			return c.httpClient.Do(reqClone)
		})
		return err
	})

	if err != nil {
		c.metrics.IncCounter("service_request_errors",
			"method", req.Method,
			"path", req.URL.Path,
			"error", err.Error(),
		)
		return nil, err
	}
	return resp, nil
}

func (c *ServiceClient) Get(ctx context.Context, path string) (*http.Response, error) {
	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.DoRequest(ctx, req)
}

// Simple HTTP server for testing
func startTestServer() *http.Server {
	var failureCount int
	var requestCount int

	mux := http.NewServeMux()

	// Regular endpoint
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		// Log trace ID to demonstrate distributed tracing
		traceID := r.Header.Get("X-Trace-ID")
		fmt.Printf("[SERVER] Received request with trace ID: %s\n", traceID)
		fmt.Fprintf(w, "Hello from the test server!")
	})

	// Endpoint to test retry mechanism - fails first two times
	mux.HandleFunc("/retry-test", func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount <= 2 {
			fmt.Printf("[SERVER] Simulating failure for retry test (attempt %d)\n", requestCount)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		fmt.Printf("[SERVER] Retry test succeeded on attempt %d\n", requestCount)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Success after retries!")
	})

	// Endpoint to test circuit breaker
	mux.HandleFunc("/circuit-test", func(w http.ResponseWriter, r *http.Request) {
		failureCount++
		fmt.Printf("[SERVER] Circuit breaker test, failure count: %d\n", failureCount)
		w.WriteHeader(http.StatusInternalServerError)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)
	return server
}

func main() {
	// Start a test server
	server := startTestServer()
	defer server.Shutdown(context.Background())

	// Create logger and metrics recorder
	logger := &SimpleLogger{}
	metrics := &SimpleMetricsRecorder{}

	// Create service client with options
	fmt.Println("\n=== BASIC REQUEST WITH TRACING AND METRICS ===")
	cb := NewCircuitBreaker("test-circuit")
	retrier := NewRetrier(3, logger)
	client := NewServiceClient("http://localhost:8080",
		WithLogger(logger),
		WithMetrics(metrics),
		WithCircuitBreaker(cb),
		WithRetrier(retrier),
	)

	// Make a request demonstrating tracing and metrics
	ctx := context.Background()
	resp, err := client.Get(ctx, "/test")
	if err != nil {
		logger.Error("Basic request failed: %v", err)
	} else {
		defer resp.Body.Close()
		logger.Info("Basic request succeeded with status: %s", resp.Status)
	}

	// Test retry mechanism
	fmt.Println("\n=== TESTING RETRY MECHANISM ===")
	retrierClient := NewServiceClient("http://localhost:8080",
		WithLogger(logger),
		WithMetrics(metrics),
		WithRetrier(retrier),
		WithCircuitBreaker(NewCircuitBreaker("retry-circuit")),
	)

	resp, err = retrierClient.Get(ctx, "/retry-test")
	if err != nil {
		logger.Error("Retry test failed after multiple attempts: %v", err)
	} else {
		defer resp.Body.Close()
		logger.Info("Retry test succeeded with status: %s", resp.Status)
	}

	// Test circuit breaker
	fmt.Println("\n=== TESTING CIRCUIT BREAKER ===")
	cbCircuit := NewCircuitBreaker("test-breaker")
	cbClient := NewServiceClient("http://localhost:8080",
		WithLogger(logger),
		WithMetrics(metrics),
		WithCircuitBreaker(cbCircuit),
		WithRetrier(retrier),
	)

	// Make multiple requests to trigger circuit breaker
	fmt.Println("Making multiple requests to trigger circuit breaker:")
	for i := 0; i < 5; i++ {
		// Brief pause to make log output cleaner
		if i > 0 {
			time.Sleep(100 * time.Millisecond)
		}

		resp, err := cbClient.Get(ctx, "/circuit-test")
		if err != nil {
			logger.Error("Circuit test request %d failed: %v", i+1, err)
		} else {
			defer resp.Body.Close()
			logger.Info("Circuit test request %d status: %s", i+1, resp.Status)
		}
	}

	fmt.Println("\nThis example demonstrates the HTTP Client Pattern with:")
	fmt.Println("1. Circuit breaker to prevent cascading failures")
	fmt.Println("2. Retry mechanism for handling transient errors")
	fmt.Println("3. Distributed tracing via headers")
	fmt.Println("4. Metrics collection for observability")
	fmt.Println("5. Configurable options using the functional options pattern")
}
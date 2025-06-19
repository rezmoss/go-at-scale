// Example 57
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Global logger instance
var logger *zap.Logger

type Logger struct {
	logger *zap.Logger
}

func NewLogger(env string) (*Logger, error) {
	var zapLogger *zap.Logger
	var err error

	if env == "production" {
		zapLogger, err = zap.NewProduction()
	} else {
		zapLogger, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, fmt.Errorf("creating logger: %w", err)
	}

	return &Logger{logger: zapLogger}, nil
}

// Custom response writer to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Structured logging middleware
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		wrapper := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// Process request
		next.ServeHTTP(wrapper, r)

		// Log request details
		logger.Info("request completed",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", wrapper.status),
			zap.Duration("duration", time.Since(start)),
		)
	})
}

// Health check handler
type HealthChecker struct {
	checks map[string]HealthCheck
}

type HealthCheck func() error

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		checks: make(map[string]HealthCheck),
	}
}

func (hc *HealthChecker) AddCheck(name string, check HealthCheck) {
	hc.checks[name] = check
}

func (hc *HealthChecker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	result := make(map[string]string)

	for name, check := range hc.checks {
		if err := check(); err != nil {
			status = http.StatusServiceUnavailable
			result[name] = fmt.Sprintf("unhealthy: %v", err)
		} else {
			result[name] = "healthy"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(result)
}

// Example endpoint handler
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	// Initialize logger
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	// Create a new health checker
	healthChecker := NewHealthChecker()

	// Add some example health checks
	healthChecker.AddCheck("database", func() error {
		// In a real app, check database connection
		return nil // Return nil for healthy
	})

	healthChecker.AddCheck("api", func() error {
		// In a real app, check external API availability
		return nil // Return nil for healthy
	})

	// Create and set up HTTP server
	mux := http.NewServeMux()
	mux.Handle("/health", healthChecker)
	mux.HandleFunc("/", helloHandler)

	// Wrap with logging middleware
	wrappedMux := LoggingMiddleware(mux)

	// Start server
	logger.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}
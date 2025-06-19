// Example 56
package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func EnableProfiling(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

// Custom performance metrics
type PerformanceMetrics struct {
	requestDuration *prometheus.HistogramVec
	errorCount      *prometheus.CounterVec
}

func NewPerformanceMetrics(reg prometheus.Registerer) *PerformanceMetrics {
	m := &PerformanceMetrics{
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "request_duration_seconds",
				Help:    "Time spent processing request",
				Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
			},
			[]string{"endpoint"},
		),
		errorCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "error_total",
				Help: "Total number of errors",
			},
			[]string{"type"},
		),
	}

	reg.MustRegister(m.requestDuration, m.errorCount)
	return m
}

// Example handler function that uses performance metrics
func exampleHandler(metrics *PerformanceMetrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Simulate some work
		time.Sleep(100 * time.Millisecond)

		// Occasionally simulate an error
		if time.Now().UnixNano()%5 == 0 {
			metrics.errorCount.WithLabelValues("timeout").Inc()
			http.Error(w, "Simulated error", http.StatusInternalServerError)
			return
		}

		// Record request duration
		duration := time.Since(start).Seconds()
		metrics.requestDuration.WithLabelValues(r.URL.Path).Observe(duration)

		fmt.Fprintf(w, "Hello from example handler! Request took %.3f seconds", duration)
	}
}

func main() {
	// Create a new prometheus registry
	reg := prometheus.NewRegistry()

	// Initialize performance metrics
	metrics := NewPerformanceMetrics(reg)

	// Create HTTP server mux
	mux := http.NewServeMux()

	// Enable profiling
	EnableProfiling(mux)

	// Register Prometheus handler
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	// Add example handler
	mux.Handle("/example", exampleHandler(metrics))

	// Add a simple home page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Go Performance Profiling Example\n\n")
		fmt.Fprintf(w, "Available endpoints:\n")
		fmt.Fprintf(w, "- /example - Example handler with metrics\n")
		fmt.Fprintf(w, "- /metrics - Prometheus metrics\n")
		fmt.Fprintf(w, "- /debug/pprof/ - Profiling index\n")
	})

	// Start the server
	fmt.Println("Server starting on http://localhost:8080")
	fmt.Println("- Visit http://localhost:8080/ for available endpoints")
	fmt.Println("- Visit http://localhost:8080/debug/pprof/ for profiling")
	fmt.Println("- Visit http://localhost:8080/metrics for Prometheus metrics")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
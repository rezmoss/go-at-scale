// Example 78
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// APIHandler interface for different API version handlers
type APIHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// MetricsRecorder interface for recording metrics
type MetricsRecorder interface {
	IncCounter(metric string, labels ...string)
}

// SimpleMetricsRecorder is a simple implementation of MetricsRecorder
type SimpleMetricsRecorder struct{}

func (m *SimpleMetricsRecorder) IncCounter(metric string, labels ...string) {
	fmt.Printf("Metric: %s, Labels: %v\n", metric, labels)
}

// Logger interface for logging
type Logger interface {
	Info(msg string, fields ...string)
}

// SimpleLogger is a simple implementation of Logger
type SimpleLogger struct{}

func (l *SimpleLogger) Info(msg string, fields ...string) {
	fmt.Printf("Log: %s, Fields: %v\n", msg, fields)
}

// VersionedAPI structure for handling versioned APIs
type VersionedAPI struct {
	versions       map[string]APIHandler
	defaultVersion string
	metrics        MetricsRecorder
	logger         Logger
	deprecated     map[string]bool
}

func (api *VersionedAPI) getRequestedVersion(r *http.Request) string {
	// First check header
	if version := r.Header.Get("Accept-Version"); version != "" {
		return version
	}

	// Then check URL path /v1/resource, /v2/resource, etc.
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) > 1 && strings.HasPrefix(parts[1], "v") {
		return parts[1]
	}

	// Return default version if nothing found
	return api.defaultVersion
}

func (api *VersionedAPI) isVersionDeprecated(version string) bool {
	deprecated, exists := api.deprecated[version]
	return exists && deprecated
}

func (api *VersionedAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	version := api.getRequestedVersion(r)

	handler, exists := api.versions[version]
	if !exists {
		api.metrics.IncCounter("api_version_not_found",
			"requested_version", version,
		)
		http.Error(w, "Version not supported", http.StatusNotFound)
		return
	}

	// Add version headers
	w.Header().Set("X-API-Version", version)
	if api.isVersionDeprecated(version) {
		w.Header().Set("Warning", "299 - API version deprecated")
	}

	handler.ServeHTTP(w, r)
}

// V1Handler handles v1 API requests
type V1Handler struct{}

func (h *V1Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is API v1 response")
}

// V2Handler handles v2 API requests
type V2Handler struct{}

func (h *V2Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is API v2 response with enhanced features")
}

func main() {
	// Create API versions
	v1Handler := &V1Handler{}
	v2Handler := &V2Handler{}

	// Create API with versions
	api := &VersionedAPI{
		versions: map[string]APIHandler{
			"v1": v1Handler,
			"v2": v2Handler,
		},
		defaultVersion: "v1",
		metrics:        &SimpleMetricsRecorder{},
		logger:         &SimpleLogger{},
		deprecated:     map[string]bool{"v1": true},
	}

	// Start HTTP server
	fmt.Println("API server running on http://localhost:8080")
	fmt.Println("Try accessing:")
	fmt.Println("- http://localhost:8080/v1/users (v1 - deprecated)")
	fmt.Println("- http://localhost:8080/v2/users (v2)")
	fmt.Println("- http://localhost:8080/unknown (not found)")

	http.Handle("/", api)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
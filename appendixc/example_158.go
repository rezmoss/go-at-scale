// Example 158
// internal/infrastructure/gateway/versioning.go
type VersionManager struct {
    versions map[string][]VersionHandler
    logger   Logger
    metrics  MetricsRecorder
}

type VersionHandler struct {
    Version     string
    Handler     http.Handler
    Deprecated  bool
    SunsetDate  *time.Time
}

func (vm *VersionManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Extract version from request (header, URL, etc.)
    version := r.Header.Get("Accept-Version")
    if version == "" {
        version = defaultVersion
    }

    // Find appropriate handler
    handlers, exists := vm.versions[version]
    if !exists {
        vm.metrics.IncCounter("version_not_found")
        http.Error(w, "Version not supported", http.StatusNotFound)
        return
    }

    // Add version headers
    for _, handler := range handlers {
        if handler.Deprecated {
            w.Header().Set("Deprecation", "true")
            if handler.SunsetDate != nil {
                w.Header().Set("Sunset", handler.SunsetDate.Format(time.RFC3339))
            }
        }
    }

    // Handle request
    handlers[0].Handler.ServeHTTP(w, r)
}
// Example 78
type VersionedAPI struct {
    versions map[string]APIHandler
    defaultVersion string
    metrics MetricsRecorder
    logger  Logger
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
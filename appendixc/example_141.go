// Example 141
// internal/monitoring/health/checker.go
type HealthChecker struct {
    checks map[string]CheckFunc
    logger *StructuredLogger
}

type CheckFunc func(context.Context) error

type HealthStatus struct {
    Status    string            `json:"status"`
    Checks    map[string]string `json:"checks"`
    Timestamp time.Time         `json:"timestamp"`
}

func (hc *HealthChecker) AddCheck(name string, check CheckFunc) {
    hc.checks[name] = check
}

func (hc *HealthChecker) RunChecks(ctx context.Context) HealthStatus {
    status := HealthStatus{
        Status:    "healthy",
        Checks:    make(map[string]string),
        Timestamp: time.Now(),
    }
    
    for name, check := range hc.checks {
        if err := check(ctx); err != nil {
            status.Status = "unhealthy"
            status.Checks[name] = err.Error()
            
            hc.logger.Error("health check failed",
                zap.String("check", name),
                zap.Error(err),
            )
        } else {
            status.Checks[name] = "ok"
        }
    }
    
    return status
}

// HTTP handler for health checks
func (hc *HealthChecker) Handler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        status := hc.RunChecks(r.Context())
        
        w.Header().Set("Content-Type", "application/json")
        if status.Status != "healthy" {
            w.WriteHeader(http.StatusServiceUnavailable)
        }
        
        json.NewEncoder(w).Encode(status)
    })
}
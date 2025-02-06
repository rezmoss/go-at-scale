// Example 57
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

// Structured logging middleware
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Wrap response writer to capture status code
        wrapper := &responseWriter{ResponseWriter: w}
        
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
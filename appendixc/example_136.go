// Example 136
// internal/monitoring/metrics/prometheus.go
type MetricsCollector struct {
    requestDuration    *prometheus.HistogramVec
    requestCount      *prometheus.CounterVec
    activeConnections *prometheus.GaugeVec
    errorCount        *prometheus.CounterVec
    queueLength       *prometheus.GaugeVec
}

func NewMetricsCollector(reg prometheus.Registerer) (*MetricsCollector, error) {
    mc := &MetricsCollector{
        requestDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name:    "request_duration_seconds",
                Help:    "Time spent processing requests",
                Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
            },
            []string{"method", "path", "status"},
        ),
        
        requestCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "request_total",
                Help: "Total number of requests processed",
            },
            []string{"method", "path", "status"},
        ),
        
        activeConnections: prometheus.NewGaugeVec(
            prometheus.GaugeOpts{
                Name: "active_connections",
                Help: "Number of active connections",
            },
            []string{"type"},
        ),
        
        errorCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "error_total",
                Help: "Total number of errors",
            },
            []string{"type", "code"},
        ),
        
        queueLength: prometheus.NewGaugeVec(
            prometheus.GaugeOpts{
                Name: "queue_length",
                Help: "Current queue length",
            },
            []string{"queue"},
        ),
    }
    
    // Register metrics
    collectors := []prometheus.Collector{
        mc.requestDuration,
        mc.requestCount,
        mc.activeConnections,
        mc.errorCount,
        mc.queueLength,
    }
    
    for _, collector := range collectors {
        if err := reg.Register(collector); err != nil {
            return nil, fmt.Errorf("registering collector: %w", err)
        }
    }
    
    return mc, nil
}

// Middleware for HTTP metrics
func (mc *MetricsCollector) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Track active connections
        mc.activeConnections.WithLabelValues("http").Inc()
        defer mc.activeConnections.WithLabelValues("http").Dec()
        
        // Use custom response writer to capture status code
        ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
        
        next.ServeHTTP(ww, r)
        
        // Record metrics
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(ww.Status())
        
        mc.requestDuration.WithLabelValues(r.Method, r.URL.Path, status).
            Observe(duration)
        
        mc.requestCount.WithLabelValues(r.Method, r.URL.Path, status).
            Inc()
    })
}
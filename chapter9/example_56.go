// Example 56
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
    errorCount     *prometheus.CounterVec
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
// Example 79
type TracingMiddleware struct {
    tracer  *trace.Tracer
    metrics MetricsRecorder
    logger  Logger
}

func (m *TracingMiddleware) Wrap(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        // Extract trace context from headers
        spanCtx, err := m.tracer.Extract(r.Header)
        if err != nil {
            spanCtx = trace.NewSpanContext()
        }

        // Create span
        span := m.tracer.StartSpan("http_request",
            trace.ChildOf(spanCtx),
            trace.Tags{
                "http.method": r.Method,
                "http.url": r.URL.String(),
            },
        )
        defer span.Finish()

        // Add trace ID to response headers
        w.Header().Set("X-Trace-ID", span.Context().TraceID.String())

        // Continue with traced context
        next.ServeHTTP(w, r.WithContext(
            trace.ContextWithSpan(ctx, span),
        ))
    })
}
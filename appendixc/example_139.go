// Example 139
// internal/monitoring/tracing/tracer.go
type Tracer struct {
    tracer     *trace.Tracer
    propagator propagation.TextMapPropagator
}

func NewTracer(serviceName string, endpoint string) (*Tracer, error) {
    ctx := context.Background()
    
    // Create exporter
    exp, err := otlptrace.New(ctx,
        otlptracegrpc.NewClient(
            otlptracegrpc.WithEndpoint(endpoint),
            otlptracegrpc.WithInsecure(),
        ),
    )
    if err != nil {
        return nil, fmt.Errorf("creating exporter: %w", err)
    }
    
    // Create tracer provider
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exp),
        trace.WithResource(
            resource.NewWithAttributes(
                semconv.SchemaURL,
                semconv.ServiceNameKey.String(serviceName),
            ),
        ),
    )
    
    otel.SetTracerProvider(tp)
    
    tracer := tp.Tracer(serviceName)
    
    return &Tracer{
        tracer:     &tracer,
        propagator: propagation.TraceContext{},
    }, nil
}

// Middleware for HTTP tracing
func (t *Tracer) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := t.propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
        
        ctx, span := (*t.tracer).Start(ctx, fmt.Sprintf("%s %s", r.Method, r.URL.Path))
        defer span.End()
        
        // Add trace info to response headers
        t.propagator.Inject(ctx, propagation.HeaderCarrier(w.Header()))
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
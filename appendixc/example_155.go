// Example 155
// internal/infrastructure/gateway/transform.go
type TransformationPipeline struct {
    transforms []Transformer
    logger     Logger
    metrics    MetricsRecorder
}

type Transformer interface {
    Transform(ctx context.Context, req *http.Request) error
}

func (p *TransformationPipeline) Transform(ctx context.Context, req *http.Request) error {
    start := time.Now()
    defer func() {
        p.metrics.ObserveLatency("request_transformation", time.Since(start))
    }()

    for _, t := range p.transforms {
        if err := t.Transform(ctx, req); err != nil {
            p.metrics.IncCounter("transformation_errors")
            return fmt.Errorf("applying transformation: %w", err)
        }
    }

    return nil
}

// Header Transformation
type HeaderTransformer struct {
    mappings map[string]string
    removes  []string
}

func (t *HeaderTransformer) Transform(ctx context.Context, req *http.Request) error {
    // Apply header mappings
    for src, dst := range t.mappings {
        if val := req.Header.Get(src); val != "" {
            req.Header.Set(dst, val)
        }
    }

    // Remove specified headers
    for _, header := range t.removes {
        req.Header.Del(header)
    }

    return nil
}
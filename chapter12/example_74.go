// Example 74
type ServiceClient struct {
    baseURL    string
    httpClient *http.Client
    retrier    Retrier
    circuitBreaker *CircuitBreaker
    metrics    MetricsRecorder
    logger     Logger
}

func NewServiceClient(baseURL string, opts ...Option) *ServiceClient {
    client := &ServiceClient{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
            Transport: &http.Transport{
                MaxIdleConns:        100,
                MaxConnsPerHost:     100,
                IdleConnTimeout:     90 * time.Second,
                DisableCompression:  false,
                DisableKeepAlives:   false,
            },
        },
    }
    
    for _, opt := range opts {
        opt(client)
    }
    
    return client
}

func (c *ServiceClient) DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
    start := time.Now()
    defer func() {
        c.metrics.ObserveLatency("service_request", time.Since(start),
            "method", req.Method,
            "path", req.URL.Path,
        )
    }()
    // Add tracing headers
    traceID := trace.FromContext(ctx).SpanContext().TraceID.String()
    req.Header.Set("X-Trace-ID", traceID)
    // Execute with circuit breaker and retry
    var resp *http.Response
    err := c.circuitBreaker.Execute(func() error {
        var err error
        resp, err = c.retrier.Do(ctx, func() (*http.Response, error) {
            return c.httpClient.Do(req)
        })
        return err
    })

    if err != nil {
        c.metrics.IncCounter("service_request_errors",
            "method", req.Method,
            "path", req.URL.Path,
            "error", err.Error(),
        )
        return nil, err
    }
    return resp, nil
}
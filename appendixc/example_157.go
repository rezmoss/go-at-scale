// Example 157
// internal/infrastructure/gateway/aggregation.go
type Aggregator struct {
    endpoints []Endpoint
    client    *http.Client
    cache     Cache
    logger    Logger
    metrics   MetricsRecorder
}

type Endpoint struct {
    Name     string
    URL      string
    Required bool
    Timeout  time.Duration
}

func (a *Aggregator) Aggregate(ctx context.Context) (map[string]interface{}, error) {
    start := time.Now()
    defer func() {
        a.metrics.ObserveLatency("request_aggregation", time.Since(start))
    }()

    results := make(map[string]interface{})
    errors := make([]error, 0)
    
    var wg sync.WaitGroup
    resultCh := make(chan struct {
        name string
        data interface{}
        err  error
    }, len(a.endpoints))

    // Launch requests in parallel
    for _, endpoint := range a.endpoints {
        wg.Add(1)
        go func(ep Endpoint) {
            defer wg.Done()

            // Check cache first
            if data, found := a.cache.Get(ep.Name); found {
                resultCh <- struct {
                    name string
                    data interface{}
                    err  error
                }{ep.Name, data, nil}
                return
            }

            // Create context with timeout
            ctx, cancel := context.WithTimeout(ctx, ep.Timeout)
            defer cancel()

            // Make request
            resp, err := a.client.Get(ep.URL)
            if err != nil {
                resultCh <- struct {
                    name string
                    data interface{}
                    err  error
                }{ep.Name, nil, err}
                return
            }
            defer resp.Body.Close()

            // Parse response
            var data interface{}
            if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
                resultCh <- struct {
                    name string
                    data interface{}
                    err  error
                }{ep.Name, nil, err}
                return
            }

            // Cache result
            a.cache.Set(ep.Name, data, defaultCacheTTL)

            resultCh <- struct {
                name string
                data interface{}
                err  error
            }{ep.Name, data, nil}
        }(endpoint)
    }

    // Wait for all requests
    go func() {
        wg.Wait()
        close(resultCh)
    }()

    // Collect results
    for result := range resultCh {
        if result.err != nil {
            errors = append(errors, fmt.Errorf("%s: %w", result.name, result.err))
            continue
        }
        results[result.name] = result.data
    }

    // Check if any required endpoints failed
    for _, endpoint := range a.endpoints {
        if endpoint.Required {
            if _, ok := results[endpoint.Name]; !ok {
                a.metrics.IncCounter("aggregation_errors")
                return nil, fmt.Errorf("required endpoint %s failed", endpoint.Name)
            }
        }
    }

    return results, nil
}
// Example 77
type CircuitBreaker struct {
    name           string
    maxFailures    int
    resetTimeout   time.Duration
    failureCount   int64
    lastFailure    time.Time
    state          State
    metrics        MetricsRecorder
    mutex          sync.RWMutex
}

func (cb *CircuitBreaker) Execute(operation func() error) error {
    if !cb.allowRequest() {
        cb.metrics.IncCounter("circuit_breaker_rejections",
            "name", cb.name,
        )
        return ErrCircuitOpen
    }

    err := operation()
    
    if err != nil {
        cb.recordFailure()
        return err
    }

    cb.recordSuccess()
    return nil
}

type Fallback struct {
    primary   Operation
    fallback  Operation
    metrics   MetricsRecorder
    logger    Logger
}

func (f *Fallback) Execute(ctx context.Context) error {
    err := f.primary(ctx)
    if err == nil {
        return nil
    }

    f.metrics.IncCounter("fallback_triggered")
    f.logger.Warn("primary operation failed, using fallback",
        "error", err,
    )

    return f.fallback(ctx)
}
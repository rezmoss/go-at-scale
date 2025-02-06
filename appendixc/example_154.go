// Example 154
// internal/infrastructure/gateway/ratelimit.go
type RateLimiter struct {
    store      RedisClient
    windowSize time.Duration
    limit      int
    metrics    MetricsRecorder
}

type RedisClient interface {
    Incr(ctx context.Context, key string) (int64, error)
    Expire(ctx context.Context, key string, ttl time.Duration) error
}

func NewRateLimiter(store RedisClient, windowSize time.Duration, limit int, metrics MetricsRecorder) *RateLimiter {
    return &RateLimiter{
        store:      store,
        windowSize: windowSize,
        limit:      limit,
        metrics:    metrics,
    }
}

func (rl *RateLimiter) Allow(ctx context.Context, identifier string) (bool, error) {
    start := time.Now()
    defer func() {
        rl.metrics.ObserveLatency("rate_limiter_check", time.Since(start))
    }()

    key := fmt.Sprintf("ratelimit:%s:%d", identifier, time.Now().Unix()/int64(rl.windowSize.Seconds()))
    
    count, err := rl.store.Incr(ctx, key)
    if err != nil {
        rl.metrics.IncCounter("rate_limiter_errors")
        return false, fmt.Errorf("incrementing counter: %w", err)
    }

    // Set expiration on first request in window
    if count == 1 {
        if err := rl.store.Expire(ctx, key, rl.windowSize); err != nil {
            rl.metrics.IncCounter("rate_limiter_errors")
            return false, fmt.Errorf("setting expiration: %w", err)
        }
    }

    allowed := count <= int64(rl.limit)
    if !allowed {
        rl.metrics.IncCounter("rate_limit_exceeded")
    }

    return allowed, nil
}
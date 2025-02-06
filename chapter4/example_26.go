// Example 26
type RateLimiter struct {
    ticker *time.Ticker
    limit  chan struct{}
}

func NewRateLimiter(rate int, burst int) *RateLimiter {
    limiter := &RateLimiter{
        ticker: time.NewTicker(time.Second / time.Duration(rate)),
        limit:  make(chan struct{}, burst),
    }
    
    // Fill token bucket
    for i := 0; i < burst; i++ {
        limiter.limit <- struct{}{}
    }
    
    // Replenish tokens
    go func() {
        for range limiter.ticker.C {
            select {
            case limiter.limit <- struct{}{}:
            default:
                // Bucket is full
            }
        }
    }()
    
    return limiter
}

func (r *RateLimiter) Wait() {
    <-r.limit
}

// Example usage
func main() {
    limiter := NewRateLimiter(10, 5) // 10 ops/sec, burst of 5
    
    for i := 0; i < 20; i++ {
        go func(i int) {
            limiter.Wait()
            log.Printf("Operation %d executed", i)
        }(i)
    }
}
// Example 156
// internal/infrastructure/gateway/circuitbreaker.go
type CircuitBreaker struct {
    name           string
    failureThreshold   float64
    resetTimeout   time.Duration
    halfOpenMax    int
    metrics        MetricsRecorder
    
    state         State
    failures      int64
    successes     int64
    lastStateChange time.Time
    mu            sync.RWMutex
}

type State int

const (
    StateClosed State = iota
    StateOpen
    StateHalfOpen
)

func (cb *CircuitBreaker) Execute(ctx context.Context, req *http.Request) (*http.Response, error) {
    if !cb.allowRequest() {
        cb.metrics.IncCounter("circuit_breaker_rejected")
        return nil, ErrCircuitOpen
    }

    resp, err := http.DefaultClient.Do(req)
    
    if err != nil || resp.StatusCode >= 500 {
        cb.recordFailure()
        return resp, err
    }

    cb.recordSuccess()
    return resp, nil
}

func (cb *CircuitBreaker) allowRequest() bool {
    cb.mu.RLock()
    defer cb.mu.RUnlock()

    switch cb.state {
    case StateClosed:
        return true
    case StateOpen:
        if time.Since(cb.lastStateChange) > cb.resetTimeout {
            cb.setState(StateHalfOpen)
            return true
        }
        return false
    case StateHalfOpen:
        return atomic.LoadInt64(&cb.successes) < int64(cb.halfOpenMax)
    default:
        return false
    }
}
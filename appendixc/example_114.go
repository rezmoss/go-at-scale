// Example 114
// internal/infrastructure/retry/retry.go
type RetryConfig struct {
    MaxAttempts int
    InitialWait time.Duration
    MaxWait     time.Duration
}

func WithRetry(operation func() error, config RetryConfig) error {
    var err error
    wait := config.InitialWait
    
    for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
        err = operation()
        if err == nil {
            return nil
        }
        
        // Check if error is retryable
        if !isRetryableError(err) {
            return err
        }
        
        // Last attempt
        if attempt == config.MaxAttempts {
            break
        }
        
        // Wait with exponential backoff
        time.Sleep(wait)
        wait = min(wait*2, config.MaxWait)
    }
    
    return fmt.Errorf("operation failed after %d attempts: %w",
        config.MaxAttempts, err)
}

// Circuit breaker
type CircuitBreaker struct {
    name     string
    timeout  time.Duration
    maxFails int
    failures int
    lastFail time.Time
    mu       sync.RWMutex
}

func (cb *CircuitBreaker) Execute(operation func() error) error {
    if !cb.canTry() {
        return fmt.Errorf("circuit breaker %s is open", cb.name)
    }
    
    err := operation()
    cb.recordResult(err)
    return err
}

func (cb *CircuitBreaker) canTry() bool {
    cb.mu.RLock()
    defer cb.mu.RUnlock()
    
    if cb.failures < cb.maxFails {
        return true
    }
    
    if time.Since(cb.lastFail) > cb.timeout {
        return true
    }
    
    return false
}
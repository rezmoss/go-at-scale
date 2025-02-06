// Example 36
type ContextAwareService struct {
    operations chan func()
    errors     chan error
}

func NewContextAwareService() *ContextAwareService {
    return &ContextAwareService{
        operations: make(chan func()),
        errors:     make(chan error, 1),
    }
}

func (s *ContextAwareService) Start(ctx context.Context) {
    go func() {
        for {
            select {
            case op := <-s.operations:
                op()
            case <-ctx.Done():
                s.errors <- ctx.Err()
                close(s.operations)
                return
            }
        }
    }()
}

// Example: Timeout-aware operations
func (s *ContextAwareService) ExecuteWithTimeout(
    ctx context.Context,
    operation func() error,
    timeout time.Duration,
) error {
    timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()
    
    done := make(chan error, 1)
    
    go func() {
        done <- operation()
    }()
    
    select {
    case err := <-done:
        return err
    case <-timeoutCtx.Done():
        return fmt.Errorf("operation timed out: %w", timeoutCtx.Err())
    }
}
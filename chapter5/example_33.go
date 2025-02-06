// Example 33
type Result[T any] struct {
    Value T
    Err   error
}

type ErrorCollector struct {
    errors []error
    mu     sync.Mutex
}

func (ec *ErrorCollector) Add(err error) {
    ec.mu.Lock()
    defer ec.mu.Unlock()
    ec.errors = append(ec.errors, err)
}

func (ec *ErrorCollector) Error() error {
    ec.mu.Lock()
    defer ec.mu.Unlock()
    
    if len(ec.errors) == 0 {
        return nil
    }
    
    errStrings := make([]string, len(ec.errors))
    for i, err := range ec.errors {
        errStrings[i] = err.Error()
    }
    
    return fmt.Errorf("multiple errors occurred: %s", strings.Join(errStrings, "; "))
}

// Example: Concurrent operations with error handling
func ConcurrentProcess[T any](items []T, process func(T) error) error {
    collector := &ErrorCollector{}
    var wg sync.WaitGroup
    
    for _, item := range items {
        wg.Add(1)
        go func(item T) {
            defer wg.Done()
            if err := process(item); err != nil {
                collector.Add(err)
            }
        }(item)
    }
    
    wg.Wait()
    return collector.Error()
}
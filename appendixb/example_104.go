// Example 104
// Pitfall 1: Unbuffered channel with no receiver
func leakyGoroutine() {
    ch := make(chan int)
    go func() {
        val := expensive()
        ch <- val  // Blocks forever if no one receives
    }()
    // Channel is never read from
}

// Solution: Use context for cancellation
func nonLeakyGoroutine(ctx context.Context) error {
    ch := make(chan int, 1)  // Buffered channel
    
    go func() {
        val := expensive()
        select {
        case ch <- val:
        case <-ctx.Done():
            return
        }
    }()
    
    select {
    case result := <-ch:
        return processResult(result)
    case <-ctx.Done():
        return ctx.Err()
    }
}

// Pitfall 2: Goroutines in loops
func leakyLoop() {
    for _, item := range items {
        go func() {
            process(item)  // item is shared across all goroutines!
        }()
    }
}

// Solution: Pass loop variables explicitly
func nonLeakyLoop() {
    for _, item := range items {
        item := item  // Create new variable for each iteration
        go func() {
            process(item)
        }()
    }
}
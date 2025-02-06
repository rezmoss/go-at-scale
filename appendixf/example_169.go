// Example 169
// Anti-pattern: Leaking goroutines
func processData(data []string) chan string {
    results := make(chan string)
    go func() {
        for _, item := range data {
            results <- process(item) // Channel might never be read
        }
    }()
    return results
}

// Proper pattern: Ensure goroutine cleanup
func processData(ctx context.Context, data []string) (<-chan string, error) {
    results := make(chan string)
    go func() {
        defer close(results)
        for _, item := range data {
            select {
            case results <- process(item):
            case <-ctx.Done():
                return
            }
        }
    }()
    return results, nil
}
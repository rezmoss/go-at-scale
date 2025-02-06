// Example 15
// Generic error handler
type ErrorHandler[T any] func(error) T

func WithErrorHandler[T any](f func() (T, error), handler ErrorHandler[T]) func() T {
    return func() T {
        result, err := f()
        if err != nil {
            return handler(err)
        }
        return result
    }
}

// Example usage
func fetchData() ([]string, error) {
    // Simulated fetch
    return nil, errors.New("network error")
}

handler := func(err error) []string {
    log.Printf("Error: %v", err)
    return []string{"fallback data"}
}

safeFetch := WithErrorHandler(fetchData, handler)
data := safeFetch()  // Returns fallback data on error
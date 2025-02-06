// Example 10
func timeTracker[T any](name string, fn func() T) func() T {
    return func() T {
        start := time.Now()
        defer func() {
            duration := time.Since(start)
            log.Printf("%s took %v to execute", name, duration)
        }()
        return fn()
    }
}

// Usage example
func expensiveOperation() int {
    time.Sleep(time.Second)
    return 42
}

func main() {
    trackedFn := timeTracker("expensive-op", expensiveOperation)
    result := trackedFn()  // Logs timing information
}
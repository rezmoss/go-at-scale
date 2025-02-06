// Example 93
// Example with race condition
type Counter struct {
    value int
}

func (c *Counter) Increment() {
    c.value++  // Race condition!
}

// Fixed version
type SafeCounter struct {
    mu    sync.Mutex
    value int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

// Run tests with race detection
// go test -race ./...
// Build with race detection
// go build -race
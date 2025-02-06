// Example 170
// Anti-pattern: Using channels for everything
type Cache struct {
    data    chan map[string]string  // Wrong: Using channel as mutex
    updates chan string
}

// Proper pattern: Use appropriate synchronization
type Cache struct {
    mu      sync.RWMutex
    data    map[string]string
    updates chan string  // Channel only for actual communication
}
// Example 34
// Mutex-based counter
type MutexCounter struct {
    mu    sync.Mutex
    value int
}

func (c *MutexCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *MutexCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

// Channel-based counter
type ChannelCounter struct {
    updates chan int
    value   chan int
    done    chan struct{}
}

func NewChannelCounter() *ChannelCounter {
    c := &ChannelCounter{
        updates: make(chan int),
        value:   make(chan int),
        done:    make(chan struct{}),
    }
    
    go func() {
        var current int
        for {
            select {
            case <-c.updates:
                current++
            case c.value <- current:
                // Value requested
            case <-c.done:
                close(c.value)
                return
            }
        }
    }()
    
    return c
}

func (c *ChannelCounter) Increment() {
    c.updates <- 1
}

func (c *ChannelCounter) Value() int {
    return <-c.value
}

func (c *ChannelCounter) Close() {
    close(c.done)
}
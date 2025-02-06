// Example 29
type BufferConfig struct {
    Size    int
    Dynamic bool
}

func NewBuffer[T any](config BufferConfig) chan T {
    if config.Dynamic {
        return make(chan T, runtime.GOMAXPROCS(0))
    }
    return make(chan T, config.Size)
}

// Example: Adaptive buffering based on load
type AdaptiveProcessor[T any] struct {
    input    chan T
    output   chan T
    capacity int
}

func (ap *AdaptiveProcessor[T]) resize() {
    currentLoad := len(ap.input) / cap(ap.input)
    if currentLoad > 0.8 && ap.capacity < 1000 {
        // Increase buffer size
        newInput := make(chan T, ap.capacity*2)
        ap.capacity *= 2
        // Transfer existing items
        close(ap.input)
        for item := range ap.input {
            newInput <- item
        }
        ap.input = newInput
    }
}
// Example 27
type Pipeline[T any] struct {
    input     chan T
    output    chan T
    processor func(T) T
    buffer    int
}

func NewPipeline[T any](processor func(T) T, buffer int) *Pipeline[T] {
    return &Pipeline[T]{
        input:     make(chan T, buffer),
        output:    make(chan T, buffer),
        processor: processor,
        buffer:    buffer,
    }
}

func (p *Pipeline[T]) Start(ctx context.Context) {
    go func() {
        defer close(p.output)
        
        for {
            select {
            case <-ctx.Done():
                return
            case item, ok := <-p.input:
                if !ok {
                    return
                }
                
                // Process with backpressure
                select {
                case p.output <- p.processor(item):
                    // Successfully processed and sent
                case <-ctx.Done():
                    return
                }
            }
        }
    }()
}
// Example: Slow consumer handling
func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    pipe := NewPipeline(func(x int) int {
        time.Sleep(100 * time.Millisecond)
        return x * 2
    }, 5)
    
    pipe.Start(ctx)
    // Producer
    go func() {
        for i := 0; i < 100; i++ {
            select {
            case pipe.input <- i:
                // Successfully sent
            case <-ctx.Done():
                return
            }
        }
        close(pipe.input)
    }()
    // Consumer
    for result := range pipe.output {
        log.Printf("Processed: %d", result)
        time.Sleep(200 * time.Millisecond) // Slow consumer
    }
}
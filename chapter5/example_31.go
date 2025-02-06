// Example 31
type Multiplexer[T any] struct {
    inputs  []<-chan T
    output  chan T
    done    chan struct{}
    timeout time.Duration
}

func NewMultiplexer[T any](timeout time.Duration) *Multiplexer[T] {
    return &Multiplexer[T]{
        output:  make(chan T),
        done:    make(chan struct{}),
        timeout: timeout,
    }
}

func (m *Multiplexer[T]) AddInput(ch <-chan T) {
    m.inputs = append(m.inputs, ch)
}

func (m *Multiplexer[T]) Start() <-chan T {
    go func() {
        defer close(m.output)
        

        // WARNING: Using reflect.SelectCase is an advanced technique and can be
        // confusing or harder to maintain in production code. If you have a small,
        // fixed set of channels, prefer a normal select {...} statement or a
        // fan-in approach with one goroutine per channel. For dynamic scenarios,
        // consider alternative designs or be aware of the complexity reflect.Select
        // introduces.

        // Create cases for select
        cases := make([]reflect.SelectCase, len(m.inputs)+1)
        cases[0] = reflect.SelectCase{
            Dir:  reflect.SelectRecv,
            Chan: reflect.ValueOf(m.done),
        }
        
        for i, ch := range m.inputs {
            cases[i+1] = reflect.SelectCase{
                Dir:  reflect.SelectRecv,
                Chan: reflect.ValueOf(ch),
            }
        }
        
        timer := time.NewTimer(m.timeout)
        defer timer.Stop()
        
        for len(cases) > 1 {
            chosen, value, ok := reflect.Select(cases)
            
            if !ok {
                // Channel closed, remove it
                cases = append(cases[:chosen], cases[chosen+1:]...)
                continue
            }
            
            if chosen == 0 {
                // Done signal received
                return
            }
            
            // Reset timer
            timer.Reset(m.timeout)
            
            select {
            case m.output <- value.Interface().(T):
                // Value sent successfully
            case <-timer.C:
                log.Println("Timeout while sending value")
            }
        }
    }()
    
    return m.output
}
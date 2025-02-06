// Example 28
type DataProcessor struct {
    input  <-chan int    // Read-only channel
    output chan<- int    // Write-only channel
    done   chan struct{} // Bidirectional signal channel
}

func NewDataProcessor(input <-chan int, output chan<- int) *DataProcessor {
    return &DataProcessor{
        input:  input,
        output: output,
        done:   make(chan struct{}),
    }
}

func (dp *DataProcessor) Process() {
    defer close(dp.done)
    
    for value := range dp.input {
        result := value * 2
        select {
        case dp.output <- result:
            // Value sent successfully
        default:
            // Handle backpressure
            log.Printf("Output channel full, dropping value: %d", result)
        }
    }
}

func (dp *DataProcessor) Wait() {
    <-dp.done
}
// Example 28
package main

import (
	"fmt"
	"log"
	"time"
)

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

func main() {
	// Create channels for input and output
	input := make(chan int, 5)
	output := make(chan int, 3) // Small buffer to demonstrate backpressure

	// Create and start the processor
	processor := NewDataProcessor(input, output)
	go processor.Process()

	// Send some data to input
	go func() {
		for i := 1; i <= 10; i++ {
			input <- i
			time.Sleep(100 * time.Millisecond)
		}
		close(input) // Signal that we're done sending data
	}()

	// Receive from output (intentionally slow to demonstrate backpressure)
	go func() {
		for result := range output {
			fmt.Printf("Received result: %d\n", result)
			time.Sleep(300 * time.Millisecond) // Slow consumer
		}
	}()

	// Wait for processing to complete
	processor.Wait()
	fmt.Println("Processing complete")

	// Close output channel and allow time for final messages to print
	close(output)
	time.Sleep(500 * time.Millisecond)
}
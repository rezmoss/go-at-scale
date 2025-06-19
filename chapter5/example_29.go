// Example 29
package main

import (
	"fmt"
	"runtime"
	"time"
)

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

func NewAdaptiveProcessor[T any](initialCapacity int) *AdaptiveProcessor[T] {
	return &AdaptiveProcessor[T]{
		input:    make(chan T, initialCapacity),
		output:   make(chan T, initialCapacity),
		capacity: initialCapacity,
	}
}

// This example demonstrates resizing logic but DOES NOT actually resize
// the input channel as that's not safely possible in Go
func (ap *AdaptiveProcessor[T]) checkLoad() {
	currentLoad := float64(len(ap.input)) / float64(cap(ap.input))
	if currentLoad > 0.8 && ap.capacity < 1000 {
		// In a real implementation, we'd need to create a new channel
		// and transfer items, but for demonstration purposes, we just
		// track what would happen
		oldCapacity := ap.capacity
		ap.capacity *= 2
		fmt.Printf("Buffer would resize from %d to %d (load: %.2f%%)\n",
			oldCapacity, ap.capacity, currentLoad*100)
	}
}

func main() {
	// Example 1: Using NewBuffer
	fmt.Println("Example 1: NewBuffer")

	// Create a static buffer
	staticBuf := NewBuffer[int](BufferConfig{Size: 5, Dynamic: false})
	fmt.Printf("Static buffer capacity: %d\n", cap(staticBuf))

	// Create a dynamic buffer based on GOMAXPROCS
	dynamicBuf := NewBuffer[int](BufferConfig{Size: 0, Dynamic: true})
	fmt.Printf("Dynamic buffer capacity: %d (based on GOMAXPROCS)\n", cap(dynamicBuf))

	// Example 2: AdaptiveProcessor demonstration
	fmt.Println("\nExample 2: AdaptiveProcessor")

	// Create test data
	testData := make([]int, 10)
	for i := range testData {
		testData[i] = i + 1
	}

	// Create adaptive processor
	processor := NewAdaptiveProcessor[int](4)
	fmt.Printf("Initial capacity: %d\n", processor.capacity)

	// Start processing in a separate goroutine
	go func() {
		for i, item := range testData {
			select {
			case processor.input <- item:
				fmt.Printf("Added item %d to input buffer\n", item)
				// Check load after each item
				processor.checkLoad()
			default:
				// Buffer is full, demonstrate backpressure
				fmt.Printf("Buffer full, waiting to add item %d...\n", item)
				processor.input <- item // This will block until space is available
				fmt.Printf("Added item %d after waiting\n", item)
				processor.checkLoad()
			}

			// Only close after all items are processed
			if i == len(testData)-1 {
				close(processor.input)
			}
		}
	}()

	// Consumer goroutine
	go func() {
		for item := range processor.input {
			// Simulate processing
			time.Sleep(20 * time.Millisecond)
			processor.output <- item

			fmt.Printf("Processed item: %d\n", item)
		}
		close(processor.output)
	}()

	// Collect and display results
	processedItems := []int{}
	for item := range processor.output {
		processedItems = append(processedItems, item)
	}

	fmt.Printf("\nFinal processed items (%d total): %v\n", len(processedItems), processedItems)
	fmt.Printf("Final theoretical capacity: %d\n", processor.capacity)
}
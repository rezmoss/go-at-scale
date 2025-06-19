// Example 25
package main

import (
	"fmt"
	"sync"
	"time"
)

func fanOut[T any](source <-chan T, workers int) []<-chan T {
	channels := make([]<-chan T, workers)
	for i := 0; i < workers; i++ {
		ch := make(chan T)
		channels[i] = ch

		go func(ch chan<- T) {
			defer close(ch)
			for item := range source {
				ch <- item
			}
		}(ch)
	}
	return channels
}

func fanIn[T any](channels []<-chan T) <-chan T {
	merged := make(chan T)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(ch <-chan T) {
			defer wg.Done()
			for item := range ch {
				merged <- item
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

// Example: Processing pipeline with fan-out/fan-in
func ProcessItems[T any](items []T, processor func(T) T, workers int) []T {
	// Create source channel
	source := make(chan T)
	go func() {
		defer close(source)
		for _, item := range items {
			source <- item
		}
	}()

	// Fan out processing
	channels := fanOut(source, workers)

	// Process items in parallel
	processedChannels := make([]<-chan T, workers)
	for i, ch := range channels {
		processedCh := make(chan T)
		processedChannels[i] = processedCh

		go func(in <-chan T, out chan<- T) {
			defer close(out)
			for item := range in {
				out <- processor(item)
			}
		}(ch, processedCh)
	}

	// Fan in results
	results := fanIn(processedChannels)

	// Collect results
	var processed []T
	for result := range results {
		processed = append(processed, result)
	}

	return processed
}

// Simple processor function that simulates work
func slowProcessor(n int) int {
	// Simulate processing work
	time.Sleep(100 * time.Millisecond)
	return n * 2
}

func main() {
	// Create a slice of items to process
	items := make([]int, 100)
	for i := range items {
		items[i] = i + 1
	}

	fmt.Println("Starting processing with fan-out/fan-in pattern...")
	start := time.Now()

	// Process items using fan-out/fan-in with 10 workers
	results := ProcessItems(items, slowProcessor, 10)

	elapsed := time.Since(start)
	fmt.Printf("Processed %d items in %v\n", len(results), elapsed)

	// Verify a few results
	fmt.Println("Sample results:")
	for i := 0; i < 5 && i < len(results); i++ {
		fmt.Printf("Result %d: %d\n", i, results[i])
	}
}
// Example 33
package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Result[T any] struct {
	Value T
	Err   error
}

type ErrorCollector struct {
	errors []error
	mu     sync.Mutex
}

func (ec *ErrorCollector) Add(err error) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.errors = append(ec.errors, err)
}

func (ec *ErrorCollector) Error() error {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	if len(ec.errors) == 0 {
		return nil
	}

	errStrings := make([]string, len(ec.errors))
	for i, err := range ec.errors {
		errStrings[i] = err.Error()
	}

	return fmt.Errorf("multiple errors occurred: %s", strings.Join(errStrings, "; "))
}

// ConcurrentProcess processes items concurrently with error handling
func ConcurrentProcess[T any](items []T, process func(T) error) error {
	collector := &ErrorCollector{}
	var wg sync.WaitGroup

	for _, item := range items {
		wg.Add(1)
		go func(item T) {
			defer wg.Done()
			if err := process(item); err != nil {
				collector.Add(err)
			}
		}(item)
	}

	wg.Wait()
	return collector.Error()
}

func main() {
	// Sample data to process
	urls := []string{
		"https://example.com/1",
		"https://example.com/2",
		"https://invalid-url", // Will cause an error
		"https://example.com/4",
		"https://another-invalid-url", // Will cause an error
	}

	// Process function that simulates HTTP fetching with potential errors
	process := func(url string) error {
		// Simulate processing time
		time.Sleep(100 * time.Millisecond)

		// Simulate errors for certain URLs
		if strings.Contains(url, "invalid") {
			return fmt.Errorf("failed to fetch %s", url)
		}

		fmt.Printf("Successfully processed: %s\n", url)
		return nil
	}

	// Run concurrent processing with error handling
	err := ConcurrentProcess(urls, process)

	// Handle the aggregated errors
	if err != nil {
		fmt.Printf("Error summary: %v\n", err)
	} else {
		fmt.Println("All items processed successfully")
	}
}
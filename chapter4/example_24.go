// Example 24
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

// processItem simulates processing a single item with potential for error
func processItem(ctx context.Context, item int) error {
	// Simulate work
	time.Sleep(time.Duration(100*item) * time.Millisecond)

	// Simulate error for item 3
	if item == 3 {
		return fmt.Errorf("error processing item %d", item)
	}

	fmt.Printf("Processed item %d\n", item)
	return nil
}

func main() {
	// Create a context
	ctx := context.Background()

	// Create items to process
	items := []int{1, 2, 3, 4, 5}

	// Create an error group with context
	g, ctx := errgroup.WithContext(ctx)

	for _, item := range items {
		it := item // local copy
		g.Go(func() error {
			return processItem(ctx, it)
		})
	}

	// If any goroutine returns an error, g.Wait() returns it.
	if err := g.Wait(); err != nil {
		// Handle the first error from any goroutine
		log.Printf("Error processing items: %v", err)
	} else {
		// All goroutines succeeded
		fmt.Println("All items processed successfully")
	}
}
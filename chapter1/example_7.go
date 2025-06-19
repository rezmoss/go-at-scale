// Example 7
package main

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("resource not found")

func fetchResource(id string) error {
	// Simplified implementation that always returns the error
	return fmt.Errorf("fetching resource %s: %w", id, ErrNotFound)
}

func main() {
	// Try to fetch a resource
	err := fetchResource("abc123")

	// Handle and print the error
	if err != nil {
		fmt.Println("Error:", err)

		// Check if the error is ErrNotFound using errors.Is
		if errors.Is(err, ErrNotFound) {
			fmt.Println("Specific error detected: resource not found")
		}

		// Unwrap to get the original error
		fmt.Println("Unwrapped error:", errors.Unwrap(err))
	}
}
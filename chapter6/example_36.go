// Example 36
package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

type ContextAwareService struct {
	operations chan func()
	errors     chan error
}

func NewContextAwareService() *ContextAwareService {
	return &ContextAwareService{
		operations: make(chan func()),
		errors:     make(chan error, 1),
	}
}

func (s *ContextAwareService) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case op := <-s.operations:
				op()
			case <-ctx.Done():
				s.errors <- ctx.Err()
				close(s.operations)
				return
			}
		}
	}()
}

// Example: Timeout-aware operations
func (s *ContextAwareService) ExecuteWithTimeout(
	ctx context.Context,
	operation func() error,
	timeout time.Duration,
) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan error, 1)

	go func() {
		done <- operation()
	}()

	select {
	case err := <-done:
		return err
	case <-timeoutCtx.Done():
		return fmt.Errorf("operation timed out: %w", timeoutCtx.Err())
	}
}

func main() {
	// Create a new service
	service := NewContextAwareService()

	// Create a parent context with cancel functionality
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure we cancel the context when main exits

	// Start the service with the context
	service.Start(ctx)

	// Create a slow operation that will take 2 seconds
	slowOperation := func() error {
		time.Sleep(2 * time.Second)
		fmt.Println("Slow operation completed successfully")
		return nil
	}

	// Create a faster operation
	fastOperation := func() error {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("Fast operation completed successfully")
		return nil
	}

	// Execute with timeout examples
	fmt.Println("Example 1: Fast operation with sufficient timeout")
	err := service.ExecuteWithTimeout(ctx, fastOperation, 1*time.Second)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	fmt.Println("\nExample 2: Slow operation with insufficient timeout")
	err = service.ExecuteWithTimeout(ctx, slowOperation, 1*time.Second)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	fmt.Println("\nExample 3: Slow operation with sufficient timeout")
	err = service.ExecuteWithTimeout(ctx, slowOperation, 3*time.Second)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	// Wait a moment before canceling the parent context
	fmt.Println("\nExample 4: Demonstrating parent context cancellation")
	fmt.Println("Canceling parent context in 1 second...")
	time.Sleep(1 * time.Second)
	cancel()

	// Check error from service
	select {
	case err := <-service.errors:
		fmt.Printf("Service shutdown with error: %v\n", err)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("No error received from service")
	}

	// Give time for final output to appear
	time.Sleep(100 * time.Millisecond)
}
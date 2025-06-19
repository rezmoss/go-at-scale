// Example 50
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Basic interface
type Handler interface {
	Handle(req *http.Request) error
}

// Base handler
type BaseHandler struct{}

func (h *BaseHandler) Handle(req *http.Request) error {
	// Basic handling
	fmt.Println("Base handler: Processing request")
	return nil
}

// Logging decorator
type LoggingDecorator struct {
	handler Handler
	logger  *log.Logger
}

func (d *LoggingDecorator) Handle(req *http.Request) error {
	start := time.Now()
	err := d.handler.Handle(req)
	d.logger.Printf("Request processed in %v", time.Since(start))
	return err
}

// Retry decorator
type RetryDecorator struct {
	handler Handler
	retries int
}

func (d *RetryDecorator) Handle(req *http.Request) (err error) {
	for i := 0; i <= d.retries; i++ {
		if err = d.handler.Handle(req); err == nil {
			return nil
		}
		fmt.Printf("Retry %d after error: %v\n", i+1, err)
		time.Sleep(time.Second << uint(i)) // Exponential backoff
	}
	return err
}

// ErrorHandler is a base handler that simulates errors for testing the retry mechanism
type ErrorHandler struct {
	failCount int
	attempts  int
}

func (h *ErrorHandler) Handle(req *http.Request) error {
	h.attempts++
	if h.attempts <= h.failCount {
		return fmt.Errorf("simulated error on attempt %d", h.attempts)
	}
	fmt.Println("Error handler: Successfully processed after", h.attempts, "attempts")
	return nil
}

func main() {
	// Create a simple request for testing
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	// Setup logger
	logger := log.New(os.Stdout, "[LOGGER] ", log.Ltime)

	fmt.Println("=== Testing Base Handler ===")
	baseHandler := &BaseHandler{}
	baseHandler.Handle(req)

	fmt.Println("\n=== Testing Logging Decorator ===")
	loggingHandler := &LoggingDecorator{
		handler: baseHandler,
		logger:  logger,
	}
	loggingHandler.Handle(req)

	fmt.Println("\n=== Testing Retry Decorator ===")
	errorHandler := &ErrorHandler{failCount: 2}
	retryHandler := &RetryDecorator{
		handler: errorHandler,
		retries: 3,
	}
	err := retryHandler.Handle(req)
	if err != nil {
		fmt.Printf("Failed after retries: %v\n", err)
	} else {
		fmt.Println("Successfully processed with retries")
	}

	fmt.Println("\n=== Testing Stacked Decorators ===")
	// Create a new error handler that will fail twice
	errorHandler = &ErrorHandler{failCount: 2}

	// Stack decorators: first retry, then log
	stackedHandler := &LoggingDecorator{
		handler: &RetryDecorator{
			handler: errorHandler,
			retries: 3,
		},
		logger: logger,
	}

	err = stackedHandler.Handle(req)
	if err != nil {
		fmt.Printf("Failed with stacked decorators: %v\n", err)
	} else {
		fmt.Println("Successfully processed with stacked decorators")
	}
}
// Example 77
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// State represents the state of the circuit breaker
type State int

const (
	StateClosed   State = iota // Circuit is closed and requests are allowed
	StateOpen                  // Circuit is open and requests are blocked
	StateHalfOpen              // Circuit is half-open, allowing a test request
)

// Errors
var (
	ErrCircuitOpen = errors.New("circuit breaker is open")
)

// MetricsRecorder interface for recording metrics
type MetricsRecorder interface {
	IncCounter(name string, labels ...string)
}

// Logger interface for logging
type Logger interface {
	Warn(msg string, keyvals ...interface{})
}

// Operation represents a function that can be executed
type Operation func(ctx context.Context) error

// Simple implementations of MetricsRecorder and Logger for demo purposes
type SimpleMetrics struct{}

func (sm *SimpleMetrics) IncCounter(name string, labels ...string) {
	fmt.Printf("Metric: %s, Labels: %v\n", name, labels)
}

type SimpleLogger struct{}

func (sl *SimpleLogger) Warn(msg string, keyvals ...interface{}) {
	fmt.Printf("WARNING: %s, %v\n", msg, keyvals)
}

// CircuitBreaker implementation
type CircuitBreaker struct {
	name         string
	maxFailures  int
	resetTimeout time.Duration
	failureCount int64
	lastFailure  time.Time
	state        State
	metrics      MetricsRecorder
	mutex        sync.RWMutex
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(name string, maxFailures int, resetTimeout time.Duration, metrics MetricsRecorder) *CircuitBreaker {
	return &CircuitBreaker{
		name:         name,
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        StateClosed,
		metrics:      metrics,
	}
}

func (cb *CircuitBreaker) allowRequest() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	if cb.state == StateClosed {
		return true
	}

	if cb.state == StateOpen {
		// Check if reset timeout has elapsed
		if time.Since(cb.lastFailure) > cb.resetTimeout {
			// Transition to half-open state
			cb.mutex.RUnlock()
			cb.mutex.Lock()
			cb.state = StateHalfOpen
			cb.mutex.Unlock()
			cb.mutex.RLock()
			return true
		}
		return false
	}

	// State is half-open, allow one test request
	return true
}

func (cb *CircuitBreaker) recordFailure() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	cb.failureCount++
	cb.lastFailure = time.Now()

	if cb.state == StateHalfOpen || cb.failureCount >= int64(cb.maxFailures) {
		cb.state = StateOpen
		fmt.Printf("Circuit breaker '%s' has opened after %d failures\n", cb.name, cb.failureCount)
	}
}

func (cb *CircuitBreaker) recordSuccess() {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if cb.state == StateHalfOpen {
		cb.state = StateClosed
		cb.failureCount = 0
		fmt.Printf("Circuit breaker '%s' has closed after successful test request\n", cb.name)
	} else if cb.state == StateClosed && cb.failureCount > 0 {
		cb.failureCount--
	}
}

// Execute applies the circuit breaker pattern to the given operation
func (cb *CircuitBreaker) Execute(operation func() error) error {
	if !cb.allowRequest() {
		cb.metrics.IncCounter("circuit_breaker_rejections", "name", cb.name)
		fmt.Printf("Circuit breaker '%s' is open, rejecting request\n", cb.name)
		return ErrCircuitOpen
	}

	err := operation()

	if err != nil {
		cb.recordFailure()
		return err
	}

	cb.recordSuccess()
	return nil
}

// Fallback implementation
type Fallback struct {
	primary  Operation
	fallback Operation
	metrics  MetricsRecorder
	logger   Logger
}

// NewFallback creates a new fallback
func NewFallback(primary, fallback Operation, metrics MetricsRecorder, logger Logger) *Fallback {
	return &Fallback{
		primary:  primary,
		fallback: fallback,
		metrics:  metrics,
		logger:   logger,
	}
}

func (f *Fallback) Execute(ctx context.Context) error {
	err := f.primary(ctx)
	if err == nil {
		return nil
	}

	f.metrics.IncCounter("fallback_triggered")
	f.logger.Warn("primary operation failed, using fallback", "error", err)

	return f.fallback(ctx)
}

// Main function to demonstrate the circuit breaker and fallback patterns
func main() {
	// Create simple metrics and logger
	metrics := &SimpleMetrics{}
	logger := &SimpleLogger{}

	// Create a circuit breaker with low threshold to demonstrate opening
	cb := NewCircuitBreaker("service-a", 3, 5*time.Second, metrics)

	// Print initial state
	fmt.Println("Circuit breaker created in CLOSED state")

	// Create a primary operation that always fails to demonstrate circuit breaker
	primaryOp := func(ctx context.Context) error {
		// Always fail to demonstrate circuit breaker opening
		return errors.New("service unavailable")
	}

	// Create a fallback operation
	fallbackOp := func(ctx context.Context) error {
		fmt.Println("Fallback operation executed")
		return nil
	}

	// Create a fallback
	fb := NewFallback(primaryOp, fallbackOp, metrics, logger)

	// Simulate multiple requests to demonstrate both patterns
	fmt.Println("Starting simulation...")
	for i := 0; i < 10; i++ {
		fmt.Printf("\nRequest %d:\n", i+1)
		fmt.Printf("Circuit state before request: %v, Failure count: %d\n", cb.state, cb.failureCount)

		// Execute primary operation with circuit breaker protection
		err := cb.Execute(func() error {
			// This simulates the primary service call
			return errors.New("service unavailable")
		})

		if err != nil {
			if errors.Is(err, ErrCircuitOpen) {
				fmt.Println("Circuit is open, request rejected")
			} else {
				// If circuit allowed the request but primary failed, use fallback
				fmt.Printf("Primary failed: %v\n", err)
				err = fb.fallback(context.Background())
				if err != nil {
					fmt.Printf("Fallback also failed: %v\n", err)
				}
			}
		}

		// Add longer pause after the circuit opens to demonstrate reset timeout
		if i == 5 {
			fmt.Println("\nWaiting for circuit breaker reset timeout...")
			time.Sleep(6 * time.Second) // Slightly longer than resetTimeout
		} else {
			time.Sleep(500 * time.Millisecond)
		}
	}
}
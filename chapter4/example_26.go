// Example 26
package main

import (
	"context"
	"log"
	"sync"
	"time"
)

type RateLimiter struct {
	ticker *time.Ticker
	limit  chan struct{}
}

func NewRateLimiter(rate int, burst int) *RateLimiter {
	limiter := &RateLimiter{
		ticker: time.NewTicker(time.Second / time.Duration(rate)),
		limit:  make(chan struct{}, burst),
	}

	// Fill token bucket
	for i := 0; i < burst; i++ {
		limiter.limit <- struct{}{}
	}

	// Replenish tokens
	go func() {
		for range limiter.ticker.C {
			select {
			case limiter.limit <- struct{}{}:
			default:
				// Bucket is full
			}
		}
	}()

	return limiter
}

// Wait blocks until a token is available or the context is cancelled.
func (r *RateLimiter) Wait(ctx context.Context) error {
	select {
	case <-r.limit:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func main() {
	limiter := NewRateLimiter(10, 5) // 10 ops/sec, burst of 5

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := limiter.Wait(ctx)
			if err != nil {
				log.Printf("Operation %d failed: %v", i, err)
				return
			}
			log.Printf("Operation %d executed", i)
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
}
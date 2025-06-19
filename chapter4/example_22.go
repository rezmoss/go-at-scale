// Example 22
package main

import (
	"context"
	"log"
	"time"
)

// Worker struct manages the lifecycle of a goroutine
type Worker struct {
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
}

// NewWorker creates a new Worker with a derived context
func NewWorker(ctx context.Context) *Worker {
	ctx, cancel := context.WithCancel(ctx)
	return &Worker{
		ctx:    ctx,
		cancel: cancel,
		done:   make(chan struct{}),
	}
}

// Start begins the worker's execution in a separate goroutine
func (w *Worker) Start() {
	go func() {
		defer close(w.done)

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-w.ctx.Done():
				log.Println("Worker: shutdown signal received")
				return
			case <-ticker.C:
				w.doWork()
			}
		}
	}()
}

// Stop signals the worker to stop and waits for it to finish
func (w *Worker) Stop() {
	w.cancel()
	<-w.done // Wait for worker to finish
}

// doWork simulates a work task
func (w *Worker) doWork() {
	// Simulate some work
	log.Println("Worker: doing work...")
	time.Sleep(100 * time.Millisecond)
}

func main() {
	// Create a base context
	ctx := context.Background()

	// Create and start a worker
	worker := NewWorker(ctx)
	worker.Start()

	log.Println("Main: worker started, running for 3 seconds...")

	// Let the worker run for a few seconds
	time.Sleep(3 * time.Second)

	// Stop the worker and wait for it to finish
	log.Println("Main: stopping worker...")
	worker.Stop()

	log.Println("Main: worker stopped, program exiting")
}
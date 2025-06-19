// Example 35
package main

import (
	"fmt"
	"sync"
	"time"
)

type BatchProcessor[T any] struct {
	maxWorkers int
	tasks      []T
	process    func(T) error
}

func NewBatchProcessor[T any](maxWorkers int, process func(T) error) *BatchProcessor[T] {
	return &BatchProcessor[T]{
		maxWorkers: maxWorkers,
		process:    process,
	}
}

func (bp *BatchProcessor[T]) AddTask(task T) {
	bp.tasks = append(bp.tasks, task)
}

func (bp *BatchProcessor[T]) Process() error {
	var (
		wg       sync.WaitGroup
		errMu    sync.Mutex
		firstErr error
	)

	// Create task channel
	taskCh := make(chan T)

	// Start workers
	for i := 0; i < bp.maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskCh {
				if err := bp.process(task); err != nil {
					errMu.Lock()
					if firstErr == nil {
						firstErr = err
					}
					errMu.Unlock()
				}
			}
		}()
	}

	// Send tasks
	for _, task := range bp.tasks {
		taskCh <- task
	}
	close(taskCh)

	// Wait for completion
	wg.Wait()
	return firstErr
}

// Custom wrapper to make the example work
type CustomWaitGroup struct {
	sync.WaitGroup
	count int
}

func (cwg *CustomWaitGroup) Add(delta int) {
	cwg.count += delta
	cwg.WaitGroup.Add(delta)
}

func (cwg *CustomWaitGroup) Done() {
	cwg.count--
	cwg.WaitGroup.Done()
}

func (cwg *CustomWaitGroup) GetCount() int {
	return cwg.count
}

// Example usage with dynamic scaling
type DynamicWorkerPool[T any] struct {
	minWorkers int
	maxWorkers int
	taskQueue  chan T
	process    func(T) error
	wg         CustomWaitGroup
}

func (p *DynamicWorkerPool[T]) scaleWorkers() {
	currentLoad := float64(len(p.taskQueue)) / float64(cap(p.taskQueue))

	if currentLoad > 0.8 && p.wg.GetCount() < p.maxWorkers {
		// Add workers
		toAdd := (p.maxWorkers - p.wg.GetCount()) / 2
		for i := 0; i < toAdd; i++ {
			p.startWorker()
		}
	}
}

// Added to make it work
func (p *DynamicWorkerPool[T]) startWorker() {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for task := range p.taskQueue {
			_ = p.process(task)
		}
	}()
}

func NewDynamicWorkerPool[T any](minWorkers, maxWorkers int, queueSize int, process func(T) error) *DynamicWorkerPool[T] {
	pool := &DynamicWorkerPool[T]{
		minWorkers: minWorkers,
		maxWorkers: maxWorkers,
		taskQueue:  make(chan T, queueSize),
		process:    process,
	}

	// Start minimum workers
	for i := 0; i < minWorkers; i++ {
		pool.startWorker()
	}

	return pool
}

func (p *DynamicWorkerPool[T]) AddTask(task T) {
	p.taskQueue <- task
	p.scaleWorkers()
}

func (p *DynamicWorkerPool[T]) Close() {
	close(p.taskQueue)
	p.wg.Wait()
}

func main() {
	// Demo of BatchProcessor
	processor := NewBatchProcessor(3, func(task int) error {
		fmt.Printf("Processing task %d\n", task)
		time.Sleep(100 * time.Millisecond)
		return nil
	})

	// Add some tasks
	for i := 1; i <= 10; i++ {
		processor.AddTask(i)
	}

	fmt.Println("Starting batch processing...")
	err := processor.Process()
	if err != nil {
		fmt.Printf("Error during processing: %v\n", err)
	} else {
		fmt.Println("Batch processing completed successfully")
	}

	// Demo of DynamicWorkerPool
	fmt.Println("\nStarting dynamic worker pool...")
	pool := NewDynamicWorkerPool[int](2, 5, 10, func(task int) error {
		fmt.Printf("Worker processing task %d\n", task)
		time.Sleep(200 * time.Millisecond)
		return nil
	})

	// Add tasks to the pool
	for i := 1; i <= 8; i++ {
		pool.AddTask(i)
	}

	// Add a few more tasks after a delay to demonstrate scaling
	time.Sleep(300 * time.Millisecond)
	for i := 9; i <= 12; i++ {
		pool.AddTask(i)
	}

	// Wait a bit to see tasks being processed
	time.Sleep(500 * time.Millisecond)

	// Close the pool and wait for all tasks to complete
	pool.Close()
	fmt.Println("Dynamic worker pool completed")
}
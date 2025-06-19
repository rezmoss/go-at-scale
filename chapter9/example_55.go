// Example 55
package main

import (
	"fmt"
	"sync"
)

type Worker struct {
	buffer []byte
}

// Pool for reusing workers
var workerPool = sync.Pool{
	New: func() interface{} {
		return &Worker{
			buffer: make([]byte, 1024),
		}
	},
}

func ProcessRequest(data []byte) error {
	worker := workerPool.Get().(*Worker)
	defer workerPool.Put(worker)

	// Reset buffer
	worker.buffer = worker.buffer[:0]

	// Use worker - added simple operation to demonstrate usage
	worker.buffer = append(worker.buffer, data...)
	fmt.Printf("Worker processed %d bytes of data\n", len(worker.buffer))

	return nil
}

// Define Item and Result types which were referenced but not defined
type Item struct {
	ID   int
	Data []byte
}

type Result struct {
	ItemID int
	Status string
}

func processBatch(items []Item) []Result {
	results := make([]Result, 0, len(items))
	for _, item := range items {
		results = append(results, Result{
			ItemID: item.ID,
			Status: fmt.Sprintf("Processed %d bytes", len(item.Data)),
		})
	}
	return results
}

// Efficient slice handling
type BatchProcessor struct {
	batchSize int
	items     []Item
}

func (bp *BatchProcessor) ProcessBatch() []Result {
	// Pre-allocate slice with capacity
	results := make([]Result, 0, len(bp.items))

	for i := 0; i < len(bp.items); i += bp.batchSize {
		end := i + bp.batchSize
		if end > len(bp.items) {
			end = len(bp.items)
		}

		batch := bp.items[i:end]
		results = append(results, processBatch(batch)...)
	}

	return results
}

func main() {
	// Example of using ProcessRequest
	fmt.Println("Worker Pool Example:")
	sampleData := []byte("This is some sample data to process")
	err := ProcessRequest(sampleData)
	if err != nil {
		fmt.Printf("Error processing request: %v\n", err)
	}

	// Example of using BatchProcessor
	fmt.Println("\nBatch Processing Example:")
	items := []Item{
		{ID: 1, Data: []byte("Item 1 data")},
		{ID: 2, Data: []byte("Item 2 data is a bit longer")},
		{ID: 3, Data: []byte("Item 3")},
		{ID: 4, Data: []byte("Item 4 has some extra info")},
		{ID: 5, Data: []byte("Item 5 data")},
	}

	processor := BatchProcessor{
		batchSize: 2,
		items:     items,
	}

	results := processor.ProcessBatch()
	fmt.Println("Processed results:")
	for i, result := range results {
		fmt.Printf("  %d. ItemID: %d, Status: %s\n", i+1, result.ItemID, result.Status)
	}
}
// Example 55
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
    
    // Use worker...
    return nil
}

// Efficient slice handling
type BatchProcessor struct {
    batchSize int
    items     []Item
}

func (bp *BatchProcessor) ProcessBatch() {
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
}
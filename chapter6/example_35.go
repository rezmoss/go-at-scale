// Example 35
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

// Example usage with dynamic scaling
type DynamicWorkerPool[T any] struct {
    minWorkers int
    maxWorkers int
    taskQueue  chan T
    process    func(T) error
    wg         sync.WaitGroup
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
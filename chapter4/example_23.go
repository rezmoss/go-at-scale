// Example 23
type Job[In, Out any] struct {
    Input  In
    Result chan Out
}

type WorkerPool[In, Out any] struct {
    workers  int
    jobs     chan Job[In, Out]
    done     chan struct{}
    processor func(In) Out
}

func NewWorkerPool[In, Out any](workers int, processor func(In) Out) *WorkerPool[In, Out] {
    return &WorkerPool[In, Out]{
        workers:   workers,
        jobs:      make(chan Job[In, Out]),
        done:      make(chan struct{}),
        processor: processor,
    }
}

func (p *WorkerPool[In, Out]) Start() {
    // Start workers
    for i := 0; i < p.workers; i++ {
        go func(workerID int) {
            for job := range p.jobs {
                result := p.processor(job.Input)
                job.Result <- result
            }
        }(i)
    }
}

func (p *WorkerPool[In, Out]) Submit(ctx context.Context, input In) (Out, error) {
    resultChan := make(chan Out, 1)
    select {
    case p.jobs <- Job[In, Out]{Input: input, Result: resultChan}:
        select {
        case result := <-resultChan:
            return result, nil
        case <-ctx.Done():
            return *new(Out), ctx.Err()
        }
    case <-ctx.Done():
        return *new(Out), ctx.Err()
    }
}

// Example usage
func main() {
    // Create a worker pool that processes integers
    pool := NewWorkerPool(5, func(x int) int {
        time.Sleep(100 * time.Millisecond) // Simulate work
        return x * 2
    })
    
    pool.Start()
    
    // Process multiple items concurrently
    results := make([]int, 10)
    var wg sync.WaitGroup
    
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            results[i] = pool.Submit(i)
        }(i)
    }
    
    wg.Wait()
}
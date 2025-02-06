// Example 22
type Worker struct {
    ctx    context.Context
    cancel context.CancelFunc
    done   chan struct{}
}

func NewWorker(ctx context.Context) *Worker {
    ctx, cancel := context.WithCancel(ctx)
    return &Worker{
        ctx:    ctx,
        cancel: cancel,
        done:   make(chan struct{}),
    }
}

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

func (w *Worker) Stop() {
    w.cancel()
    <-w.done // Wait for worker to finish
}

func (w *Worker) doWork() {
    // Simulate some work
    time.Sleep(100 * time.Millisecond)
}
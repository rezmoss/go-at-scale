// Example 32
type Process struct {
    done chan struct{}
    err  error
    once sync.Once
}

func NewProcess() *Process {
    return &Process{
        done: make(chan struct{}),
    }
}

func (p *Process) Stop(err error) {
    p.once.Do(func() {
        p.err = err
        close(p.done)
    })
}

func (p *Process) Wait() error {
    <-p.done
    return p.err
}

// Example: Graceful shutdown
type Server struct {
    processes []*Process
    done      chan struct{}
}

func (s *Server) Shutdown(ctx context.Context) error {
    // Signal shutdown to all processes
    close(s.done)
    
    // Wait for all processes with timeout
    for _, proc := range s.processes {
        select {
        case <-proc.done:
            // Process completed
        case <-ctx.Done():
            return ctx.Err()
        }
    }
    
    return nil
}
// Example 50
// Basic interface
type Handler interface {
    Handle(req *http.Request) error
}

// Base handler
type BaseHandler struct{}

func (h *BaseHandler) Handle(req *http.Request) error {
    // Basic handling
    return nil
}

// Logging decorator
type LoggingDecorator struct {
    handler Handler
    logger  *log.Logger
}

func (d *LoggingDecorator) Handle(req *http.Request) error {
    start := time.Now()
    err := d.handler.Handle(req)
    d.logger.Printf("Request processed in %v", time.Since(start))
    return err
}

// Retry decorator
type RetryDecorator struct {
    handler Handler
    retries int
}

func (d *RetryDecorator) Handle(req *http.Request) (err error) {
    for i := 0; i <= d.retries; i++ {
        if err = d.handler.Handle(req); err == nil {
            return nil
        }
        time.Sleep(time.Second << uint(i)) // Exponential backoff
    }
    return err
}
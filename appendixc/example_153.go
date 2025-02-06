// Example 153
// internal/infrastructure/messaging/router.go
type MessageRouter struct {
    routes   map[string][]MessageHandler
    filters  []MessageFilter
    logger   Logger
    metrics  MetricsRecorder
}

func (r *MessageRouter) Route(ctx context.Context, message Message) error {
    // Apply filters
    for _, filter := range r.filters {
        if !filter.Accept(message) {
            r.metrics.IncCounter("messages_filtered")
            return nil
        }
    }

    // Find handlers
    handlers, exists := r.routes[message.Type]
    if !exists {
        r.metrics.IncCounter("messages_unrouted")
        return ErrNoHandlerFound
    }

    // Execute handlers
    for _, handler := range handlers {
        if err := handler.Handle(ctx, message); err != nil {
            return fmt.Errorf("handling message: %w", err)
        }
    }

    return nil
}
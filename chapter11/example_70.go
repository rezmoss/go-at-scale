// Example 70
type CommandHandler interface {
    Handle(ctx context.Context, cmd Command) error
}

type QueryHandler interface {
    Handle(ctx context.Context, query Query) (interface{}, error)
}

type UserCommandHandler struct {
    eventStore  EventStore
    publisher   EventPublisher
    validator   Validator
    metrics     MetricsRecorder
}

func (h *UserCommandHandler) Handle(ctx context.Context, cmd Command) error {
    start := time.Now()
    defer func() {
        h.metrics.ObserveLatency("command_processing", time.Since(start))
    }()

    // Validate command
    if err := h.validator.Validate(cmd); err != nil {
        return fmt.Errorf("invalid command: %w", err)
    }

    // Generate events
    events, err := cmd.ToEvents()
    if err != nil {
        return fmt.Errorf("generating events: %w", err)
    }

    // Store events
    if err := h.eventStore.Save(ctx, events); err != nil {
        return fmt.Errorf("saving events: %w", err)
    }

    // Publish events
    for _, event := range events {
        if err := h.publisher.Publish(ctx, event); err != nil {
            h.metrics.IncCounter("event_publish_errors")
            h.logger.Error("failed to publish event", 
                "error", err,
                "event_id", event.ID)
        }
    }

    return nil
}
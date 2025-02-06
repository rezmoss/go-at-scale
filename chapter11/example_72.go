// Example 72
type EventProcessor struct {
    consumer   kafka.Consumer
    handler    EventHandler
    metrics    MetricsRecorder
    logger     Logger
    errorsChan chan error
}

func (p *EventProcessor) Start(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            msg, err := p.consumer.ReadMessage(ctx)
            if err != nil {
                p.metrics.IncCounter("message_read_errors")
                p.logger.Error("failed to read message", "error", err)
                continue
            }

            if err := p.processMessage(ctx, msg); err != nil {
                p.metrics.IncCounter("message_processing_errors")
                p.logger.Error("failed to process message",
                    "error", err,
                    "message_id", msg.ID)
                p.errorsChan <- err
            }
        }
    }
}

func (p *EventProcessor) processMessage(ctx context.Context, msg kafka.Message) error {
    start := time.Now()
    defer func() {
        p.metrics.ObserveLatency("message_processing", time.Since(start))
    }()

    var event Event
    if err := json.Unmarshal(msg.Value, &event); err != nil {
        return fmt.Errorf("unmarshaling event: %w", err)
    }

    return p.handler.HandleEvent(ctx, event)
}
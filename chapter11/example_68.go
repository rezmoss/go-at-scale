// Example 68
type EventPublisher interface {
    Publish(ctx context.Context, event Event) error
    PublishBatch(ctx context.Context, events []Event) error
}

type KafkaEventPublisher struct {
    producer  kafka.Producer
    topic     string
    metrics   MetricsRecorder
    logger    Logger
}

func (p *KafkaEventPublisher) Publish(ctx context.Context, event Event) error {
    start := time.Now()
    defer func() {
        p.metrics.ObserveLatency("event_publish", time.Since(start))
    }()

    data, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("marshaling event: %w", err)
    }

    msg := kafka.Message{
        Key:   []byte(event.AggregateID),
        Value: data,
        Headers: []kafka.Header{
            {Key: "event_type", Value: []byte(event.Type)},
            {Key: "event_version", Value: []byte(strconv.Itoa(event.Version))},
        },
    }

    if err := p.producer.Produce(ctx, p.topic, msg); err != nil {
        p.metrics.IncCounter("event_publish_errors")
        return fmt.Errorf("producing event: %w", err)
    }

    p.metrics.IncCounter("events_published")
    return nil
}
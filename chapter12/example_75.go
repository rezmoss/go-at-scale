// Example 75
type MessageBroker interface {
    Publish(ctx context.Context, topic string, msg Message) error
    Subscribe(ctx context.Context, topic string, handler MessageHandler) error
}

type KafkaBroker struct {
    producer *kafka.Writer
    consumer *kafka.Reader
    metrics  MetricsRecorder
    logger   Logger
}

func (b *KafkaBroker) Publish(ctx context.Context, topic string, msg Message) error {
    start := time.Now()
    defer func() {
        b.metrics.ObserveLatency("message_publish", time.Since(start),
            "topic", topic,
        )
    }()

    // Add tracing context
    headers := []kafka.Header{
        {Key: "trace_id", Value: []byte(trace.FromContext(ctx).SpanContext().TraceID.String())},
    }

    // Publish message
    err := b.producer.WriteMessages(ctx, kafka.Message{
        Topic:   topic,
        Key:     []byte(msg.ID),
        Value:   msg.Data,
        Headers: headers,
    })

    if err != nil {
        b.metrics.IncCounter("message_publish_errors", "topic", topic)
        return fmt.Errorf("publishing message: %w", err)
    }

    b.metrics.IncCounter("messages_published", "topic", topic)
    return nil
}
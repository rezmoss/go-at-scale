// Example 151
// internal/infrastructure/messaging/stream.go
type StreamProcessor struct {
    client      *RabbitMQClient
    handler     MessageHandler
    concurrency int
    logger      Logger
    metrics     MetricsRecorder
}

func (p *StreamProcessor) Start(ctx context.Context) error {
    msgs, err := p.client.channel.Consume(
        p.config.QueueName,
        "",    // consumer
        false, // auto-ack
        false, // exclusive
        false, // no-local
        false, // no-wait
        nil,   // args
    )
    if err != nil {
        return fmt.Errorf("starting consumer: %w", err)
    }

    for i := 0; i < p.concurrency; i++ {
        go p.processMessages(ctx, msgs)
    }

    return nil
}

func (p *StreamProcessor) processMessages(ctx context.Context, msgs <-chan amqp.Delivery) {
    for msg := range msgs {
        start := time.Now()
        err := p.handler.Handle(ctx, msg.Body)
        
        if err != nil {
            p.metrics.IncCounter("message_processing_errors")
            p.logger.Error("failed to process message", "error", err)
            msg.Nack(false, true) // requeue message
            continue
        }
        
        msg.Ack(false)
        p.metrics.ObserveLatency("message_processing", time.Since(start))
    }
}
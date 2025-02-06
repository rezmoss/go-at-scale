// Example 150
// internal/infrastructure/messaging/deadletter.go
type DeadLetterConfig struct {
    MainQueue         string
    DeadLetterQueue   string
    RetryCount        int
    RetryDelay        time.Duration
}

type DeadLetterHandler struct {
    client  *RabbitMQClient
    config  DeadLetterConfig
    logger  Logger
    metrics MetricsRecorder
}

func (h *DeadLetterHandler) Setup(ctx context.Context) error {
    // Declare dead letter exchange
    err := h.client.channel.ExchangeDeclare(
        "dlx",    // name
        "direct", // type
        true,     // durable
        false,    // auto-deleted
        false,    // internal
        false,    // no-wait
        nil,      // arguments
    )
    if err != nil {
        return fmt.Errorf("declaring dead letter exchange: %w", err)
    }

    // Declare dead letter queue
    _, err = h.client.channel.QueueDeclare(
        h.config.DeadLetterQueue,
        true,  // durable
        false, // auto-deleted
        false, // exclusive
        false, // no-wait
        amqp.Table{
            "x-dead-letter-exchange":    "",
            "x-dead-letter-routing-key": h.config.MainQueue,
            "x-message-ttl":            int32(h.config.RetryDelay.Milliseconds()),
        },
    )
    
    return err
}
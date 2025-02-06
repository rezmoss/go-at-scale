// Example 149
// internal/infrastructure/messaging/rabbitmq.go
type RabbitMQConfig struct {
    URI               string
    ReconnectDelay    time.Duration
    MaxReconnectTries int
    ExchangeName      string
    QueueName         string
    RoutingKey        string
}

type RabbitMQClient struct {
    conn    *amqp.Connection
    channel *amqp.Channel
    config  RabbitMQConfig
    logger  Logger
    metrics MetricsRecorder
}

func NewRabbitMQClient(config RabbitMQConfig, logger Logger, metrics MetricsRecorder) (*RabbitMQClient, error) {
    client := &RabbitMQClient{
        config:  config,
        logger:  logger,
        metrics: metrics,
    }
    
    if err := client.connect(); err != nil {
        return nil, fmt.Errorf("connecting to RabbitMQ: %w", err)
    }
    
    // Start connection monitoring
    go client.monitorConnection()
    
    return client, nil
}

func (c *RabbitMQClient) Publish(ctx context.Context, message []byte) error {
    start := time.Now()
    defer func() {
        c.metrics.ObserveLatency("rabbitmq_publish", time.Since(start))
    }()

    err := c.channel.PublishWithContext(ctx,
        c.config.ExchangeName,
        c.config.RoutingKey,
        false, // mandatory
        false, // immediate
        amqp.Publishing{
            ContentType:  "application/json",
            Body:        message,
            DeliveryMode: amqp.Persistent,
            Timestamp:   time.Now(),
        },
    )
    
    if err != nil {
        c.metrics.IncCounter("rabbitmq_publish_errors")
        return fmt.Errorf("publishing message: %w", err)
    }
    
    c.metrics.IncCounter("rabbitmq_messages_published")
    return nil
}
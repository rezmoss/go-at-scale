// Example 75
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

// Message represents a message to be published
type Message struct {
	ID   string
	Data []byte
}

// MessageHandler is a function that processes messages
type MessageHandler func(ctx context.Context, msg Message) error

// Logger interface for logging
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// SimpleLogger implements Logger
type SimpleLogger struct{}

func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	log.Printf("INFO: "+msg, args...)
}

func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	log.Printf("ERROR: "+msg, args...)
}

// MetricsRecorder interface for recording metrics
type MetricsRecorder interface {
	ObserveLatency(name string, duration time.Duration, labels ...string)
	IncCounter(name string, labels ...string)
}

// SimpleMetrics implements MetricsRecorder
type SimpleMetrics struct{}

func (m *SimpleMetrics) ObserveLatency(name string, duration time.Duration, labels ...string) {
	log.Printf("METRIC: %s took %v with labels %v", name, duration, labels)
}

func (m *SimpleMetrics) IncCounter(name string, labels ...string) {
	log.Printf("COUNTER: %s incremented with labels %v", name, labels)
}

// Trace context implementation
type trace struct {
	spanContext traceSpanContext
}

type traceSpanContext struct {
	TraceID uuid.UUID
}

func (t *trace) SpanContext() traceSpanContext {
	return t.spanContext
}

// FromContext extracts trace from context
func FromContext(ctx context.Context) *trace {
	val := ctx.Value("trace")
	if val == nil {
		return &trace{spanContext: traceSpanContext{TraceID: uuid.New()}}
	}
	return val.(*trace)
}

// MessageBroker interface
type MessageBroker interface {
	Publish(ctx context.Context, topic string, msg Message) error
	Subscribe(ctx context.Context, topic string, handler MessageHandler) error
}

// KafkaBroker struct
type KafkaBroker struct {
	producer *kafka.Writer
	consumer *kafka.Reader
	metrics  MetricsRecorder
	logger   Logger
}

func NewKafkaBroker(kafkaURL string) *KafkaBroker {
	producer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaBroker{
		producer: producer,
		metrics:  &SimpleMetrics{},
		logger:   &SimpleLogger{},
	}
}

// CreateTopic creates a Kafka topic if it doesn't exist
func CreateTopic(kafkaURL, topic string, partitions int) error {
	conn, err := kafka.Dial("tcp", kafkaURL)
	if err != nil {
		return fmt.Errorf("dial failed: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("get controller failed: %w", err)
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, fmt.Sprintf("%d", controller.Port)))
	if err != nil {
		return fmt.Errorf("connect to controller failed: %w", err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     partitions,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		return fmt.Errorf("create topic failed: %w", err)
	}

	return nil
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
		{Key: "trace_id", Value: []byte(FromContext(ctx).SpanContext().TraceID.String())},
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

func (b *KafkaBroker) Subscribe(ctx context.Context, topic string, handler MessageHandler) error {
	// Implementation of Subscribe method
	b.consumer = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{b.producer.Addr.String()},
		Topic:    topic,
		GroupID:  "example-consumer-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				m, err := b.consumer.ReadMessage(ctx)
				if err != nil {
					b.logger.Error("Failed to read message: %v", err)
					continue
				}

				msg := Message{
					ID:   string(m.Key),
					Data: m.Value,
				}

				traceCtx := context.WithValue(ctx, "trace", &trace{
					spanContext: traceSpanContext{
						TraceID: uuid.New(),
					},
				})

				if err := handler(traceCtx, msg); err != nil {
					b.logger.Error("Failed to process message: %v", err)
				}
			}
		}
	}()

	return nil
}

func main() {
	kafkaURL := "localhost:9092"
	topicName := "example-topic"

	// Create the topic first
	if err := CreateTopic(kafkaURL, topicName, 1); err != nil {
		log.Printf("Warning: Failed to create topic: %v", err)
		// Continue anyway as the topic might already exist or be auto-created
	}

	// Example usage
	broker := NewKafkaBroker(kafkaURL)

	// Create a context with tracing
	ctx := context.WithValue(context.Background(), "trace", &trace{
		spanContext: traceSpanContext{
			TraceID: uuid.New(),
		},
	})

	// Example message
	msg := Message{
		ID:   uuid.New().String(),
		Data: []byte("Hello Kafka"),
	}

	// Publish example
	if err := broker.Publish(ctx, topicName, msg); err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	// Subscribe example
	if err := broker.Subscribe(ctx, topicName, func(ctx context.Context, msg Message) error {
		log.Printf("Received message: %s", string(msg.Data))
		return nil
	}); err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	// Keep the application running
	select {}
}
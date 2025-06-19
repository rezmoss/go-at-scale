// Example 72
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Event represents a domain event to be processed
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// EventHandler interface for processing events
type EventHandler interface {
	HandleEvent(ctx context.Context, event Event) error
}

// Simple implementation of the EventHandler
type SimpleEventHandler struct{}

func (h *SimpleEventHandler) HandleEvent(ctx context.Context, event Event) error {
	// Simple handler that just logs the event
	fmt.Printf("Processing event: %s of type: %s\n", event.ID, event.Type)
	return nil
}

// Logger interface for logging
type Logger interface {
	Error(msg string, keysAndValues ...interface{})
}

// Simple logger implementation
type SimpleLogger struct{}

func (l *SimpleLogger) Error(msg string, keysAndValues ...interface{}) {
	log.Printf("ERROR: %s %v\n", msg, keysAndValues)
}

// MetricsRecorder interface for recording metrics
type MetricsRecorder interface {
	IncCounter(name string)
	ObserveLatency(name string, duration time.Duration)
}

// Simple metrics recorder implementation
type SimpleMetricsRecorder struct{}

func (m *SimpleMetricsRecorder) IncCounter(name string) {
	fmt.Printf("Incrementing counter: %s\n", name)
}

func (m *SimpleMetricsRecorder) ObserveLatency(name string, duration time.Duration) {
	fmt.Printf("Observed latency for %s: %v\n", name, duration)
}

// Kafka Message wrapper to match the original code
type Message struct {
	ID    string
	Value []byte
}

// Consumer interface for Kafka
type Consumer interface {
	ReadMessage(ctx context.Context) (Message, error)
}

// Kafka consumer implementation
type KafkaConsumer struct {
	consumer *kafka.Consumer
}

func NewKafkaConsumer(brokers, group, topic string) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          group,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{consumer: c}, nil
}

func (c *KafkaConsumer) ReadMessage(ctx context.Context) (Message, error) {
	// Using a small timeout to be able to check context cancellation
	msg, err := c.consumer.ReadMessage(1000 * time.Millisecond)
	if err != nil {
		if err.(kafka.Error).Code() == kafka.ErrTimedOut {
			// Check if context is done
			select {
			case <-ctx.Done():
				return Message{}, ctx.Err()
			default:
				// Just a timeout, wait longer before retrying to avoid spamming logs
				time.Sleep(1 * time.Second)
				return Message{}, fmt.Errorf("timed out reading message")
			}
		}
		return Message{}, err
	}

	// Convert to our Message type
	return Message{
		ID:    string(msg.Key),
		Value: msg.Value,
	}, nil
}

// EventProcessor from the original example
type EventProcessor struct {
	consumer   Consumer
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
				// Only increment counter and log if it's not a context cancellation
				if err != context.Canceled {
					p.metrics.IncCounter("message_read_errors")
					p.logger.Error("failed to read message", "error", err)
				}
				continue
			}

			// Successfully read a message
			fmt.Printf("Successfully read message with ID: %s\n", msg.ID)

			if err := p.processMessage(ctx, msg); err != nil {
				p.metrics.IncCounter("message_processing_errors")
				p.logger.Error("failed to process message",
					"error", err,
					"message_id", msg.ID)
				p.errorsChan <- err
			} else {
				p.metrics.IncCounter("message_processed_success")
				fmt.Printf("Successfully processed message with ID: %s\n", msg.ID)
			}
		}
	}
}

func (p *EventProcessor) processMessage(ctx context.Context, msg Message) error {
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

func main() {
	// Get Kafka configuration from environment variables or use defaults
	kafkaBrokers := getEnv("KAFKA_BROKERS", "localhost:9092")
	kafkaTopic := getEnv("KAFKA_TOPIC", "events")
	kafkaGroup := getEnv("KAFKA_GROUP", "event-processor")

	log.Printf("Starting event processor with Kafka broker: %s, topic: %s, group: %s",
		kafkaBrokers, kafkaTopic, kafkaGroup)
	log.Printf("Make sure Kafka is running, or you'll see timeout errors")

	// Create Kafka consumer
	consumer, err := NewKafkaConsumer(kafkaBrokers, kafkaGroup, kafkaTopic)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	// Create event processor
	processor := &EventProcessor{
		consumer:   consumer,
		handler:    &SimpleEventHandler{},
		metrics:    &SimpleMetricsRecorder{},
		logger:     &SimpleLogger{},
		errorsChan: make(chan error, 100), // Buffer for errors
	}

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start error handling goroutine
	go func() {
		for err := range processor.errorsChan {
			log.Printf("Error channel received: %v", err)
			// Here you could implement retry logic or dead-letter queue
		}
	}()

	// Start processor in a goroutine
	go func() {
		if err := processor.Start(ctx); err != nil && err != context.Canceled {
			log.Printf("Processor stopped with error: %v", err)
		}
	}()

	// Wait for termination signal
	sig := <-sigChan
	log.Printf("Received signal %v, shutting down...", sig)
	cancel()                    // Cancel the context to stop the processor
	time.Sleep(1 * time.Second) // Give some time for cleanup
}

// Helper function to get environment variable with default fallback
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
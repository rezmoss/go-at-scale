// Example 68
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Event represents something interesting that occurred in the system
type Event struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	AggregateID string                 `json:"aggregate_id"`
	Version     int                    `json:"version"`
	Payload     map[string]interface{} `json:"payload"`
	Metadata    map[string]string      `json:"metadata"`
	Timestamp   time.Time              `json:"timestamp"`
}

// EventVersioner interface addresses how to deal with changes to event definitions over time
type EventVersioner interface {
	UpgradeEvent(event Event) (Event, error)
	GetLatestVersion() int
}

// MetricsRecorder interface for recording metrics
type MetricsRecorder interface {
	IncCounter(name string)
	ObserveLatency(name string, duration time.Duration)
}

// Logger interface for logging
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// SimpleMetricsRecorder is a simple implementation of MetricsRecorder
type SimpleMetricsRecorder struct{}

func (m *SimpleMetricsRecorder) IncCounter(name string) {
	fmt.Printf("Counter %s incremented\n", name)
}

func (m *SimpleMetricsRecorder) ObserveLatency(name string, duration time.Duration) {
	fmt.Printf("Latency %s: %v\n", name, duration)
}

// SimpleLogger is a simple implementation of Logger
type SimpleLogger struct{}

func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	fmt.Printf("INFO: "+msg+"\n", args...)
}

func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	fmt.Printf("ERROR: "+msg+"\n", args...)
}

// EventPublisher interface
type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
	PublishBatch(ctx context.Context, events []Event) error
}

// KafkaEventPublisher implements EventPublisher using Kafka
type KafkaEventPublisher struct {
	producer *kafka.Producer // Changed from kafka.Producer to *kafka.Producer
	topic    string
	metrics  MetricsRecorder
	logger   Logger
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
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
	}

	if err := p.producer.Produce(&msg, nil); err != nil {
		p.metrics.IncCounter("event_publish_errors")
		return fmt.Errorf("producing event: %w", err)
	}

	p.metrics.IncCounter("events_published")
	return nil
}

// Implementation of PublishBatch to fulfill the EventPublisher interface
func (p *KafkaEventPublisher) PublishBatch(ctx context.Context, events []Event) error {
	for _, event := range events {
		if err := p.Publish(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// NewKafkaEventPublisher creates a new KafkaEventPublisher
func NewKafkaEventPublisher(producer *kafka.Producer, topic string, metrics MetricsRecorder, logger Logger) *KafkaEventPublisher {
	return &KafkaEventPublisher{
		producer: producer,
		topic:    topic,
		metrics:  metrics,
		logger:   logger,
	}
}

func main() {
	// Create a Kafka producer configuration
	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	}

	// Create a producer instance
	producer, err := kafka.NewProducer(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create producer: %s", err))
	}
	defer producer.Close()

	// Create the publisher with required dependencies
	publisher := NewKafkaEventPublisher(
		producer,
		"events",
		&SimpleMetricsRecorder{},
		&SimpleLogger{},
	)

	// Create and publish a test event
	event := Event{
		ID:          "event-123",
		Type:        "order.created",
		AggregateID: "order-456",
		Version:     1,
		Payload: map[string]interface{}{
			"orderID": "order-456",
			"amount":  99.99,
		},
		Metadata: map[string]string{
			"source": "test",
		},
		Timestamp: time.Now(),
	}

	ctx := context.Background()
	err = publisher.Publish(ctx, event)
	if err != nil {
		fmt.Printf("Failed to publish event: %v\n", err)
	} else {
		fmt.Println("Event published successfully!")
	}
}
// Example 70
package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Command and Query interfaces
type Command interface {
	ToEvents() ([]Event, error)
}

type Query interface {
	// Query interface marker
}

// Handler interfaces
type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) error
}

type QueryHandler interface {
	Handle(ctx context.Context, query Query) (interface{}, error)
}

// Event related types
type Event struct {
	ID        string
	Type      string
	Data      interface{}
	Timestamp time.Time
}

type EventStore interface {
	Save(ctx context.Context, events []Event) error
}

type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
}

// Support types
type Validator interface {
	Validate(cmd Command) error
}

type MetricsRecorder interface {
	ObserveLatency(metric string, duration time.Duration)
	IncCounter(metric string)
}

type Logger interface {
	Error(msg string, args ...interface{})
}

// Implementation of UserCommandHandler
type UserCommandHandler struct {
	eventStore EventStore
	publisher  EventPublisher
	validator  Validator
	metrics    MetricsRecorder
	logger     Logger
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

// Simple implementations for demonstration purposes
type SimpleValidator struct{}

func (v *SimpleValidator) Validate(cmd Command) error {
	// Basic validation logic would go here
	return nil
}

type SimpleEventStore struct{}

func (s *SimpleEventStore) Save(ctx context.Context, events []Event) error {
	for _, event := range events {
		fmt.Printf("Saving event: %s\n", event.ID)
	}
	return nil
}

type SimpleEventPublisher struct{}

func (p *SimpleEventPublisher) Publish(ctx context.Context, event Event) error {
	fmt.Printf("Publishing event: %s\n", event.ID)
	return nil
}

type SimpleMetricsRecorder struct{}

func (m *SimpleMetricsRecorder) ObserveLatency(metric string, duration time.Duration) {
	fmt.Printf("Metric %s: %v\n", metric, duration)
}

func (m *SimpleMetricsRecorder) IncCounter(metric string) {
	fmt.Printf("Increasing counter: %s\n", metric)
}

type SimpleLogger struct{}

func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	log.Printf("ERROR: "+msg, args...)
}

// Example command
type CreateUserCommand struct {
	UserID    string
	Email     string
	FirstName string
	LastName  string
}

func (c *CreateUserCommand) ToEvents() ([]Event, error) {
	return []Event{
		{
			ID:        "evt-" + c.UserID,
			Type:      "UserCreated",
			Data:      c,
			Timestamp: time.Now(),
		},
	}, nil
}

func main() {
	// Set up the command handler with dependencies
	handler := &UserCommandHandler{
		eventStore: &SimpleEventStore{},
		publisher:  &SimpleEventPublisher{},
		validator:  &SimpleValidator{},
		metrics:    &SimpleMetricsRecorder{},
		logger:     &SimpleLogger{},
	}

	// Create and handle a command
	cmd := &CreateUserCommand{
		UserID:    "user-123",
		Email:     "example@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	ctx := context.Background()
	if err := handler.Handle(ctx, cmd); err != nil {
		log.Fatalf("Error handling command: %v", err)
	}

	fmt.Println("Command processed successfully!")
}
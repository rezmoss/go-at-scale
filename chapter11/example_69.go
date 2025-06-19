// Example 69
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// Event represents a domain event
type Event struct {
	ID          string          `json:"id"`
	AggregateID string          `json:"aggregate_id"`
	Type        string          `json:"type"`
	Version     int             `json:"version"`
	Data        json.RawMessage `json:"data"`
	CreatedAt   time.Time       `json:"created_at"`
}

// Aggregate interface defines methods that an aggregate should implement
type Aggregate interface {
	ID() string
	Version() int
	ApplyEvent(event Event) error
}

// EventStore interface defines methods for storing and loading events
type EventStore interface {
	Save(ctx context.Context, events []Event) error
	Load(ctx context.Context, aggregateID string) ([]Event, error)
}

// Logger is a simple interface for logging
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// MetricsRecorder is a simple interface for recording metrics
type MetricsRecorder interface {
	RecordEventSaved(eventType string)
	RecordEventLoaded(eventType string)
}

// SimpleLogger implements the Logger interface
type SimpleLogger struct{}

func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	log.Printf("[INFO] "+msg, args...)
}

func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	log.Printf("[ERROR] "+msg, args...)
}

// SimpleMetricsRecorder implements the MetricsRecorder interface
type SimpleMetricsRecorder struct{}

func (m *SimpleMetricsRecorder) RecordEventSaved(eventType string) {
	log.Printf("Metric: Event saved: %s", eventType)
}

func (m *SimpleMetricsRecorder) RecordEventLoaded(eventType string) {
	log.Printf("Metric: Event loaded: %s", eventType)
}

// PostgresEventStore implements the EventStore interface
type PostgresEventStore struct {
	db      *sql.DB
	metrics MetricsRecorder
	logger  Logger
}

// NewPostgresEventStore creates a new PostgresEventStore
func NewPostgresEventStore(db *sql.DB, metrics MetricsRecorder, logger Logger) *PostgresEventStore {
	return &PostgresEventStore{
		db:      db,
		metrics: metrics,
		logger:  logger,
	}
}

// saveEvent saves a single event to the database
func (s *PostgresEventStore) saveEvent(ctx context.Context, tx *sql.Tx, event Event) error {
	query := `
		INSERT INTO events (id, aggregate_id, type, version, data, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := tx.ExecContext(ctx, query, event.ID, event.AggregateID, event.Type, event.Version, event.Data, event.CreatedAt)
	if err != nil {
		return fmt.Errorf("inserting event: %w", err)
	}

	s.metrics.RecordEventSaved(event.Type)
	s.logger.Info("Event saved", "type", event.Type, "aggregateID", event.AggregateID)
	return nil
}

// Save stores multiple events in a single transaction
func (s *PostgresEventStore) Save(ctx context.Context, events []Event) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	for _, event := range events {
		if err := s.saveEvent(ctx, tx, event); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

// Load retrieves all events for an aggregate
func (s *PostgresEventStore) Load(ctx context.Context, aggregateID string) ([]Event, error) {
	query := `
		SELECT id, aggregate_id, type, version, data, created_at
		FROM events
		WHERE aggregate_id = $1
		ORDER BY version ASC
	`
	rows, err := s.db.QueryContext(ctx, query, aggregateID)
	if err != nil {
		return nil, fmt.Errorf("querying events: %w", err)
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		if err := rows.Scan(&event.ID, &event.AggregateID, &event.Type, &event.Version, &event.Data, &event.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning event: %w", err)
		}
		events = append(events, event)
		s.metrics.RecordEventLoaded(event.Type)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating over events: %w", err)
	}

	s.logger.Info("Events loaded", "aggregateID", aggregateID, "count", len(events))
	return events, nil
}

// Example Order aggregate
type Order struct {
	OrderID      string    `json:"order_id"`
	CustomerID   string    `json:"customer_id"`
	Status       string    `json:"status"`
	TotalAmount  float64   `json:"total_amount"`
	CreatedAt    time.Time `json:"created_at"`
	LastModified time.Time `json:"last_modified"`
	version      int
}

// ID returns the aggregate ID
func (o *Order) ID() string {
	return o.OrderID
}

// Version returns the aggregate version
func (o *Order) Version() int {
	return o.version
}

// ApplyEvent applies an event to the order aggregate
func (o *Order) ApplyEvent(event Event) error {
	switch event.Type {
	case "OrderCreated":
		var data struct {
			CustomerID  string    `json:"customer_id"`
			TotalAmount float64   `json:"total_amount"`
			CreatedAt   time.Time `json:"created_at"`
		}
		if err := json.Unmarshal(event.Data, &data); err != nil {
			return err
		}
		o.OrderID = event.AggregateID
		o.CustomerID = data.CustomerID
		o.TotalAmount = data.TotalAmount
		o.Status = "Created"
		o.CreatedAt = data.CreatedAt
		o.LastModified = data.CreatedAt
	case "OrderUpdated":
		var data struct {
			TotalAmount  float64   `json:"total_amount"`
			LastModified time.Time `json:"last_modified"`
		}
		if err := json.Unmarshal(event.Data, &data); err != nil {
			return err
		}
		o.TotalAmount = data.TotalAmount
		o.LastModified = data.LastModified
	case "OrderCanceled":
		var data struct {
			LastModified time.Time `json:"last_modified"`
		}
		if err := json.Unmarshal(event.Data, &data); err != nil {
			return err
		}
		o.Status = "Canceled"
		o.LastModified = data.LastModified
	default:
		return fmt.Errorf("unknown event type: %s", event.Type)
	}
	o.version = event.Version
	return nil
}

// LoadOrder loads an order from its event stream
func LoadOrder(ctx context.Context, orderID string, store EventStore) (*Order, error) {
	events, err := store.Load(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("loading events: %w", err)
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("order not found: %s", orderID)
	}

	order := &Order{}
	for _, event := range events {
		if err := order.ApplyEvent(event); err != nil {
			return nil, fmt.Errorf("applying event: %w", err)
		}
	}

	return order, nil
}

// CreateOrder creates a new order and saves the event
func CreateOrder(ctx context.Context, store EventStore, orderID, customerID string, totalAmount float64) (*Order, error) {
	now := time.Now()
	data, err := json.Marshal(map[string]interface{}{
		"customer_id":  customerID,
		"total_amount": totalAmount,
		"created_at":   now,
	})
	if err != nil {
		return nil, fmt.Errorf("marshaling event data: %w", err)
	}

	event := Event{
		ID:          fmt.Sprintf("evt-%d", time.Now().UnixNano()),
		AggregateID: orderID,
		Type:        "OrderCreated",
		Version:     1,
		Data:        data,
		CreatedAt:   now,
	}

	if err := store.Save(ctx, []Event{event}); err != nil {
		return nil, fmt.Errorf("saving event: %w", err)
	}

	order := &Order{}
	if err := order.ApplyEvent(event); err != nil {
		return nil, fmt.Errorf("applying event: %w", err)
	}

	return order, nil
}

// UpdateOrder updates an order and saves the event
func UpdateOrder(ctx context.Context, store EventStore, order *Order, totalAmount float64) error {
	now := time.Now()
	data, err := json.Marshal(map[string]interface{}{
		"total_amount":  totalAmount,
		"last_modified": now,
	})
	if err != nil {
		return fmt.Errorf("marshaling event data: %w", err)
	}

	event := Event{
		ID:          fmt.Sprintf("evt-%d", time.Now().UnixNano()),
		AggregateID: order.ID(),
		Type:        "OrderUpdated",
		Version:     order.Version() + 1,
		Data:        data,
		CreatedAt:   now,
	}

	if err := store.Save(ctx, []Event{event}); err != nil {
		return fmt.Errorf("saving event: %w", err)
	}

	if err := order.ApplyEvent(event); err != nil {
		return fmt.Errorf("applying event: %w", err)
	}

	return nil
}

// CancelOrder cancels an order and saves the event
func CancelOrder(ctx context.Context, store EventStore, order *Order) error {
	now := time.Now()
	data, err := json.Marshal(map[string]interface{}{
		"last_modified": now,
	})
	if err != nil {
		return fmt.Errorf("marshaling event data: %w", err)
	}

	event := Event{
		ID:          fmt.Sprintf("evt-%d", time.Now().UnixNano()),
		AggregateID: order.ID(),
		Type:        "OrderCanceled",
		Version:     order.Version() + 1,
		Data:        data,
		CreatedAt:   now,
	}

	if err := store.Save(ctx, []Event{event}); err != nil {
		return fmt.Errorf("saving event: %w", err)
	}

	if err := order.ApplyEvent(event); err != nil {
		return fmt.Errorf("applying event: %w", err)
	}

	return nil
}

func setupDatabase(db *sql.DB) error {
	// Create events table if not exists
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS events (
			id TEXT PRIMARY KEY,
			aggregate_id TEXT NOT NULL,
			type TEXT NOT NULL,
			version INTEGER NOT NULL,
			data JSONB NOT NULL,
			created_at TIMESTAMP NOT NULL,
			UNIQUE(aggregate_id, version)
		)
	`)
	return err
}

func main() {
	// Connect to PostgreSQL
	db, err := sql.Open("postgres", "postgresql://pg:pg@localhost:5432/eventsourcing?sslmode=disable")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Setup the database
	if err := setupDatabase(db); err != nil {
		log.Fatalf("Error setting up database: %v", err)
	}

	// Create event store
	logger := &SimpleLogger{}
	metrics := &SimpleMetricsRecorder{}
	store := NewPostgresEventStore(db, metrics, logger)

	// Create a new order
	ctx := context.Background()
	orderID := fmt.Sprintf("order-%d", time.Now().UnixNano())
	order, err := CreateOrder(ctx, store, orderID, "customer-123", 100.50)
	if err != nil {
		log.Fatalf("Error creating order: %v", err)
	}
	log.Printf("Order created: %+v", order)

	// Update the order
	if err := UpdateOrder(ctx, store, order, 150.75); err != nil {
		log.Fatalf("Error updating order: %v", err)
	}
	log.Printf("Order updated: %+v", order)

	// Load the order from events
	loadedOrder, err := LoadOrder(ctx, orderID, store)
	if err != nil {
		log.Fatalf("Error loading order: %v", err)
	}
	log.Printf("Order loaded: %+v", loadedOrder)

	// Cancel the order
	if err := CancelOrder(ctx, store, order); err != nil {
		log.Fatalf("Error canceling order: %v", err)
	}
	log.Printf("Order canceled: %+v", order)

	// Load the order again to see the final state
	finalOrder, err := LoadOrder(ctx, orderID, store)
	if err != nil {
		log.Fatalf("Error loading final order: %v", err)
	}
	log.Printf("Final order state: %+v", finalOrder)
}
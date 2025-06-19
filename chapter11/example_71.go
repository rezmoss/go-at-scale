// Example 71
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
	ID        string
	Type      string
	Data      json.RawMessage
	Timestamp time.Time
}

// Projector interface for projecting events
type Projector interface {
	Project(ctx context.Context, event Event) error
	Rebuild(ctx context.Context) error
}

// EventStore interface for loading events
type EventStore interface {
	LoadAll(ctx context.Context) ([]Event, error)
	Append(ctx context.Context, event Event) error
}

// MetricsRecorder interface for recording metrics
type MetricsRecorder interface {
	RecordProjection(eventType string, duration time.Duration)
}

// Logger interface for logging
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// PostgresEventStore implements EventStore using PostgreSQL
type PostgresEventStore struct {
	db *sql.DB
}

func NewPostgresEventStore(db *sql.DB) *PostgresEventStore {
	return &PostgresEventStore{db: db}
}

func (s *PostgresEventStore) LoadAll(ctx context.Context) ([]Event, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, type, data, timestamp FROM events ORDER BY timestamp ASC")
	if err != nil {
		return nil, fmt.Errorf("querying events: %w", err)
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.Type, &e.Data, &e.Timestamp); err != nil {
			return nil, fmt.Errorf("scanning event: %w", err)
		}
		events = append(events, e)
	}

	return events, nil
}

func (s *PostgresEventStore) Append(ctx context.Context, event Event) error {
	//Important: In a real-world application, a secure migration tool would parse and execute these operations, below is a simplified example and should not be used in production and is not secure
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO events (id, type, data, timestamp) VALUES ($1, $2, $3, $4)",
		event.ID, event.Type, event.Data, event.Timestamp)
	if err != nil {
		return fmt.Errorf("inserting event: %w", err)
	}
	return nil
}

// SimpleLogger implements Logger
type SimpleLogger struct{}

func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	log.Printf("INFO: "+msg, args...)
}

func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	log.Printf("ERROR: "+msg, args...)
}

// SimpleMetrics implements MetricsRecorder
type SimpleMetrics struct{}

func (m *SimpleMetrics) RecordProjection(eventType string, duration time.Duration) {
	log.Printf("METRIC: Projected %s in %v", eventType, duration)
}

// UserProjector implements the Projector interface for User events
type UserProjector struct {
	db         *sql.DB
	eventStore EventStore
	metrics    MetricsRecorder
	logger     Logger
}

func NewUserProjector(db *sql.DB, eventStore EventStore, metrics MetricsRecorder, logger Logger) *UserProjector {
	return &UserProjector{
		db:         db,
		eventStore: eventStore,
		metrics:    metrics,
		logger:     logger,
	}
}

func (p *UserProjector) Project(ctx context.Context, event Event) error {
	start := time.Now()
	defer func() {
		p.metrics.RecordProjection(event.Type, time.Since(start))
	}()

	switch event.Type {
	case "UserCreated":
		return p.handleUserCreated(ctx, event)
	case "UserUpdated":
		return p.handleUserUpdated(ctx, event)
	case "UserDeleted":
		return p.handleUserDeleted(ctx, event)
	default:
		return fmt.Errorf("unknown event type: %s", event.Type)
	}
}

func (p *UserProjector) Rebuild(ctx context.Context) error {
	// Start transaction
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	// Clear existing projections
	if err := p.clearProjections(ctx, tx); err != nil {
		return err
	}

	// Load all events
	events, err := p.eventStore.LoadAll(ctx)
	if err != nil {
		return fmt.Errorf("loading events: %w", err)
	}

	// Project all events
	for _, event := range events {
		if err := p.projectWithTx(ctx, tx, event); err != nil {
			return err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

// Helper method to clear projections
func (p *UserProjector) clearProjections(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "DELETE FROM user_projection")
	if err != nil {
		return fmt.Errorf("clearing user projection: %w", err)
	}
	return nil
}

// Project with transaction
func (p *UserProjector) projectWithTx(ctx context.Context, tx *sql.Tx, event Event) error {
	switch event.Type {
	case "UserCreated":
		return p.handleUserCreatedTx(ctx, tx, event)
	case "UserUpdated":
		return p.handleUserUpdatedTx(ctx, tx, event)
	case "UserDeleted":
		return p.handleUserDeletedTx(ctx, tx, event)
	default:
		return nil // Skip unknown events during rebuild
	}
}

// Event handlers
func (p *UserProjector) handleUserCreated(ctx context.Context, event Event) error {
	var userData struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.Unmarshal(event.Data, &userData); err != nil {
		return fmt.Errorf("unmarshaling user data: %w", err)
	}
	//Important: In a real-world application, a secure migration tool would parse and execute these operations, below is a simplified example and should not be used in production and is not secure
	_, err := p.db.ExecContext(ctx,
		"INSERT INTO user_projection (id, name, email) VALUES ($1, $2, $3)",
		userData.ID, userData.Name, userData.Email)
	if err != nil {
		return fmt.Errorf("inserting user: %w", err)
	}

	p.logger.Info("Projected UserCreated for user %s", userData.ID)
	return nil
}

func (p *UserProjector) handleUserUpdated(ctx context.Context, event Event) error {
	var userData struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.Unmarshal(event.Data, &userData); err != nil {
		return fmt.Errorf("unmarshaling user data: %w", err)
	}
	//Important: In a real-world application, a secure migration tool would parse and execute these operations, below is a simplified example and should not be used in production and is not secure
	_, err := p.db.ExecContext(ctx,
		"UPDATE user_projection SET name = $2, email = $3 WHERE id = $1",
		userData.ID, userData.Name, userData.Email)
	if err != nil {
		return fmt.Errorf("updating user: %w", err)
	}

	p.logger.Info("Projected UserUpdated for user %s", userData.ID)
	return nil
}

func (p *UserProjector) handleUserDeleted(ctx context.Context, event Event) error {
	var userData struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(event.Data, &userData); err != nil {
		return fmt.Errorf("unmarshaling user data: %w", err)
	}
	//Important: In a real-world application, a secure migration tool would parse and execute these operations, below is a simplified example and should not be used in production and is not secure
	_, err := p.db.ExecContext(ctx,
		"DELETE FROM user_projection WHERE id = $1",
		userData.ID)
	if err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}

	p.logger.Info("Projected UserDeleted for user %s", userData.ID)
	return nil
}

// Transaction-based handlers for rebuild
func (p *UserProjector) handleUserCreatedTx(ctx context.Context, tx *sql.Tx, event Event) error {
	var userData struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.Unmarshal(event.Data, &userData); err != nil {
		return fmt.Errorf("unmarshaling user data: %w", err)
	}

	_, err := tx.ExecContext(ctx,
		"INSERT INTO user_projection (id, name, email) VALUES ($1, $2, $3)",
		userData.ID, userData.Name, userData.Email)
	if err != nil {
		return fmt.Errorf("inserting user: %w", err)
	}

	return nil
}

func (p *UserProjector) handleUserUpdatedTx(ctx context.Context, tx *sql.Tx, event Event) error {
	var userData struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.Unmarshal(event.Data, &userData); err != nil {
		return fmt.Errorf("unmarshaling user data: %w", err)
	}

	_, err := tx.ExecContext(ctx,
		"UPDATE user_projection SET name = $2, email = $3 WHERE id = $1",
		userData.ID, userData.Name, userData.Email)
	if err != nil {
		return fmt.Errorf("updating user: %w", err)
	}

	return nil
}

func (p *UserProjector) handleUserDeletedTx(ctx context.Context, tx *sql.Tx, event Event) error {
	var userData struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(event.Data, &userData); err != nil {
		return fmt.Errorf("unmarshaling user data: %w", err)
	}

	_, err := tx.ExecContext(ctx,
		"DELETE FROM user_projection WHERE id = $1",
		userData.ID)
	if err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}

	return nil
}

func main() {
	ctx := context.Background()

	// Connect to database
	db, err := sql.Open("postgres", "postgres://pg:pg@localhost:5432/eventprojection?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create dependencies
	eventStore := NewPostgresEventStore(db)
	logger := &SimpleLogger{}
	metrics := &SimpleMetrics{}

	// Create projector
	projector := NewUserProjector(db, eventStore, metrics, logger)

	// Example: rebuild projection
	if err := projector.Rebuild(ctx); err != nil {
		log.Fatalf("Failed to rebuild projection: %v", err)
	}
	logger.Info("Projection rebuilt successfully")

	// Example: handle a new event
	newEvent := Event{
		ID:        "evt-123",
		Type:      "UserCreated",
		Data:      json.RawMessage(`{"id":"user-123","name":"John Doe","email":"john@example.com"}`),
		Timestamp: time.Now(),
	}

	if err := eventStore.Append(ctx, newEvent); err != nil {
		log.Fatalf("Failed to append event: %v", err)
	}

	if err := projector.Project(ctx, newEvent); err != nil {
		log.Fatalf("Failed to project event: %v", err)
	}

	logger.Info("New event projected successfully")
}
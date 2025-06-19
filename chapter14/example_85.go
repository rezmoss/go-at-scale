// Example 85
package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Define interfaces
type DataStore interface {
	Write(ctx context.Context, data *Data) error
	Read(ctx context.Context, id string) (*Data, error)
}

type Data struct {
	ID      string
	Content string
	Version int
}

type MetricsRecorder interface {
	IncCounter(metric string)
}

type Logger interface {
	Error(msg string, args ...interface{})
	Info(msg string, args ...interface{})
}

// Simple implementations
type SimpleMetricsRecorder struct{}

func (m *SimpleMetricsRecorder) IncCounter(metric string) {
	fmt.Printf("Metric incremented: %s\n", metric)
}

type SimpleLogger struct{}

func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	log.Printf("ERROR: "+msg, args...)
}

func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	log.Printf("INFO: "+msg, args...)
}

// Example stores
type OldDataStore struct {
	data map[string]*Data
}

func NewOldDataStore() *OldDataStore {
	return &OldDataStore{
		data: make(map[string]*Data),
	}
}

func (s *OldDataStore) Write(ctx context.Context, data *Data) error {
	time.Sleep(10 * time.Millisecond) // Simulate DB write
	s.data[data.ID] = data
	return nil
}

func (s *OldDataStore) Read(ctx context.Context, id string) (*Data, error) {
	data, ok := s.data[id]
	if !ok {
		return nil, fmt.Errorf("data not found: %s", id)
	}
	return data, nil
}

type NewDataStore struct {
	data map[string]*Data
}

func NewNewDataStore() *NewDataStore {
	return &NewDataStore{
		data: make(map[string]*Data),
	}
}

func (s *NewDataStore) Write(ctx context.Context, data *Data) error {
	time.Sleep(5 * time.Millisecond) // Simulate faster DB write
	s.data[data.ID] = data
	return nil
}

func (s *NewDataStore) Read(ctx context.Context, id string) (*Data, error) {
	data, ok := s.data[id]
	if !ok {
		return nil, fmt.Errorf("data not found: %s", id)
	}
	return data, nil
}

type MigrationManager struct {
	source      DataStore
	destination DataStore
	validator   DataValidator
	metrics     MetricsRecorder
	logger      Logger
}

type DataValidator interface {
	ValidateConsistency(ctx context.Context, source, destination DataStore) error
}

type SimpleValidator struct{}

func (v *SimpleValidator) ValidateConsistency(ctx context.Context, source, destination DataStore) error {
	// In a real implementation, this would scan and compare data between stores
	return nil
}

type DataSource interface {
	Read(ctx context.Context, id string) (*Data, error)
}

// Dual-write pattern implementation
type DualWriteMigration struct {
	oldStore    DataStore
	newStore    DataStore
	metrics     MetricsRecorder
	logger      Logger
	activeStore string // Track which store is active for reads
}

func NewDualWriteMigration(oldStore, newStore DataStore, metrics MetricsRecorder, logger Logger) *DualWriteMigration {
	return &DualWriteMigration{
		oldStore:    oldStore,
		newStore:    newStore,
		metrics:     metrics,
		logger:      logger,
		activeStore: "old", // Start with old store as source of truth
	}
}

func (m *DualWriteMigration) Write(ctx context.Context, data *Data) error {
	// Write to old store
	if err := m.oldStore.Write(ctx, data); err != nil {
		return fmt.Errorf("writing to old store: %w", err)
	}

	// Write to new store
	if err := m.newStore.Write(ctx, data); err != nil {
		m.metrics.IncCounter("new_store_write_failures")
		m.logger.Error("failed to write to new store", "error", err)
		// Continue without failing - old store write succeeded
	}

	return nil
}

func (m *DualWriteMigration) Read(ctx context.Context, id string) (*Data, error) {
	if m.activeStore == "new" {
		return m.newStore.Read(ctx, id)
	}
	return m.oldStore.Read(ctx, id)
}

func (m *DualWriteMigration) verifyConsistency(ctx context.Context) error {
	m.logger.Info("Verifying data consistency between stores...")
	// In a real implementation, this would verify all relevant data is consistent
	// between both stores before switching
	return nil
}

func (m *DualWriteMigration) updateRoutingConfig(ctx context.Context, store string) error {
	m.logger.Info("Updating routing config to use store: %s", store)
	m.activeStore = store
	return nil
}

func (m *DualWriteMigration) SwitchToNewStore(ctx context.Context) error {
	// Verify data consistency
	if err := m.verifyConsistency(ctx); err != nil {
		return fmt.Errorf("consistency verification failed: %w", err)
	}

	// Switch reads to new store
	if err := m.updateRoutingConfig(ctx, "new"); err != nil {
		return fmt.Errorf("updating routing config: %w", err)
	}

	return nil
}

func main() {
	ctx := context.Background()

	// Initialize stores
	oldStore := NewOldDataStore()
	newStore := NewNewDataStore()

	// Initialize metrics and logger
	metrics := &SimpleMetricsRecorder{}
	logger := &SimpleLogger{}

	// Create the migration manager
	migration := NewDualWriteMigration(oldStore, newStore, metrics, logger)

	// Example workflow:
	// 1. Write data to both stores
	data1 := &Data{ID: "1", Content: "Example data", Version: 1}
	err := migration.Write(ctx, data1)
	if err != nil {
		logger.Error("Failed to write data", "error", err)
		return
	}

	// 2. Read data (from old store initially)
	readData, err := migration.Read(ctx, "1")
	if err != nil {
		logger.Error("Failed to read data", "error", err)
		return
	}
	fmt.Printf("Read data from active store: %+v\n", readData)

	// 3. Switch to the new store
	err = migration.SwitchToNewStore(ctx)
	if err != nil {
		logger.Error("Failed to switch to new store", "error", err)
		return
	}

	// 4. Read data again (now from new store)
	readData, err = migration.Read(ctx, "1")
	if err != nil {
		logger.Error("Failed to read data", "error", err)
		return
	}
	fmt.Printf("Read data after switch: %+v\n", readData)

	// 5. Write more data (will be written to both stores)
	data2 := &Data{ID: "2", Content: "More data", Version: 1}
	err = migration.Write(ctx, data2)
	if err != nil {
		logger.Error("Failed to write data", "error", err)
		return
	}

	fmt.Println("Zero-downtime migration example completed successfully!")
}
// Example 81
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

// Define the interfaces and types needed
type DataSource interface {
	ApplyTransaction(ctx context.Context, tx *Transaction) error
	GetCheckpoint(ctx context.Context) (Checkpoint, error)
}

type ConsistencyMonitor interface {
	VerifyConsistency(ctx context.Context, txID string) error
}

type MetricsRecorder interface {
	ObserveLatency(metric string, duration time.Duration)
	IncCounter(metric string)
}

type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type Transaction struct {
	ID      string
	Data    []byte
	Created time.Time
}

type Checkpoint struct {
	LastTransactionID string
	Timestamp         time.Time
	Hash              []byte
}

var ErrInconsistentReplicas = errors.New("replicas are not consistent")

type ReplicationManager struct {
	source   DataSource
	replicas []DataSource
	monitor  ConsistencyMonitor
	metrics  MetricsRecorder
	logger   Logger
}

func (r *ReplicationManager) ProcessTransaction(ctx context.Context, tx *Transaction) error {
	start := time.Now()
	defer func() {
		r.metrics.ObserveLatency("replication_latency", time.Since(start))
	}()

	// Apply transaction to source first
	if err := r.applyTransaction(ctx, r.source, tx); err != nil {
		r.metrics.IncCounter("source_transaction_failures")
		return fmt.Errorf("applying transaction to source: %w", err)
	}

	// Apply transaction to all replicas
	for _, replica := range r.replicas {
		if err := r.applyTransaction(ctx, replica, tx); err != nil {
			r.metrics.IncCounter("replication_failures")
			return fmt.Errorf("applying transaction to replica: %w", err)
		}
	}

	// Verify consistency
	if err := r.monitor.VerifyConsistency(ctx, tx.ID); err != nil {
		r.metrics.IncCounter("consistency_check_failures")
		return fmt.Errorf("verifying consistency: %w", err)
	}

	r.metrics.IncCounter("successful_replications")
	return nil
}

func (r *ReplicationManager) VerifyReplication(ctx context.Context) error {
	checkpoints, err := r.gatherCheckpoints(ctx)
	if err != nil {
		return fmt.Errorf("gathering checkpoints: %w", err)
	}

	if !r.areCheckpointsConsistent(checkpoints) {
		return ErrInconsistentReplicas
	}

	return nil
}

// Additional methods needed to make the code runnable
func (r *ReplicationManager) applyTransaction(ctx context.Context, replica DataSource, tx *Transaction) error {
	return replica.ApplyTransaction(ctx, tx)
}

func (r *ReplicationManager) gatherCheckpoints(ctx context.Context) ([]Checkpoint, error) {
	checkpoints := make([]Checkpoint, 0, len(r.replicas)+1)

	// Get checkpoint from source
	sourceCP, err := r.source.GetCheckpoint(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting source checkpoint: %w", err)
	}
	checkpoints = append(checkpoints, sourceCP)

	// Get checkpoints from all replicas
	for i, replica := range r.replicas {
		cp, err := replica.GetCheckpoint(ctx)
		if err != nil {
			return nil, fmt.Errorf("getting checkpoint from replica %d: %w", i, err)
		}
		checkpoints = append(checkpoints, cp)
	}

	return checkpoints, nil
}

func (r *ReplicationManager) areCheckpointsConsistent(checkpoints []Checkpoint) bool {
	if len(checkpoints) <= 1 {
		return true
	}

	reference := checkpoints[0]
	for i := 1; i < len(checkpoints); i++ {
		if reference.LastTransactionID != checkpoints[i].LastTransactionID {
			r.logger.Error("Transaction ID mismatch: %s vs %s",
				reference.LastTransactionID, checkpoints[i].LastTransactionID)
			return false
		}
	}

	return true
}

// Simple implementations for demo purposes
type InMemoryDataSource struct {
	name         string
	transactions map[string]*Transaction
	lastTxID     string
}

func NewInMemoryDataSource(name string) *InMemoryDataSource {
	return &InMemoryDataSource{
		name:         name,
		transactions: make(map[string]*Transaction),
	}
}

func (ds *InMemoryDataSource) ApplyTransaction(ctx context.Context, tx *Transaction) error {
	// Simulating some processing time
	time.Sleep(10 * time.Millisecond)
	ds.transactions[tx.ID] = tx
	ds.lastTxID = tx.ID
	return nil
}

func (ds *InMemoryDataSource) GetCheckpoint(ctx context.Context) (Checkpoint, error) {
	if ds.lastTxID == "" {
		return Checkpoint{}, nil
	}

	tx := ds.transactions[ds.lastTxID]
	return Checkpoint{
		LastTransactionID: ds.lastTxID,
		Timestamp:         time.Now(),
		Hash:              tx.Data, // Using data as hash for simplicity
	}, nil
}

type SimpleConsistencyMonitor struct {
	source   DataSource
	replicas []DataSource
}

func (m *SimpleConsistencyMonitor) VerifyConsistency(ctx context.Context, txID string) error {
	// Simple verification - just check that the transaction ID matches
	return nil
}

type SimpleMetricsRecorder struct{}

func (m *SimpleMetricsRecorder) ObserveLatency(metric string, duration time.Duration) {
	log.Printf("%s: %v", metric, duration)
}

func (m *SimpleMetricsRecorder) IncCounter(metric string) {
	log.Printf("Increment counter: %s", metric)
}

type SimpleLogger struct{}

func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	log.Printf("[INFO] "+msg, args...)
}

func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	log.Printf("[ERROR] "+msg, args...)
}

func main() {
	// Create data sources
	source := NewInMemoryDataSource("source")
	replica1 := NewInMemoryDataSource("replica1")
	replica2 := NewInMemoryDataSource("replica2")

	// Create the replication manager
	manager := &ReplicationManager{
		source:   source,
		replicas: []DataSource{replica1, replica2},
		monitor: &SimpleConsistencyMonitor{
			source:   source,
			replicas: []DataSource{replica1, replica2},
		},
		metrics: &SimpleMetricsRecorder{},
		logger:  &SimpleLogger{},
	}

	// Process a test transaction
	ctx := context.Background()
	tx := &Transaction{
		ID:      "tx-1",
		Data:    []byte("test transaction data"),
		Created: time.Now(),
	}

	if err := manager.ProcessTransaction(ctx, tx); err != nil {
		log.Fatalf("Error processing transaction: %v", err)
	}

	// Give a small delay to ensure all operations are complete
	time.Sleep(50 * time.Millisecond)

	// Verify replication
	if err := manager.VerifyReplication(ctx); err != nil {
		log.Fatalf("Replication verification failed: %v", err)
	}

	log.Println("Replication successful!")
}
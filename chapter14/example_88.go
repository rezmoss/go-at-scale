// Example 88
package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

// MigrationStatus represents the current status of the migration
type MigrationStatus string

const (
	StatusPending   MigrationStatus = "PENDING"
	StatusRunning   MigrationStatus = "RUNNING"
	StatusCompleted MigrationStatus = "COMPLETED"
	StatusFailed    MigrationStatus = "FAILED"
	StatusAborted   MigrationStatus = "ABORTED"
)

// Storage interface for persisting migration progress
type Storage interface {
	SaveProgress(ctx context.Context, progress MigrationProgress) error
}

// MetricsRecorder interface for recording metrics
type MetricsRecorder interface {
	SetGauge(name string, value float64)
}

// Logger interface for logging
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// MigrationTracker handles tracking and reporting migration progress
type MigrationTracker struct {
	storage  Storage
	metrics  MetricsRecorder
	logger   Logger
	progress atomic.Value
}

// MigrationProgress holds the current state of the migration
type MigrationProgress struct {
	TotalRecords   int64
	ProcessedCount int64
	FailedCount    int64
	StartTime      time.Time
	LastUpdateTime time.Time
	Status         MigrationStatus
}

// NewMigrationTracker creates a new migration tracker
func NewMigrationTracker(storage Storage, metrics MetricsRecorder, logger Logger) *MigrationTracker {
	tracker := &MigrationTracker{
		storage: storage,
		metrics: metrics,
		logger:  logger,
	}

	// Initialize with default progress
	initialProgress := MigrationProgress{
		StartTime:      time.Now(),
		LastUpdateTime: time.Now(),
		Status:         StatusPending,
	}
	tracker.progress.Store(initialProgress)

	return tracker
}

// getProgress retrieves the current progress
func (t *MigrationTracker) getProgress() MigrationProgress {
	return t.progress.Load().(MigrationProgress)
}

// UpdateProgress allows updating the progress from the migration process
func (t *MigrationTracker) UpdateProgress(processed, failed int64) {
	current := t.getProgress()

	updated := MigrationProgress{
		TotalRecords:   current.TotalRecords,
		ProcessedCount: processed,
		FailedCount:    failed,
		StartTime:      current.StartTime,
		LastUpdateTime: time.Now(),
		Status:         StatusRunning,
	}

	t.progress.Store(updated)
}

// storeProgress persists the progress to storage
func (t *MigrationTracker) storeProgress(ctx context.Context, progress MigrationProgress) error {
	return t.storage.SaveProgress(ctx, progress)
}

// TrackProgress periodically tracks and records migration progress
func (t *MigrationTracker) TrackProgress(ctx context.Context) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			progress := t.getProgress()

			// Update metrics
			t.metrics.SetGauge("migration_processed_count",
				float64(progress.ProcessedCount))
			t.metrics.SetGauge("migration_failed_count",
				float64(progress.FailedCount))

			// Calculate rate
			elapsed := time.Since(progress.StartTime)
			rate := float64(progress.ProcessedCount) / elapsed.Seconds()
			t.metrics.SetGauge("migration_rate", rate)

			// Estimate remaining time
			remaining := float64(progress.TotalRecords-progress.ProcessedCount) / rate
			t.metrics.SetGauge("migration_estimated_remaining_seconds", remaining)

			// Store progress
			if err := t.storeProgress(ctx, progress); err != nil {
				t.logger.Error("failed to store progress", "error", err)
			}
		}
	}
}

// Simple implementations for the interfaces

// InMemoryStorage implements Storage interface
type InMemoryStorage struct {
	savedProgress MigrationProgress
}

func (s *InMemoryStorage) SaveProgress(ctx context.Context, progress MigrationProgress) error {
	s.savedProgress = progress
	return nil
}

// ConsoleMetricsRecorder implements MetricsRecorder interface
type ConsoleMetricsRecorder struct{}

func (m *ConsoleMetricsRecorder) SetGauge(name string, value float64) {
	fmt.Printf("METRIC: %s = %.2f\n", name, value)
}

// ConsoleLogger implements Logger interface
type ConsoleLogger struct{}

func (l *ConsoleLogger) Info(msg string, args ...interface{}) {
	fmt.Printf("INFO: "+msg+"\n", args...)
}

func (l *ConsoleLogger) Error(msg string, args ...interface{}) {
	fmt.Printf("ERROR: "+msg+"\n", args...)
}

func main() {
	// Create dependencies
	storage := &InMemoryStorage{}
	metrics := &ConsoleMetricsRecorder{}
	logger := &ConsoleLogger{}

	// Create tracker
	tracker := NewMigrationTracker(storage, metrics, logger)

	// Set initial values
	progress := tracker.getProgress()
	progress.TotalRecords = 1000
	tracker.progress.Store(progress)

	// Create a cancelable context
	ctx, cancel := context.WithCancel(context.Background())

	// Start tracking in a goroutine
	go func() {
		if err := tracker.TrackProgress(ctx); err != nil {
			logger.Error("Tracking stopped: %v", err)
		}
	}()

	// Simulate a migration process
	go func() {
		for i := 0; i < 100; i++ {
			// Update progress with some processed and failed records
			tracker.UpdateProgress(int64(i*10), int64(i))
			time.Sleep(1 * time.Second)
		}

		// Complete the migration
		finalProgress := tracker.getProgress()
		finalProgress.Status = StatusCompleted
		tracker.progress.Store(finalProgress)

		// Stop tracking
		cancel()
	}()

	// Wait for completion
	time.Sleep(120 * time.Second)
}
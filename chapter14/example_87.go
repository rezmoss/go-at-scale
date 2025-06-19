// Example 87
package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// ErrChecksumMismatch is returned when source and destination checksums don't match
var ErrChecksumMismatch = fmt.Errorf("checksum mismatch between source and destination")

// Record represents a data record
type Record struct {
	ID   string
	Data []byte
}

// DataReader interface for reading data
type DataReader interface {
	Stream(ctx context.Context) <-chan Record
	Read(ctx context.Context) ([]byte, error)
}

// MetricsRecorder interface for recording metrics
type MetricsRecorder interface {
	IncCounter(name string)
}

// Logger interface for logging
type Logger interface {
	Error(msg string, err error)
	Info(msg string)
}

// DataValidator struct remains the same
type DataValidator struct {
	source      DataReader
	destination DataReader
	metrics     MetricsRecorder
	logger      Logger
}

// ValidateMigration implements the main validation logic
func (v *DataValidator) ValidateMigration(ctx context.Context) error {
	// Get source and destination checksums
	sourceHash, err := v.calculateChecksum(ctx, v.source)
	if err != nil {
		return fmt.Errorf("calculating source checksum: %w", err)
	}

	destHash, err := v.calculateChecksum(ctx, v.destination)
	if err != nil {
		return fmt.Errorf("calculating destination checksum: %w", err)
	}

	// Compare checksums
	if sourceHash != destHash {
		v.metrics.IncCounter("validation_failures")
		return ErrChecksumMismatch
	}

	// Validate data integrity
	if err := v.validateDataIntegrity(ctx); err != nil {
		return fmt.Errorf("validating data integrity: %w", err)
	}

	// Validate business rules
	if err := v.validateBusinessRules(ctx); err != nil {
		return fmt.Errorf("validating business rules: %w", err)
	}

	return nil
}

func (v *DataValidator) validateDataIntegrity(ctx context.Context) error {
	stream := v.source.Stream(ctx)
	for record := range stream {
		if err := v.validateRecord(ctx, record); err != nil {
			v.metrics.IncCounter("record_validation_failures")
			return err
		}
	}
	return nil
}

func (v *DataValidator) calculateChecksum(ctx context.Context, reader DataReader) (string, error) {
	data, err := reader.Read(ctx)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

func (v *DataValidator) validateRecord(ctx context.Context, record Record) error {
	if record.ID == "" || len(record.Data) == 0 {
		return fmt.Errorf("invalid record: empty ID or data")
	}
	return nil
}

func (v *DataValidator) validateBusinessRules(ctx context.Context) error {
	// Simple implementation for demonstration
	return nil
}

// Example implementations of interfaces for testing
type MockDataReader struct {
	data []byte
}

func (m *MockDataReader) Stream(ctx context.Context) <-chan Record {
	ch := make(chan Record)
	go func() {
		defer close(ch)
		ch <- Record{ID: "1", Data: []byte("test data")}
	}()
	return ch
}

func (m *MockDataReader) Read(ctx context.Context) ([]byte, error) {
	return m.data, nil
}

type MockMetricsRecorder struct{}

func (m *MockMetricsRecorder) IncCounter(name string) {
	fmt.Printf("Incrementing counter: %s\n", name)
}

type MockLogger struct{}

func (m *MockLogger) Error(msg string, err error) {
	fmt.Printf("Error: %s - %v\n", msg, err)
}

func (m *MockLogger) Info(msg string) {
	fmt.Printf("Info: %s\n", msg)
}

func main() {
	validator := &DataValidator{
		source:      &MockDataReader{data: []byte("test data")},
		destination: &MockDataReader{data: []byte("test data")},
		metrics:     &MockMetricsRecorder{},
		logger:      &MockLogger{},
	}

	ctx := context.Background()
	if err := validator.ValidateMigration(ctx); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
		return
	}
	fmt.Println("Validation successful")
}
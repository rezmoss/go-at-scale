// Example 86
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Migration represents a single database migration
type Migration struct {
	ID   int
	Name string
	Up   func(ctx context.Context, tx *sql.Tx) error
	Down func(ctx context.Context, tx *sql.Tx) error
}

// SchemaValidator interface for validating schema integrity
type SchemaValidator interface {
	Validate(ctx context.Context, tx *sql.Tx) error
}

// MetricsRecorder interface for recording metrics
type MetricsRecorder interface {
	IncCounter(name string, labelName string, labelValue string)
	ObserveLatency(name string, duration time.Duration, labelName string, labelValue string)
}

// Logger interface for logging
type Logger interface {
	Error(msg string, keysAndValues ...interface{})
}

// SimpleLogger implements Logger interface
type SimpleLogger struct{}

func (l *SimpleLogger) Error(msg string, keysAndValues ...interface{}) {
	log.Printf("ERROR: %s %v", msg, keysAndValues)
}

// SimpleMetrics implements MetricsRecorder interface
type SimpleMetrics struct{}

func (m *SimpleMetrics) IncCounter(name string, labelName string, labelValue string) {
	log.Printf("Metric: %s increased [%s=%s]", name, labelName, labelValue)
}

func (m *SimpleMetrics) ObserveLatency(name string, duration time.Duration, labelName string, labelValue string) {
	log.Printf("Metric: %s = %v [%s=%s]", name, duration, labelName, labelValue)
}

// SimpleValidator implements SchemaValidator interface
type SimpleValidator struct{}

func (v *SimpleValidator) Validate(ctx context.Context, tx *sql.Tx) error {
	// In a real implementation, this would check schema integrity
	return nil
}

// SchemaManager manages database schema migrations
type SchemaManager struct {
	db         *sql.DB
	migrations []Migration
	validator  SchemaValidator
	metrics    MetricsRecorder
	logger     Logger
}

func (sm *SchemaManager) lockMigrationsTable(ctx context.Context, tx *sql.Tx) error {
	// In SQLite, we can simulate a lock by creating a table if it doesn't exist
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at TIMESTAMP NOT NULL
		)
	`)
	return err
}

func (sm *SchemaManager) isMigrationApplied(ctx context.Context, tx *sql.Tx, migrationID int) (bool, error) {
	var count int
	err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE id = ?", migrationID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("checking if migration is applied: %w", err)
	}
	return count > 0, nil
}

func (sm *SchemaManager) recordMigration(ctx context.Context, tx *sql.Tx, migration Migration) error {
	_, err := tx.ExecContext(ctx,
		"INSERT INTO schema_migrations (id, name, applied_at) VALUES (?, ?, ?)",
		migration.ID, migration.Name, time.Now())
	if err != nil {
		return fmt.Errorf("recording migration: %w", err)
	}
	return nil
}

func (sm *SchemaManager) ApplyMigrations(ctx context.Context) error {
	// Start transaction
	tx, err := sm.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	// Lock migrations table
	if err := sm.lockMigrationsTable(ctx, tx); err != nil {
		return fmt.Errorf("locking migrations table: %w", err)
	}

	// Apply each migration
	for _, migration := range sm.migrations {
		start := time.Now()

		if err := sm.applyMigration(ctx, tx, migration); err != nil {
			sm.metrics.IncCounter("migration_failures",
				"migration", migration.Name)
			return fmt.Errorf("applying migration %s: %w",
				migration.Name, err)
		}

		sm.metrics.ObserveLatency("migration_duration",
			time.Since(start),
			"migration", migration.Name)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func (sm *SchemaManager) applyMigration(ctx context.Context, tx *sql.Tx,
	migration Migration) error {
	// Check if already applied
	applied, err := sm.isMigrationApplied(ctx, tx, migration.ID)
	if err != nil {
		return err
	}
	if applied {
		return nil
	}

	// Apply forward migration
	if err := migration.Up(ctx, tx); err != nil {
		if rbErr := migration.Down(ctx, tx); rbErr != nil {
			sm.logger.Error("rollback failed",
				"error", rbErr,
				"migration", migration.Name)
		}
		return err
	}

	// Record migration
	return sm.recordMigration(ctx, tx, migration)
}

func main() {
	// Open SQLite database connection
	db, err := sql.Open("sqlite3", "./migrations.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create migrations
	migrations := []Migration{
		{
			ID:   1,
			Name: "create_users_table",
			Up: func(ctx context.Context, tx *sql.Tx) error {
				_, err := tx.ExecContext(ctx, `
					CREATE TABLE IF NOT EXISTS users (
						id INTEGER PRIMARY KEY,
						username TEXT NOT NULL,
						email TEXT NOT NULL,
						created_at TIMESTAMP NOT NULL
					)
				`)
				return err
			},
			Down: func(ctx context.Context, tx *sql.Tx) error {
				_, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS users")
				return err
			},
		},
		{
			ID:   2,
			Name: "add_user_status",
			Up: func(ctx context.Context, tx *sql.Tx) error {
				_, err := tx.ExecContext(ctx, `
					ALTER TABLE users ADD COLUMN status TEXT DEFAULT 'active'
				`)
				return err
			},
			Down: func(ctx context.Context, tx *sql.Tx) error {
				// SQLite doesn't support dropping columns directly
				// In a real implementation, we would recreate the table without the column
				return nil
			},
		},
	}

	// Initialize SchemaManager
	schemaManager := &SchemaManager{
		db:         db,
		migrations: migrations,
		validator:  &SimpleValidator{},
		metrics:    &SimpleMetrics{},
		logger:     &SimpleLogger{},
	}

	// Apply migrations
	ctx := context.Background()
	if err := schemaManager.ApplyMigrations(ctx); err != nil {
		log.Fatalf("Error applying migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
}
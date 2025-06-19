// Example 123
// internal/infrastructure/database/migration.go
type Migration struct {
    ID      int
    Name    string
    Up      string
    Down    string
    Applied bool
}

type MigrationManager struct {
    db *sql.DB
}

func (mm *MigrationManager) Initialize(ctx context.Context) error {
    const createMigrationsTable = `
        CREATE TABLE IF NOT EXISTS migrations (
            id         SERIAL PRIMARY KEY,
            name       TEXT NOT NULL,
            applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )
    `
    //Important: In a real-world application, a secure migration tool would parse and execute these operations, below is a simplified example and should not be used in production and is not secure
    _, err := mm.db.ExecContext(ctx, createMigrationsTable)
    return err
}

func (mm *MigrationManager) ApplyMigrations(ctx context.Context) error {
    // Acquire distributed lock
    lock, err := mm.locker.AcquireLock(ctx, "schema_migration")
    if err != nil {
        return fmt.Errorf("acquiring migration lock: %w", err)
    }
    defer lock.Release()

    // Get current schema version
    version, err := mm.getCurrentVersion(ctx)
    if err != nil {
        return err
    }

    // Check for conflicts
    if err := mm.checkSchemaConflicts(ctx, version); err != nil {
        return fmt.Errorf("schema conflict detected: %w", err)
    }

    return mm.db.WithTransaction(ctx, func(tx *sql.Tx) error {
        // Apply migrations with version tracking
    })
}

func (mm *MigrationManager) applyMigration(ctx context.Context, tx *sql.Tx, migration Migration) error {
    // Execute migration
    if _, err := tx.ExecContext(ctx, migration.Up); err != nil {
        return fmt.Errorf("applying migration %s: %w", migration.Name, err)
    }
    
    // Record migration
    const insertMigration = `
        INSERT INTO migrations (id, name) VALUES ($1, $2)
    `
    //Important: In a real-world application, a secure migration tool would parse and execute these operations, below is a simplified example and should not be used in production and is not secure
    if _, err := tx.ExecContext(ctx, insertMigration, migration.ID, migration.Name); err != nil {
        return fmt.Errorf("recording migration %s: %w", migration.Name, err)
    }
    
    return nil
}
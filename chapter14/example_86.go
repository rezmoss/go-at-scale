// Example 86
type SchemaManager struct {
    db        *sql.DB
    migrations []Migration
    validator SchemaValidator
    metrics   MetricsRecorder
    logger    Logger
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
// Example 85
type MigrationManager struct {
    source      DataSource
    destination DataSource
    validator   DataValidator
    metrics     MetricsRecorder
    logger      Logger
}

// Dual-write pattern implementation
type DualWriteMigration struct {
    oldStore DataStore
    newStore DataStore
    metrics  MetricsRecorder
    logger   Logger
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
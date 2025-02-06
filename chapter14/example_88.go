// Example 88
type MigrationTracker struct {
    storage  Storage
    metrics  MetricsRecorder
    logger   Logger
    progress atomic.Value
}

type MigrationProgress struct {
    TotalRecords   int64
    ProcessedCount int64
    FailedCount    int64
    StartTime      time.Time
    LastUpdateTime time.Time
    Status         MigrationStatus
}

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
            remaining := float64(progress.TotalRecords - progress.ProcessedCount) / rate
            t.metrics.SetGauge("migration_estimated_remaining_seconds", remaining)

            // Store progress
            if err := t.storeProgress(ctx, progress); err != nil {
                t.logger.Error("failed to store progress", "error", err)
            }
        }
    }
}
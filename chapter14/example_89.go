// Example 89
type RollbackManager struct {
    db          *sql.DB
    snapshotter Snapshotter
    metrics     MetricsRecorder
    logger      Logger
}

func (rm *RollbackManager) PrepareRollback(ctx context.Context) error {
    // Create snapshot
    snapshot, err := rm.snapshotter.CreateSnapshot(ctx)
    if err != nil {
        return fmt.Errorf("creating snapshot: %w", err)
    }

    // Store rollback information
    if err := rm.storeRollbackInfo(ctx, snapshot); err != nil {
        return fmt.Errorf("storing rollback info: %w", err)
    }

    return nil
}

func (rm *RollbackManager) ExecuteRollback(ctx context.Context) error {
    start := time.Now()
    defer func() {
        rm.metrics.ObserveLatency("rollback_duration", time.Since(start))
    }()

    // Get rollback information
    info, err := rm.getRollbackInfo(ctx)
    if err != nil {
        return fmt.Errorf("getting rollback info: %w", err)
    }

    // Execute rollback steps
    for _, step := range info.Steps {
        if err := rm.executeRollbackStep(ctx, step); err != nil {
            rm.metrics.IncCounter("rollback_step_failures")
            return fmt.Errorf("executing rollback step: %w", err)
        }
    }

    return nil
}
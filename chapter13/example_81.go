// Example 81
type ReplicationManager struct {
    source    DataSource
    replicas  []DataSource
    monitor   ConsistencyMonitor
    metrics   MetricsRecorder
    logger    Logger
}

func (r *ReplicationManager) ProcessTransaction(ctx context.Context, tx *Transaction) error {
    start := time.Now()
    defer func() {
        r.metrics.ObserveLatency("replication_latency", time.Since(start))
    }()

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
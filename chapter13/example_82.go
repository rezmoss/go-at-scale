// Example 82
type FailoverManager struct {
    primary     *ServiceInstance
    standby     []*ServiceInstance
    discovery   ServiceDiscovery
    healthCheck HealthChecker
    metrics     MetricsRecorder
    logger      Logger
}

func (f *FailoverManager) MonitorAndFailover(ctx context.Context) error {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-ticker.C:
            healthy, err := f.healthCheck.CheckPrimary(ctx)
            if err != nil || !healthy {
                if err := f.initiateFailover(ctx); err != nil {
                    f.metrics.IncCounter("failover_failures")
                    f.logger.Error("failover failed", "error", err)
                    continue
                }
            }
        }
    }
}

func (f *FailoverManager) initiateFailover(ctx context.Context) error {
    start := time.Now()
    defer func() {
        f.metrics.ObserveLatency("failover_duration", time.Since(start))
    }()

    // Select new primary
    newPrimary, err := f.selectNewPrimary(ctx)
    if err != nil {
        return fmt.Errorf("selecting new primary: %w", err)
    }

    // Update service discovery
    if err := f.discovery.UpdatePrimary(ctx, newPrimary); err != nil {
        return fmt.Errorf("updating service discovery: %w", err)
    }

    // Promote standby to primary
    if err := f.promoteStandby(ctx, newPrimary); err != nil {
        return fmt.Errorf("promoting standby: %w", err)
    }

    f.metrics.IncCounter("successful_failovers")
    return nil
}
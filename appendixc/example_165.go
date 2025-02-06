// Example 165
// internal/cloud/multiregion/manager.go
type RegionManager struct {
    regions        []Region
    healthCheckers map[string]*HealthChecker
    router         *TrafficRouter
    logger         Logger
    metrics        MetricsRecorder
}

type Region struct {
    Name     string
    Priority int
    Weight   float64
    Healthy  bool
}

func (rm *RegionManager) Start(ctx context.Context) error {
    // Initialize health checkers for each region
    for _, region := range rm.regions {
        checker := NewHealthChecker(HealthCheckConfig{
            Region:       region.Name,
            Interval:     30 * time.Second,
            Timeout:      5 * time.Second,
            FailureThreshold: 3,
        })
        rm.healthCheckers[region.Name] = checker

        go func(r Region) {
            for {
                select {
                case health := <-checker.HealthC():
                    rm.handleHealthUpdate(r.Name, health)
                case <-ctx.Done():
                    return
                }
            }
        }(region)
    }

    return nil
}

func (rm *RegionManager) handleHealthUpdate(region string, healthy bool) {
    rm.metrics.IncCounter(fmt.Sprintf("region_%s_health_check", region))
    
    previousState := rm.regions[region].Healthy
    rm.regions[region].Healthy = healthy

    if previousState != healthy {
        rm.logger.Info("region health changed",
            "region", region,
            "healthy", healthy)
        
        // Update routing weights
        rm.updateRoutingWeights()
    }
}

func (rm *RegionManager) updateRoutingWeights() {
    var totalWeight float64
    weights := make(map[string]float64)

    // Calculate weights based on health and priority
    for _, region := range rm.regions {
        if !region.Healthy {
            weights[region.Name] = 0
            continue
        }

        weight := region.Weight * float64(region.Priority)
        weights[region.Name] = weight
        totalWeight += weight
    }

    // Normalize weights
    for region := range weights {
        if totalWeight > 0 {
            weights[region] /= totalWeight
        }
    }

    rm.router.UpdateWeights(weights)
}
// Example 147
// internal/deployment/strategy.go
type DeploymentStrategy interface {
    Deploy(ctx context.Context) error
    Rollback(ctx context.Context) error
}

// Blue-Green deployment
type BlueGreenDeployment struct {
    k8sClient  *kubernetes.Clientset
    newVersion string
    oldVersion string
}

func (d *BlueGreenDeployment) Deploy(ctx context.Context) error {
    // Deploy new version
    if err := d.deployNewVersion(ctx); err != nil {
        return fmt.Errorf("deploying new version: %w", err)
    }
    
    // Wait for new version to be ready
    if err := d.waitForReadiness(ctx); err != nil {
        return fmt.Errorf("waiting for readiness: %w", err)
    }
    
    // Switch traffic
    if err := d.switchTraffic(ctx); err != nil {
        return fmt.Errorf("switching traffic: %w", err)
    }
    
    // Clean up old version
    if err := d.cleanup(ctx); err != nil {
        log.Printf("cleanup failed: %v", err)
    }
    
    return nil
}

// Canary deployment
type CanaryDeployment struct {
    k8sClient  *kubernetes.Clientset
    newVersion string
    steps      []float64 // Traffic percentages
    interval   time.Duration
}

func (d *CanaryDeployment) Deploy(ctx context.Context) error {
    // Deploy canary version
    if err := d.deployCanary(ctx); err != nil {
        return fmt.Errorf("deploying canary: %w", err)
    }
    
    // Gradually increase traffic
    for _, percentage := range d.steps {
        if err := d.setTrafficPercentage(ctx, percentage); err != nil {
            return fmt.Errorf("setting traffic percentage: %w", err)
        }
        
        // Monitor health
        if err := d.monitorHealth(ctx); err != nil {
            return fmt.Errorf("health check failed: %w", err)
        }
        
        time.Sleep(d.interval)
    }
    
    // Finalize deployment
    return d.finalize(ctx)
}
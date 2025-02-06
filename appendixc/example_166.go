// Example 166
// internal/cloud/autoscaling/scaler.go
type AutoScaler struct {
    client    kubernetes.Interface
    metrics   MetricsClient
    config    ScalerConfig
    logger    Logger
}

type ScalerConfig struct {
    Namespace     string
    Deployment    string
    MinReplicas   int32
    MaxReplicas   int32
    TargetCPU     int32
    ScaleUpCooldown   time.Duration
    ScaleDownCooldown time.Duration
}

func (s *AutoScaler) Start(ctx context.Context) error {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            if err := s.evaluate(ctx); err != nil {
                s.logger.Error("scaling evaluation failed", "error", err)
            }
        case <-ctx.Done():
            return nil
        }
    }
}

func (s *AutoScaler) evaluate(ctx context.Context) error {
    // Get current metrics
    cpu, err := s.metrics.GetCPUUtilization(ctx, s.config.Namespace, s.config.Deployment)
    if err != nil {
        return fmt.Errorf("getting CPU metrics: %w", err)
    }

    // Get current deployment
    deployment, err := s.client.AppsV1().Deployments(s.config.Namespace).Get(ctx, s.config.Deployment, metav1.GetOptions{})
    if err != nil {
        return fmt.Errorf("getting deployment: %w", err)
    }

    currentReplicas := *deployment.Spec.Replicas
    desiredReplicas := s.calculateDesiredReplicas(cpu, currentReplicas)

    // Apply scaling limits
    desiredReplicas = int32(math.Max(float64(s.config.MinReplicas),
        math.Min(float64(s.config.MaxReplicas), float64(desiredReplicas))))

    if desiredReplicas != currentReplicas {
        // Check cooldown periods
        if !s.canScale(desiredReplicas > currentReplicas) {
            return nil
        }

        // Update deployment
        deployment.Spec.Replicas = &desiredReplicas
        if _, err := s.client.AppsV1().Deployments(s.config.Namespace).Update(ctx, deployment, metav1.UpdateOptions{}); err != nil {
            return fmt.Errorf("updating deployment: %w", err)
        }

        s.logger.Info("scaled deployment",
            "from", currentReplicas,
            "to", desiredReplicas,
            "cpu_utilization", cpu)
    }

    return nil
}
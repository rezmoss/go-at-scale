// Example 162
// internal/testing/chaos/executor.go
type ChaosExecutor struct {
    k8sClient  kubernetes.Interface
    metrics    MetricsRecorder
    logger     Logger
}

type ChaosConfig struct {
    TargetNamespace string
    Duration        time.Duration
    Failures        []FailureSpec
}

type FailureSpec struct {
    Type     string
    Selector map[string]string
    Rate     float64
}

func (e *ChaosExecutor) InjectFailures(ctx context.Context, config ChaosConfig) error {
    for _, failure := range config.Failures {
        switch failure.Type {
        case "pod-kill":
            if err := e.injectPodFailure(ctx, config.TargetNamespace, failure); err != nil {
                return fmt.Errorf("injecting pod failure: %w", err)
            }
        case "network-latency":
            if err := e.injectNetworkLatency(ctx, config.TargetNamespace, failure); err != nil {
                return fmt.Errorf("injecting network latency: %w", err)
            }
        case "cpu-pressure":
            if err := e.injectCPUPressure(ctx, config.TargetNamespace, failure); err != nil {
                return fmt.Errorf("injecting CPU pressure: %w", err)
            }
        }
    }

    // Wait for duration
    time.Sleep(config.Duration)

    // Cleanup
    if err := e.cleanup(ctx, config.TargetNamespace); err != nil {
        return fmt.Errorf("cleanup failed: %w", err)
    }

    return nil
}

func (e *ChaosExecutor) injectPodFailure(ctx context.Context, namespace string, spec FailureSpec) error {
    pods, err := e.k8sClient.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
        LabelSelector: labels.SelectorFromSet(spec.Selector).String(),
    })
    if err != nil {
        return fmt.Errorf("listing pods: %w", err)
    }

    // Randomly select pods based on rate
    for _, pod := range pods.Items {
        if rand.Float64() < spec.Rate {
            err := e.k8sClient.CoreV1().Pods(namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{})
            if err != nil {
                return fmt.Errorf("deleting pod %s: %w", pod.Name, err)
            }
            e.metrics.IncCounter("chaos_pod_kills")
            e.logger.Info("killed pod", "pod", pod.Name)
        }
    }

    return nil
}
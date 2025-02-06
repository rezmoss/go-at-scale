// Example 84
type RecoveryTester struct {
    scenarios []RecoveryScenario
    executor  ScenarioExecutor
    validator RecoveryValidator
    metrics   MetricsRecorder
    logger    Logger
}

func (t *RecoveryTester) TestRecovery(ctx context.Context) error {
    for _, scenario := range t.scenarios {
        if err := t.runScenario(ctx, scenario); err != nil {
            t.metrics.IncCounter("scenario_failures",
                "scenario", scenario.Name)
            return fmt.Errorf("running scenario %s: %w", scenario.Name, err)
        }
    }
    return nil
}

func (t *RecoveryTester) runScenario(ctx context.Context, scenario RecoveryScenario) error {
    start := time.Now()
    defer func() {
        t.metrics.ObserveLatency("scenario_duration", time.Since(start))
    }()

    // Initialize test environment
    env, err := t.executor.SetupEnvironment(ctx, scenario)
    if err != nil {
        return fmt.Errorf("setting up environment: %w", err)
    }
    defer env.Cleanup()

    // Execute scenario
    recovery, err := t.executor.ExecuteScenario(ctx, env, scenario)
    if err != nil {
        return fmt.Errorf("executing scenario: %w", err)
    }

    // Validate recovery
    if err := t.validator.ValidateRecovery(ctx, recovery); err != nil {
        return fmt.Errorf("validating recovery: %w", err)
    }

    t.metrics.IncCounter("successful_scenarios",
        "scenario", scenario.Name)
    return nil
}
// Example 161
// internal/testing/performance/loadtest.go
type LoadTest struct {
    config  LoadTestConfig
    metrics MetricsRecorder
    logger  Logger
}

type LoadTestConfig struct {
    BaseURL       string
    NumUsers      int
    RampUpTime    time.Duration
    Duration      time.Duration
    ThinkTime     time.Duration
    Scenarios     []Scenario
}

type Scenario struct {
    Name       string
    Weight     int
    Executable func(context.Context) error
}

func (lt *LoadTest) Run(ctx context.Context) error {
    start := time.Now()
    defer func() {
        lt.metrics.ObserveLatency("load_test_duration", time.Since(start))
    }()

    // Calculate user start times
    startTimes := lt.calculateStartTimes()

    // Create user groups
    var wg sync.WaitGroup
    errors := make(chan error, lt.config.NumUsers)

    for i := 0; i < lt.config.NumUsers; i++ {
        wg.Add(1)
        go func(userID int, startDelay time.Duration) {
            defer wg.Done()

            // Wait for scheduled start time
            time.Sleep(startDelay)

            // Run user scenarios
            if err := lt.runUserScenarios(ctx, userID); err != nil {
                errors <- fmt.Errorf("user %d: %w", userID, err)
            }
        }(i, startTimes[i])
    }

    // Wait for completion
    wg.Wait()
    close(errors)

    // Collect errors
    var errs []error
    for err := range errors {
        errs = append(errs, err)
    }

    if len(errs) > 0 {
        return fmt.Errorf("load test failures: %v", errs)
    }

    // errors.Join returns an error wrapping multiple underlying errors.
    // You can still use errors.Is/errors.As on the joined error.
    // If using Go 1.20 or newer:
    // if len(errs) > 0 {
    //     return errors.Join(errs...)
    // }

    return nil
}

func (lt *LoadTest) runUserScenarios(ctx context.Context, userID int) error {
    client := &http.Client{
        Timeout: 30 * time.Second,
    }

    endTime := time.Now().Add(lt.config.Duration)
    for time.Now().Before(endTime) {
        // Select scenario based on weights
        scenario := lt.selectScenario()

        // Execute scenario
        start := time.Now()
        err := scenario.Executable(ctx)
        duration := time.Since(start)

        // Record metrics
        lt.metrics.ObserveLatency(fmt.Sprintf("scenario_%s", scenario.Name), duration)
        if err != nil {
            lt.metrics.IncCounter(fmt.Sprintf("scenario_%s_errors", scenario.Name))
            lt.logger.Error("scenario failed",
                "user_id", userID,
                "scenario", scenario.Name,
                "error", err)
        }

        // Think time
        time.Sleep(lt.config.ThinkTime)
    }

    return nil
}
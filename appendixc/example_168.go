// Example 168
// internal/cloud/secrets/manager.go
type SecretManager struct {
    client    secretmanager.Client
    cache     *sync.Map
    config    SecretConfig
    logger    Logger
    metrics   MetricsRecorder
}

type SecretConfig struct {
    Project     string
    Region      string
    CacheTTL    time.Duration
}

type cachedSecret struct {
    value      string
    expiration time.Time
}

func (sm *SecretManager) GetSecret(ctx context.Context, name string) (string, error) {
    start := time.Now()
    defer func() {
        sm.metrics.ObserveLatency("secret_retrieval", time.Since(start))
    }()

    // Check cache
    if value, ok := sm.checkCache(name); ok {
        sm.metrics.IncCounter("secret_cache_hits")
        return value, nil
    }

    // Build the secret name
    secretName := fmt.Sprintf("projects/%s/secrets/%s/versions/latest",
        sm.config.Project, name)

    // Access the secret version
    result, err := sm.client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
        Name: secretName,
    })
    if err != nil {
        sm.metrics.IncCounter("secret_retrieval_failures")
        return "", fmt.Errorf("accessing secret: %w", err)
    }

    secretValue := string(result.Payload.Data)

    // Cache the result
    sm.cache.Store(name, cachedSecret{
        value:      secretValue,
        expiration: time.Now().Add(sm.config.CacheTTL),
    })

    sm.metrics.IncCounter("secret_retrievals")
    return secretValue, nil
}

func (sm *SecretManager) checkCache(name string) (string, bool) {
    if value, ok := sm.cache.Load(name); ok {
        cached := value.(cachedSecret)
        if time.Now().Before(cached.expiration) {
            return cached.value, true
        }
        // Expired
        sm.cache.Delete(name)
    }
    return "", false
}
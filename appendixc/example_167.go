// Example 167
// internal/cloud/storage/manager.go
type StorageManager struct {
    client    *storage.Client
    config    StorageConfig
    cache     Cache
    logger    Logger
    metrics   MetricsRecorder
}

type StorageConfig struct {
    Bucket          string
    Region          string
    CacheEnabled    bool
    CacheTTL        time.Duration
    RetryAttempts   int
    RetryDelay      time.Duration
}

func (sm *StorageManager) Upload(ctx context.Context, key string, data []byte) error {
    start := time.Now()
    defer func() {
        sm.metrics.ObserveLatency("storage_upload", time.Since(start))
    }()

    // Create bucket handle
    bucket := sm.client.Bucket(sm.config.Bucket)

    // Upload with retry
    var err error
    for attempt := 1; attempt <= sm.config.RetryAttempts; attempt++ {
        obj := bucket.Object(key)
        writer := obj.NewWriter(ctx)

        if _, err = writer.Write(data); err != nil {
            sm.logger.Error("upload failed",
                "attempt", attempt,
                "error", err)
            
            time.Sleep(sm.config.RetryDelay)
            continue
        }

        if err = writer.Close(); err != nil {
            sm.logger.Error("closing writer failed",
                "attempt", attempt,
                "error", err)
            
            time.Sleep(sm.config.RetryDelay)
            continue
        }

        // Upload successful
        break
    }

    if err != nil {
        sm.metrics.IncCounter("storage_upload_failures")
        return fmt.Errorf("upload failed after %d attempts: %w",
            sm.config.RetryAttempts, err)
    }

    // Invalidate cache if enabled
    if sm.config.CacheEnabled {
        sm.cache.Delete(key)
    }

    sm.metrics.IncCounter("storage_uploads")
    return nil
}

func (sm *StorageManager) Download(ctx context.Context, key string) ([]byte, error) {
    start := time.Now()
    defer func() {
        sm.metrics.ObserveLatency("storage_download", time.Since(start))
    }()

    // Check cache first if enabled
    if sm.config.CacheEnabled {
        if data, found := sm.cache.Get(key); found {
            sm.metrics.IncCounter("storage_cache_hits")
            return data.([]byte), nil
        }
    }

    // Create bucket handle
    bucket := sm.client.Bucket(sm.config.Bucket)

    // Download with retry
    var data []byte
    var err error
    for attempt := 1; attempt <= sm.config.RetryAttempts; attempt++ {
        obj := bucket.Object(key)
        reader, err := obj.NewReader(ctx)
        if err != nil {
            if err == storage.ErrObjectNotExist {
                return nil, err
            }
            
            sm.logger.Error("creating reader failed",
                "attempt", attempt,
                "error", err)
            
            time.Sleep(sm.config.RetryDelay)
            continue
        }
        defer reader.Close()

        data, err = io.ReadAll(reader)
        if err != nil {
            sm.logger.Error("reading data failed",
                "attempt", attempt,
                "error", err)
            
            time.Sleep(sm.config.RetryDelay)
            continue
        }

        // Download successful
        break
    }

    if err != nil {
        sm.metrics.IncCounter("storage_download_failures")
        return nil, fmt.Errorf("download failed after %d attempts: %w",
            sm.config.RetryAttempts, err)
    }

    // Update cache if enabled
    if sm.config.CacheEnabled {
        sm.cache.Set(key, data, sm.config.CacheTTL)
    }

    sm.metrics.IncCounter("storage_downloads")
    return data, nil
}
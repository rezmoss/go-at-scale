// Example 127
// internal/infrastructure/cache/consistency.go
type WriteThrough struct {
    cache       *RedisCache
    repository  Repository
    eventBus    EventBus
}

func (wt *WriteThrough) Save(ctx context.Context, key string, value interface{}) error {
    // Save to database first
    if err := wt.repository.Save(ctx, value); err != nil {
        return fmt.Errorf("saving to repository: %w", err)
    }
    
    // Update cache
    if err := wt.cache.Set(ctx, key, value, DefaultTTL); err != nil {
        // Log error but don't fail the operation
        log.Printf("failed to update cache: %v", err)
    }
    
    // Notify other instances
    event := CacheUpdateEvent{
        Key:   key,
        Value: value,
    }
    if err := wt.eventBus.Publish(ctx, "cache.updated", event); err != nil {
        log.Printf("failed to publish cache update event: %v", err)
    }
    
    return nil
}

// Cache-Aside Pattern
type CacheAside struct {
    cache      *RedisCache
    repository Repository
    loader     sync.Map // For preventing thundering herd
}

func (ca *CacheAside) Get(ctx context.Context, key string, value interface{}) error {
    // Try cache first
    err := ca.cache.Get(ctx, key, value)
    if err == nil {
        return nil
    }
    
    // Handle cache miss
    if err != ErrNotFound {
        return fmt.Errorf("getting from cache: %w", err)
    }
    
    // Check if we're already loading this key
    if loader, ok := ca.loader.Load(key); ok {
        // Wait for the other loader to complete
        <-loader.(chan struct{})
        return ca.cache.Get(ctx, key, value)
    }
    
    // Create loader
    done := make(chan struct{})
    if actual, loaded := ca.loader.LoadOrStore(key, done); loaded {
        // Another goroutine beat us to it
        <-actual.(chan struct{})
        return ca.cache.Get(ctx, key, value)
    }
    
    defer func() {
        close(done)
        ca.loader.Delete(key)
    }()
    
    // Load from repository
    if err := ca.repository.Get(ctx, key, value); err != nil {
        return fmt.Errorf("getting from repository: %w", err)
    }
    
    // Update cache
    if err := ca.cache.Set(ctx, key, value, DefaultTTL); err != nil {
        log.Printf("failed to update cache: %v", err)
    }
    
    return nil
}
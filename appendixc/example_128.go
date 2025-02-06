// Example 128
// internal/infrastructure/cache/invalidation.go
type CacheInvalidator struct {
    cache     *RedisCache
    patterns  []string
    eventBus  EventBus
}

func (ci *CacheInvalidator) InvalidatePattern(ctx context.Context, pattern string) error {
    keys, err := ci.cache.client.Keys(ctx, pattern).Result()
    if err != nil {
        return fmt.Errorf("getting keys for pattern: %w", err)
    }
    
    pipe := ci.cache.client.Pipeline()
    for _, key := range keys {
        pipe.Del(ctx, key)
    }
    
    if _, err := pipe.Exec(ctx); err != nil {
        return fmt.Errorf("deleting keys: %w", err)
    }
    
    // Notify other instances
    event := CacheInvalidationEvent{Pattern: pattern}
    if err := ci.eventBus.Publish(ctx, "cache.invalidated", event); err != nil {
        log.Printf("failed to publish invalidation event: %v", err)
    }
    
    return nil
}

// Time-based invalidation
type TTLManager struct {
    cache *RedisCache
}

func (tm *TTLManager) UpdateTTL(ctx context.Context, key string, ttl time.Duration) error {
    return tm.cache.client.Expire(ctx, key, ttl).Err()
}

func (tm *TTLManager) TouchKey(ctx context.Context, key string) error {
    return tm.cache.client.Touch(ctx, key).Err()
}
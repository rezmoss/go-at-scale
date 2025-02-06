// Example 129
// internal/infrastructure/cache/optimization.go
type CacheWarmer struct {
    cache      *RedisCache
    repository Repository
    patterns   []string
}

func (cw *CacheWarmer) WarmUp(ctx context.Context) error {
    for _, pattern := range cw.patterns {
        keys, err := cw.repository.GetKeys(ctx, pattern)
        if err != nil {
            return fmt.Errorf("getting keys for pattern %s: %w", pattern, err)
        }
        
        for _, key := range keys {
            var value interface{}
            if err := cw.repository.Get(ctx, key, &value); err != nil {
                log.Printf("failed to get value for key %s: %v", key, err)
                continue
            }
            
            if err := cw.cache.Set(ctx, key, value, DefaultTTL); err != nil {
                log.Printf("failed to warm cache for key %s: %v", key, err)
            }
        }
    }
    
    return nil
}

// Cache compression
type CompressedCache struct {
    cache     *RedisCache
    compress  func([]byte) ([]byte, error)
    decompress func([]byte) ([]byte, error)
}

func (cc *CompressedCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    data, err := cc.cache.codec.Encode(value)
    if err != nil {
        return fmt.Errorf("encoding value: %w", err)
    }
    
    compressed, err := cc.compress(data)
    if err != nil {
        return fmt.Errorf("compressing data: %w", err)
    }
    
    return cc.cache.client.Set(ctx, key, compressed, ttl).Err()
}
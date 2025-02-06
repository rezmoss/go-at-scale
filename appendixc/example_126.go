// Example 126
// internal/infrastructure/cache/redis.go
type RedisCache struct {
    client  *redis.Client
    codec   encoding.Codec
    metrics *metrics.Reporter
}

func NewRedisCache(client *redis.Client, codec encoding.Codec, metrics *metrics.Reporter) *RedisCache {
    return &RedisCache{
        client:  client,
        codec:   codec,
        metrics: metrics,
    }
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
    start := time.Now()
    defer func() {
        c.metrics.ObserveLatency("cache_set", time.Since(start))
    }()
    
    data, err := c.codec.Encode(value)
    if err != nil {
        return fmt.Errorf("encoding value: %w", err)
    }
    
    if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
        c.metrics.IncCounter("cache_set_errors")
        return fmt.Errorf("setting cache key: %w", err)
    }
    
    c.metrics.IncCounter("cache_sets")
    return nil
}

func (c *RedisCache) Get(ctx context.Context, key string, value interface{}) error {
    start := time.Now()
    defer func() {
        c.metrics.ObserveLatency("cache_get", time.Since(start))
    }()
    
    data, err := c.client.Get(ctx, key).Bytes()
    if err != nil {
        if err == redis.Nil {
            c.metrics.IncCounter("cache_misses")
            return ErrNotFound
        }
        c.metrics.IncCounter("cache_get_errors")
        return fmt.Errorf("getting cache key: %w", err)
    }
    
    if err := c.codec.Decode(data, value); err != nil {
        return fmt.Errorf("decoding value: %w", err)
    }
    
    c.metrics.IncCounter("cache_hits")
    return nil
}
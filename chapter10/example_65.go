// Example 65
type CacheResolver struct {
    underlying Resolver
    cache      *redis.Client
    ttl        time.Duration
}

func (r *CacheResolver) Query_user(ctx context.Context, id string) (*User, error) {
    cacheKey := fmt.Sprintf("user:%s", id)
    
    // Try cache first
    if cached, err := r.getFromCache(ctx, cacheKey); err == nil {
        return cached, nil
    }
    
    // Get from underlying resolver
    user, err := r.underlying.Query_user(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Cache result
    go r.cacheResult(ctx, cacheKey, user)
    
    return user, nil
}
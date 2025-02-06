// Example 125
// internal/infrastructure/cache/memory.go
type Cache[K comparable, V any] struct {
    data     map[K]cacheEntry[V]
    mu       sync.RWMutex
    maxSize  int
    onEvict  func(K, V)
}

type cacheEntry[V any] struct {
    value      V
    expiration time.Time
    lastAccess time.Time
}

func NewCache[K comparable, V any](maxSize int, onEvict func(K, V)) *Cache[K, V] {
    cache := &Cache[K, V]{
        data:    make(map[K]cacheEntry[V]),
        maxSize: maxSize,
        onEvict: onEvict,
    }
    
    // Start cleanup routine
    go cache.cleanup()
    
    return cache
}

func (c *Cache[K, V]) Set(key K, value V, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    // Check if we need to evict
    if len(c.data) >= c.maxSize {
        c.evictOldest()
    }
    
    c.data[key] = cacheEntry[V]{
        value:      value,
        expiration: time.Now().Add(ttl),
        lastAccess: time.Now(),
    }
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
    c.mu.RLock()
    entry, exists := c.data[key]
    c.mu.RUnlock()
    
    if !exists {
        var zero V
        return zero, false
    }
    
    if time.Now().After(entry.expiration) {
        c.mu.Lock()
        delete(c.data, key)
        c.mu.Unlock()
        var zero V
        return zero, false
    }
    // Update last access time
    c.mu.Lock()
    entry.lastAccess = time.Now()
    c.data[key] = entry
    c.mu.Unlock()
    
    return entry.value, true
}

func (c *Cache[K, V]) evictOldest() {
    var oldestKey K
    var oldestAccess time.Time
    // Find oldest entry
    for k, v := range c.data {
        if oldestAccess.IsZero() || v.lastAccess.Before(oldestAccess) {
            oldestKey = k
            oldestAccess = v.lastAccess
        }
    }
    // Evict oldest
    if c.onEvict != nil {
        c.onEvict(oldestKey, c.data[oldestKey].value)
    }
    delete(c.data, oldestKey)
}
// Example 51
type ExpensiveResource interface {
    Request() ([]byte, error)
}

type RealResource struct {
    data []byte
}

func (r *RealResource) Request() ([]byte, error) {
    // Expensive operation simulation
    time.Sleep(2 * time.Second)
    return r.data, nil
}

type CachingProxy struct {
    resource ExpensiveResource
    cache    map[string][]byte
    mu       sync.RWMutex
}

func (p *CachingProxy) Request() ([]byte, error) {
    p.mu.RLock()
    if data, ok := p.cache["key"]; ok {
        p.mu.RUnlock()
        return data, nil
    }
    p.mu.RUnlock()
    // Cache miss - get from real resource
    data, err := p.resource.Request()
    if err != nil {
        return nil, err
    }
    p.mu.Lock()
    p.cache["key"] = data
    p.mu.Unlock()
    
    return data, nil
}
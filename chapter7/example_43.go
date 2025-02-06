// Example 43
// Constructor injection
type Service struct {
    repo   Repository
    cache  Cache
    config ServiceConfig
}

func NewService(repo Repository, cache Cache, config ServiceConfig) *Service {
    return &Service{
        repo:   repo,
        cache:  cache,
        config: config,
    }
}

// Functional options pattern for optional dependencies
type ServiceOption func(*Service)

func WithCache(cache Cache) ServiceOption {
    return func(s *Service) {
        s.cache = cache
    }
}

func WithMetrics(metrics MetricsCollector) ServiceOption {
    return func(s *Service) {
        s.metrics = metrics
    }
}

func NewServiceWithOptions(repo Repository, opts ...ServiceOption) *Service {
    s := &Service{
        repo: repo,
    }
    
    for _, opt := range opts {
        opt(s)
    }
    
    return s
}

// Usage
service := NewServiceWithOptions(
    repo,
    WithCache(redisCache),
    WithMetrics(prometheusCollector),
)
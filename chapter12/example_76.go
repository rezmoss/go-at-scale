// Example 76
type ServiceDiscovery interface {
    Register(ctx context.Context, service Service) error
    Deregister(ctx context.Context, serviceID string) error
    GetService(ctx context.Context, name string) (*Service, error)
    GetServices(ctx context.Context) ([]Service, error)
}

type ConsulDiscovery struct {
    client   *api.Client
    cache    *sync.Map
    watcher  ServiceWatcher
    metrics  MetricsRecorder
    logger   Logger
}

func (d *ConsulDiscovery) GetService(ctx context.Context, name string) (*Service, error) {
    // Check cache first
    if service, ok := d.cache.Load(name); ok {
        return service.(*Service), nil
    }

    services, _, err := d.client.Health().Service(name, "", true, &api.QueryOptions{
        Context: ctx,
    })
    if err != nil {
        d.metrics.IncCounter("service_discovery_errors")
        return nil, fmt.Errorf("getting service: %w", err)
    }

    if len(services) == 0 {
        return nil, ErrServiceNotFound
    }

    // Apply load balancing strategy
    service := d.loadBalancer.Choose(services)
    d.cache.Store(name, service)

    return service, nil
}
// Example 115
// internal/infrastructure/discovery/consul.go
type ServiceRegistry struct {
    client    *api.Client
    serviceName string
}

func (sr *ServiceRegistry) Register(serviceID string, address string, port int) error {
    registration := &api.AgentServiceRegistration{
        ID:      serviceID,
        Name:    sr.serviceName,
        Port:    port,
        Address: address,
        Check: &api.AgentServiceCheck{
            HTTP:     fmt.Sprintf("http://%s:%d/health", address, port),
            Interval: "10s",
            Timeout:  "5s",
        },
    }
    
    return sr.client.Agent().ServiceRegister(registration)
}

func (sr *ServiceRegistry) Deregister(serviceID string) error {
    return sr.client.Agent().ServiceDeregister(serviceID)
}

// Load balancer
type LoadBalancer struct {
    services []string
    mu       sync.RWMutex
    current  int
}

func (lb *LoadBalancer) Next() string {
    lb.mu.Lock()
    defer lb.mu.Unlock()
    
    service := lb.services[lb.current]
    lb.current = (lb.current + 1) % len(lb.services)
    return service
}
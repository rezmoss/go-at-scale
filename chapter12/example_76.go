// Example 76
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/hashicorp/consul/api"
)

// Service represents a service instance
type Service struct {
	ID      string
	Name    string
	Address string
	Port    int
	Tags    []string
	Meta    map[string]string
}

// LoadBalancer interface for selecting services
type LoadBalancer interface {
	Choose(services []*api.ServiceEntry) *Service
}

// RoundRobinBalancer implements a simple round-robin selection strategy
type RoundRobinBalancer struct {
	counter int
	mu      sync.Mutex
}

func (lb *RoundRobinBalancer) Choose(services []*api.ServiceEntry) *Service {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if len(services) == 0 {
		return nil
	}

	idx := lb.counter % len(services)
	lb.counter++

	entry := services[idx]
	return &Service{
		ID:      entry.Service.ID,
		Name:    entry.Service.Service,
		Address: entry.Service.Address,
		Port:    entry.Service.Port,
		Tags:    entry.Service.Tags,
		Meta:    entry.Service.Meta,
	}
}

// Logger interface for logging
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// SimpleLogger implements the Logger interface
type SimpleLogger struct{}

func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	log.Printf("INFO: "+msg, args...)
}

func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	log.Printf("ERROR: "+msg, args...)
}

// MetricsRecorder interface for recording metrics
type MetricsRecorder interface {
	IncCounter(name string)
	RecordTiming(name string, value float64)
}

// SimpleMetrics implements a basic metrics recorder
type SimpleMetrics struct{}

func (m *SimpleMetrics) IncCounter(name string) {
	log.Printf("Metric incremented: %s", name)
}

func (m *SimpleMetrics) RecordTiming(name string, value float64) {
	log.Printf("Timing recorded: %s = %.2f", name, value)
}

// ServiceWatcher interface for watching service changes
type ServiceWatcher interface {
	Watch(ctx context.Context, service string) error
}

// SimpleWatcher implements a basic service watcher
type SimpleWatcher struct{}

func (w *SimpleWatcher) Watch(ctx context.Context, service string) error {
	log.Printf("Watching service: %s", service)
	return nil
}

// Common errors
var (
	ErrServiceNotFound = errors.New("service not found")
)

// ServiceDiscovery interface
type ServiceDiscovery interface {
	Register(ctx context.Context, service Service) error
	Deregister(ctx context.Context, serviceID string) error
	GetService(ctx context.Context, name string) (*Service, error)
	GetServices(ctx context.Context) ([]Service, error)
}

// ConsulDiscovery implements ServiceDiscovery using Consul
type ConsulDiscovery struct {
	client       *api.Client
	cache        *sync.Map
	watcher      ServiceWatcher
	metrics      MetricsRecorder
	logger       Logger
	loadBalancer LoadBalancer
}

// NewConsulDiscovery creates a new Consul service discovery instance
func NewConsulDiscovery(consulAddr string) (*ConsulDiscovery, error) {
	config := api.DefaultConfig()
	if consulAddr != "" {
		config.Address = consulAddr
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("creating consul client: %w", err)
	}

	return &ConsulDiscovery{
		client:       client,
		cache:        &sync.Map{},
		watcher:      &SimpleWatcher{},
		metrics:      &SimpleMetrics{},
		logger:       &SimpleLogger{},
		loadBalancer: &RoundRobinBalancer{},
	}, nil
}

// Register registers a service with Consul
func (d *ConsulDiscovery) Register(ctx context.Context, service Service) error {
	reg := &api.AgentServiceRegistration{
		ID:      service.ID,
		Name:    service.Name,
		Address: service.Address,
		Port:    service.Port,
		Tags:    service.Tags,
		Meta:    service.Meta,
	}

	if err := d.client.Agent().ServiceRegister(reg); err != nil {
		d.metrics.IncCounter("service_registration_errors")
		return fmt.Errorf("registering service: %w", err)
	}

	d.logger.Info("Service registered: %s (%s)", service.Name, service.ID)
	return nil
}

// Deregister removes a service from Consul
func (d *ConsulDiscovery) Deregister(ctx context.Context, serviceID string) error {
	if err := d.client.Agent().ServiceDeregister(serviceID); err != nil {
		d.metrics.IncCounter("service_deregistration_errors")
		return fmt.Errorf("deregistering service: %w", err)
	}

	d.logger.Info("Service deregistered: %s", serviceID)
	return nil
}

// GetService retrieves a service by name
func (d *ConsulDiscovery) GetService(ctx context.Context, name string) (*Service, error) {
	// Check cache first
	if service, ok := d.cache.Load(name); ok {
		return service.(*Service), nil
	}

	queryOpts := &api.QueryOptions{}
	queryOpts = queryOpts.WithContext(ctx)
	services, _, err := d.client.Health().Service(name, "", true, queryOpts)
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

// GetServices returns all registered services
func (d *ConsulDiscovery) GetServices(ctx context.Context) ([]Service, error) {
	queryOpts := &api.QueryOptions{}
	queryOpts = queryOpts.WithContext(ctx)
	serviceMap, _, err := d.client.Catalog().Services(queryOpts)
	if err != nil {
		d.metrics.IncCounter("service_discovery_errors")
		return nil, fmt.Errorf("listing services: %w", err)
	}

	var result []Service
	for serviceName := range serviceMap {
		service, err := d.GetService(ctx, serviceName)
		if err != nil {
			continue
		}
		result = append(result, *service)
	}

	return result, nil
}

func main() {
	// Create a new Consul discovery client
	discovery, err := NewConsulDiscovery("localhost:8500")
	if err != nil {
		log.Fatalf("Failed to create service discovery: %v", err)
	}

	// Register a sample service
	err = discovery.Register(context.Background(), Service{
		ID:      "web-service-1",
		Name:    "web-service",
		Address: "127.0.0.1",
		Port:    8080,
		Tags:    []string{"http", "web"},
		Meta:    map[string]string{"version": "1.0.0"},
	})
	if err != nil {
		log.Printf("Failed to register service: %v", err)
	}

	// Get a service instance
	service, err := discovery.GetService(context.Background(), "web-service")
	if err != nil {
		log.Printf("Failed to get service: %v", err)
	} else {
		log.Printf("Found service: %s at %s:%d", service.Name, service.Address, service.Port)
	}

	// For demonstration only - in a real app you'd defer this to shutdown
	err = discovery.Deregister(context.Background(), "web-service-1")
	if err != nil {
		log.Printf("Failed to deregister service: %v", err)
	}
}
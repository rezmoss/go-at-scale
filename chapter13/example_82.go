// Example 82
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
)

// ServiceInstance represents a service that can be primary or standby
type ServiceInstance struct {
	ID        string
	Endpoint  string
	IsHealthy bool
}

// ServiceDiscovery interface for updating service discovery
type ServiceDiscovery interface {
	UpdatePrimary(ctx context.Context, instance *ServiceInstance) error
}

// HealthChecker interface for checking health of primary
type HealthChecker interface {
	CheckPrimary(ctx context.Context) (bool, error)
}

// MetricsRecorder interface for recording metrics
type MetricsRecorder interface {
	IncCounter(metric string)
	ObserveLatency(metric string, duration time.Duration)
}

// Logger interface for logging
type Logger interface {
	Error(msg string, args ...interface{})
	Info(msg string, args ...interface{})
}

type FailoverManager struct {
	primary     *ServiceInstance
	standby     []*ServiceInstance
	discovery   ServiceDiscovery
	healthCheck HealthChecker
	metrics     MetricsRecorder
	logger      Logger
}

func (f *FailoverManager) MonitorAndFailover(ctx context.Context) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			healthy, err := f.healthCheck.CheckPrimary(ctx)
			if err != nil || !healthy {
				if err := f.initiateFailover(ctx); err != nil {
					f.metrics.IncCounter("failover_failures")
					f.logger.Error("failover failed", "error", err)
					continue
				}
			}
		}
	}
}

func (f *FailoverManager) initiateFailover(ctx context.Context) error {
	start := time.Now()
	defer func() {
		f.metrics.ObserveLatency("failover_duration", time.Since(start))
	}()

	// Select new primary
	newPrimary, err := f.selectNewPrimary(ctx)
	if err != nil {
		return fmt.Errorf("selecting new primary: %w", err)
	}

	// Update service discovery
	if err := f.discovery.UpdatePrimary(ctx, newPrimary); err != nil {
		return fmt.Errorf("updating service discovery: %w", err)
	}

	// Promote standby to primary
	if err := f.promoteStandby(ctx, newPrimary); err != nil {
		return fmt.Errorf("promoting standby: %w", err)
	}

	f.metrics.IncCounter("successful_failovers")
	f.logger.Info("failover completed successfully", "newPrimary", newPrimary.ID)
	return nil
}

// Added implementation for selectNewPrimary
func (f *FailoverManager) selectNewPrimary(ctx context.Context) (*ServiceInstance, error) {
	// Find a healthy standby to promote
	for _, instance := range f.standby {
		if instance.IsHealthy {
			return instance, nil
		}
	}
	return nil, fmt.Errorf("no healthy standby instances available")
}

// Added implementation for promoteStandby
func (f *FailoverManager) promoteStandby(ctx context.Context, instance *ServiceInstance) error {
	// Update our records
	f.primary = instance

	// Remove the instance from standby list
	var newStandby []*ServiceInstance
	for _, s := range f.standby {
		if s.ID != instance.ID {
			newStandby = append(newStandby, s)
		}
	}
	f.standby = newStandby

	return nil
}

// Implementations of required interfaces
type simpleDiscovery struct{}

func (s *simpleDiscovery) UpdatePrimary(ctx context.Context, instance *ServiceInstance) error {
	fmt.Printf("Discovery updated with new primary: %s at %s\n", instance.ID, instance.Endpoint)
	return nil
}

type simpleHealthChecker struct {
	failureChance float64
}

func (s *simpleHealthChecker) CheckPrimary(ctx context.Context) (bool, error) {
	// Simulate occasional failures for demonstration
	if rand.Float64() < s.failureChance {
		return false, nil
	}
	return true, nil
}

type simpleMetrics struct{}

func (s *simpleMetrics) IncCounter(metric string) {
	fmt.Printf("Metric increased: %s\n", metric)
}

func (s *simpleMetrics) ObserveLatency(metric string, duration time.Duration) {
	fmt.Printf("Latency observed for %s: %v\n", metric, duration)
}

type simpleLogger struct{}

func (s *simpleLogger) Error(msg string, args ...interface{}) {
	log.Printf("ERROR: "+msg, args...)
}

func (s *simpleLogger) Info(msg string, args ...interface{}) {
	log.Printf("INFO: "+msg, args...)
}

func main() {
	// Setup instances
	primary := &ServiceInstance{
		ID:        "service-1",
		Endpoint:  "http://service-1:8080",
		IsHealthy: true,
	}

	standby := []*ServiceInstance{
		{
			ID:        "service-2",
			Endpoint:  "http://service-2:8080",
			IsHealthy: true,
		},
		{
			ID:        "service-3",
			Endpoint:  "http://service-3:8080",
			IsHealthy: true,
		},
	}

	// Create a failover manager with basic implementations
	manager := &FailoverManager{
		primary:     primary,
		standby:     standby,
		discovery:   &simpleDiscovery{},
		healthCheck: &simpleHealthChecker{failureChance: 0.3}, // 30% chance of failure for demo
		metrics:     &simpleMetrics{},
		logger:      &simpleLogger{},
	}

	// Setup context with cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Run the failover manager in a goroutine
	go func() {
		if err := manager.MonitorAndFailover(ctx); err != nil {
			log.Printf("Failover manager stopped: %v", err)
		}
	}()

	// Let it run for 1 minute for demonstration
	fmt.Println("Failover manager running. Will stop after 1 minute...")
	time.Sleep(1 * time.Minute)

	// Cancel the context to stop the manager
	cancel()
	fmt.Println("Failover manager stopped.")
}
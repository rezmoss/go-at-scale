// Example 43
package main

import (
	"fmt"
)

// Interfaces for dependencies
type Repository interface {
	GetData(id string) string
}

type Cache interface {
	Get(key string) string
	Set(key string, value string)
}

type MetricsCollector interface {
	Increment(metric string)
}

// Concrete implementations
type PostgresRepository struct{}

func (r *PostgresRepository) GetData(id string) string {
	return fmt.Sprintf("Data for ID: %s from PostgreSQL", id)
}

type RedisCache struct{}

func (c *RedisCache) Get(key string) string {
	return fmt.Sprintf("Value for key: %s from Redis", key)
}

func (c *RedisCache) Set(key string, value string) {
	fmt.Printf("Set %s = %s in Redis\n", key, value)
}

type PrometheusCollector struct{}

func (m *PrometheusCollector) Increment(metric string) {
	fmt.Printf("Incremented metric: %s in Prometheus\n", metric)
}

// Config struct
type ServiceConfig struct {
	Timeout    int
	MaxRetries int
}

// Service with dependencies
type Service struct {
	repo    Repository
	cache   Cache
	config  ServiceConfig
	metrics MetricsCollector
}

// Constructor injection
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

// Sample method that uses the dependencies
func (s *Service) GetDataWithCaching(id string) string {
	// Try cache first
	if s.cache != nil {
		cachedData := s.cache.Get(id)
		if cachedData != "" {
			if s.metrics != nil {
				s.metrics.Increment("cache_hit")
			}
			return cachedData
		}
	}

	// Get from repository
	data := s.repo.GetData(id)

	// Update cache
	if s.cache != nil {
		s.cache.Set(id, data)
	}

	if s.metrics != nil {
		s.metrics.Increment("repository_hit")
	}

	return data
}

func main() {
	// Initialize dependencies
	repo := &PostgresRepository{}
	cache := &RedisCache{}
	config := ServiceConfig{Timeout: 30, MaxRetries: 3}
	metrics := &PrometheusCollector{}

	// Example 1: Constructor injection
	service1 := NewService(repo, cache, config)
	fmt.Println("Service 1 result:", service1.GetDataWithCaching("123"))

	// Example 2: Functional options pattern
	service2 := NewServiceWithOptions(
		repo,
		WithCache(cache),
		WithMetrics(metrics),
	)
	fmt.Println("Service 2 result:", service2.GetDataWithCaching("456"))

	// Example 3: With only repository
	service3 := NewServiceWithOptions(repo)
	fmt.Println("Service 3 result:", service3.GetDataWithCaching("789"))
}
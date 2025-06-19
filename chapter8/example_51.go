// Example 51
package main

import (
	"fmt"
	"sync"
	"time"
)

// ExpensiveResource interface defines the contract
type ExpensiveResource interface {
	Request() ([]byte, error)
}

// RealResource implements the ExpensiveResource interface
type RealResource struct {
	data []byte
}

func (r *RealResource) Request() ([]byte, error) {
	// Expensive operation simulation
	time.Sleep(2 * time.Second)
	return r.data, nil
}

// CachingProxy also implements ExpensiveResource but adds caching
type CachingProxy struct {
	resource ExpensiveResource
	cache    map[string][]byte
	mu       sync.RWMutex
}

func (p *CachingProxy) Request() ([]byte, error) {
	p.mu.RLock()
	if data, ok := p.cache["key"]; ok {
		p.mu.RUnlock()
		fmt.Println("Returning cached data")
		return data, nil
	}
	p.mu.RUnlock()

	// Cache miss - get from real resource
	fmt.Println("Cache miss, getting from real resource...")
	data, err := p.resource.Request()
	if err != nil {
		return nil, err
	}
	p.mu.Lock()
	p.cache["key"] = data
	p.mu.Unlock()

	return data, nil
}

func main() {
	// Create the real resource with some data
	realResource := &RealResource{
		data: []byte("This is some expensive data"),
	}

	// Create the proxy with the real resource
	proxy := &CachingProxy{
		resource: realResource,
		cache:    make(map[string][]byte),
	}

	// First request - should go to the real resource
	fmt.Println("Making first request...")
	start := time.Now()
	data, err := proxy.Request()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("First request time: %v\n", time.Since(start))
	fmt.Printf("Data received: %s\n\n", data)

	// Second request - should be cached
	fmt.Println("Making second request...")
	start = time.Now()
	data, err = proxy.Request()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Second request time: %v\n", time.Since(start))
	fmt.Printf("Data received: %s\n", data)
}
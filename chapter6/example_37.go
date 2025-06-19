// Example 37
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Example 1: AtomicFlag
type AtomicFlag struct {
	value uint32
}

func (f *AtomicFlag) Set() bool {
	return atomic.CompareAndSwapUint32(&f.value, 0, 1)
}

func (f *AtomicFlag) Clear() {
	atomic.StoreUint32(&f.value, 0)
}

func (f *AtomicFlag) IsSet() bool {
	return atomic.LoadUint32(&f.value) == 1
}

// Example 2: Lock-free counter
type LockFreeCounter struct {
	value uint64
}

func (c *LockFreeCounter) Increment() uint64 {
	return atomic.AddUint64(&c.value, 1)
}

func (c *LockFreeCounter) Value() uint64 {
	return atomic.LoadUint64(&c.value)
}

// Example 3: Double-checked locking
type expensive struct {
	data string
}

func newExpensive() *expensive {
	// Simulate expensive initialization
	time.Sleep(100 * time.Millisecond)
	return &expensive{data: "initialized"}
}

type Singleton struct {
	initialized uint32
	instance    *expensive
	mu          sync.Mutex
}

func (s *Singleton) getInstance() *expensive {
	if atomic.LoadUint32(&s.initialized) == 0 {
		s.mu.Lock()
		defer s.mu.Unlock()

		if s.initialized == 0 {
			s.instance = newExpensive()
			atomic.StoreUint32(&s.initialized, 1)
		}
	}
	return s.instance
}

// Example 4: Using sync.Once
var once sync.Once
var instance *expensive

func GetInstance() *expensive {
	once.Do(func() {
		instance = newExpensive()
	})
	return instance
}

func main() {
	// Test AtomicFlag
	fmt.Println("Testing AtomicFlag:")
	flag := &AtomicFlag{}
	fmt.Printf("Initial state: %v\n", flag.IsSet())
	fmt.Printf("Set flag: %v\n", flag.Set())
	fmt.Printf("State after setting: %v\n", flag.IsSet())
	fmt.Printf("Set again: %v\n", flag.Set()) // Will return false since already set
	flag.Clear()
	fmt.Printf("State after clearing: %v\n", flag.IsSet())

	// Test LockFreeCounter
	fmt.Println("\nTesting LockFreeCounter:")
	counter := &LockFreeCounter{}
	var wg sync.WaitGroup

	// Increment concurrently
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				counter.Increment()
			}
		}()
	}
	wg.Wait()
	fmt.Printf("Final counter value: %d\n", counter.Value())

	// Test Singleton with double-checked locking
	fmt.Println("\nTesting Singleton with double-checked locking:")
	singleton := &Singleton{}

	var instances []*expensive
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			instances = append(instances, singleton.getInstance())
		}()
	}
	wg.Wait()

	// Check that all instances are the same
	fmt.Println("All instances are the same object:", instances[0] == instances[len(instances)-1])

	// Test sync.Once
	fmt.Println("\nTesting sync.Once:")
	var onceInstances []*expensive
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			onceInstances = append(onceInstances, GetInstance())
		}()
	}
	wg.Wait()

	// Check that all instances are the same
	fmt.Println("All sync.Once instances are the same object:", onceInstances[0] == onceInstances[len(onceInstances)-1])
}
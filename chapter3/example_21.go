// Example 21
package main

import (
	"fmt"
	"sync"
	"time"
)

type Hashable interface {
	comparable
}

func Memoize[T Hashable, U any](f func(T) U) func(T) U {
	cache := make(map[T]U)
	var mu sync.RWMutex

	return func(x T) U {
		mu.RLock()
		if val, ok := cache[x]; ok {
			mu.RUnlock()
			return val
		}
		mu.RUnlock()

		mu.Lock()
		defer mu.Unlock()

		// Double-check after acquiring write lock
		if val, ok := cache[x]; ok {
			return val
		}

		result := f(x)
		cache[x] = result
		return result
	}
}

// Example: Memoized Fibonacci
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
	// Regular fibonacci calculation
	start := time.Now()
	result := fibonacci(30)
	fmt.Printf("Regular fibonacci(30) = %d, took: %v\n", result, time.Since(start))

	// Memoized fibonacci calculation
	memoFib := Memoize(fibonacci)

	start = time.Now()
	result = memoFib(40)
	fmt.Printf("Memoized fibonacci(40) = %d, took: %v\n", result, time.Since(start))

	// Call again to demonstrate caching
	start = time.Now()
	result = memoFib(40)
	fmt.Printf("Cached memoized fibonacci(40) = %d, took: %v\n", result, time.Since(start))
}
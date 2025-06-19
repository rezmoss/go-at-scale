// Example 93
// Example with race condition
package main

import (
	"fmt"
	"sync"
)

// Example with race condition
type Counter struct {
	value int
}

func (c *Counter) Increment() {
	c.value++ // Race condition!
}

func (c *Counter) Value() int {
	return c.value
}

func main() {
	// Demonstrate race condition
	counter := Counter{}

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	wg.Wait()
	fmt.Println("Counter value:", counter.Value()) // May not be 1000

	fmt.Println("\nRun this program with -race flag to detect race conditions:")
	fmt.Println("go run -race main.go")
}
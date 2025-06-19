// Example 34
package main

import (
	"fmt"
	"sync"
	"time"
)

// Mutex-based counter
type MutexCounter struct {
	mu    sync.Mutex
	value int
}

func (c *MutexCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *MutexCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// Channel-based counter
type ChannelCounter struct {
	updates chan int
	value   chan int
	done    chan struct{}
}

func NewChannelCounter() *ChannelCounter {
	c := &ChannelCounter{
		updates: make(chan int),
		value:   make(chan int),
		done:    make(chan struct{}),
	}

	go func() {
		var current int
		for {
			select {
			case <-c.updates:
				current++
			case c.value <- current:
				// Value requested
			case <-c.done:
				close(c.value)
				return
			}
		}
	}()

	return c
}

func (c *ChannelCounter) Increment() {
	c.updates <- 1
}

func (c *ChannelCounter) Value() int {
	return <-c.value
}

func (c *ChannelCounter) Close() {
	close(c.done)
}

func main() {
	// Demonstrate mutex-based counter
	fmt.Println("Testing Mutex-based counter:")
	mutexCounter := &MutexCounter{}

	var wg sync.WaitGroup

	// Increment the mutex counter concurrently
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mutexCounter.Increment()
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Final mutex counter value: %d\n\n", mutexCounter.Value())

	// Demonstrate channel-based counter
	fmt.Println("Testing Channel-based counter:")
	channelCounter := NewChannelCounter()
	defer channelCounter.Close()

	// Reset wait group
	wg = sync.WaitGroup{}

	// Increment the channel counter concurrently
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				channelCounter.Increment()
			}
		}()
	}

	// Give some time for all increments to be processed
	wg.Wait()
	time.Sleep(100 * time.Millisecond) // Ensure all channel operations complete

	fmt.Printf("Final channel counter value: %d\n", channelCounter.Value())
}
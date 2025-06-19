// Example 52
package main

import (
	"fmt"
	"sync"
	"time"
)

// Observer interface defines the Update method that observers must implement
type Observer interface {
	Update(message string)
}

// Subject interface defines methods for managing observers
type Subject interface {
	Register(observer Observer)
	Deregister(observer Observer)
	NotifyAll(message string)
}

// EventManager implements the Subject interface
type EventManager struct {
	observers map[Observer]struct{}
	mu        sync.RWMutex
}

// NewEventManager creates and returns a new EventManager
func NewEventManager() *EventManager {
	return &EventManager{
		observers: make(map[Observer]struct{}),
	}
}

func (em *EventManager) Register(observer Observer) {
	em.mu.Lock()
	defer em.mu.Unlock()
	em.observers[observer] = struct{}{}
}

func (em *EventManager) Deregister(observer Observer) {
	em.mu.Lock()
	defer em.mu.Unlock()
	delete(em.observers, observer)
}

func (em *EventManager) NotifyAll(message string) {
	em.mu.RLock()
	defer em.mu.RUnlock()
	for observer := range em.observers {
		observer.Update(message)
	}
}

// LogObserver is a concrete implementation of the Observer interface
type LogObserver struct {
	id string
}

// Update implements the Observer interface
func (l *LogObserver) Update(message string) {
	fmt.Printf("Observer %s received: %s\n", l.id, message)
}

func main() {
	// Create a new event manager (subject)
	eventManager := NewEventManager()

	// Create observers
	observer1 := &LogObserver{id: "1"}
	observer2 := &LogObserver{id: "2"}
	observer3 := &LogObserver{id: "3"}

	// Register observers
	eventManager.Register(observer1)
	eventManager.Register(observer2)
	eventManager.Register(observer3)

	// Notify all observers
	fmt.Println("Notifying all observers:")
	eventManager.NotifyAll("Hello, observers!")

	// Deregister one observer
	fmt.Println("\nDeregistering observer 2")
	eventManager.Deregister(observer2)

	// Notify all observers again
	fmt.Println("\nNotifying remaining observers:")
	eventManager.NotifyAll("Hello again!")

	// Simulate asynchronous notification
	fmt.Println("\nSimulating asynchronous notifications:")
	go func() {
		time.Sleep(1 * time.Second)
		eventManager.NotifyAll("Async message 1")
	}()

	go func() {
		time.Sleep(2 * time.Second)
		eventManager.NotifyAll("Async message 2")
	}()

	// Wait for async notifications to complete
	time.Sleep(3 * time.Second)
}
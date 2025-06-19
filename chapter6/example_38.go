// Example 38
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// SafeResource implements resource hierarchy to prevent deadlocks
type SafeResource struct {
	id     int
	mu     sync.Mutex
	parent *SafeResource
}

func (r *SafeResource) Lock() {
	if r.parent != nil {
		r.parent.Lock()
	}
	r.mu.Lock()
}

func (r *SafeResource) Unlock() {
	r.mu.Unlock()
	if r.parent != nil {
		r.parent.Unlock()
	}
}

// TimeoutMutex demonstrates deadlock detection with timeouts
type TimeoutMutex struct {
	mu      sync.Mutex
	timeout time.Duration
}

func (t *TimeoutMutex) Lock() error {
	done := make(chan struct{})

	go func() {
		t.mu.Lock()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(t.timeout):
		return fmt.Errorf("lock acquisition timed out after %v", t.timeout)
	}
}

func (t *TimeoutMutex) Unlock() {
	t.mu.Unlock()
}

// TryMutex implements try-lock pattern
type TryMutex struct {
	locked uint32
}

func (t *TryMutex) TryLock() bool {
	return atomic.CompareAndSwapUint32(&t.locked, 0, 1)
}

func (t *TryMutex) Unlock() {
	atomic.StoreUint32(&t.locked, 0)
}

func main() {
	// Example 1: Resource Hierarchy
	fmt.Println("Example 1: Resource Hierarchy")
	parent := &SafeResource{id: 1}
	child := &SafeResource{id: 2, parent: parent}

	// Safe locking order
	child.Lock()
	fmt.Println("Acquired locks in hierarchy")
	child.Unlock()

	// Example 2: TimeoutMutex
	fmt.Println("\nExample 2: TimeoutMutex")
	tm := &TimeoutMutex{timeout: 500 * time.Millisecond}

	// Normal lock acquisition
	err := tm.Lock()
	if err == nil {
		fmt.Println("Successfully acquired lock with timeout")
		tm.Unlock()
	} else {
		fmt.Println("Error:", err)
	}

	// Simulate deadlock
	tm2 := &TimeoutMutex{timeout: 500 * time.Millisecond}
	tm2.mu.Lock() // Lock it first to create contention

	go func() {
		time.Sleep(1 * time.Second)
		tm2.mu.Unlock() // Unlock after timeout should occur
	}()

	err = tm2.Lock()
	if err != nil {
		fmt.Println("Expected timeout error:", err)
	} else {
		fmt.Println("Lock acquired (unexpected)")
		tm2.Unlock()
	}

	// Example 3: TryMutex
	fmt.Println("\nExample 3: TryMutex")
	tryMu := &TryMutex{}

	// First attempt should succeed
	if tryMu.TryLock() {
		fmt.Println("First try-lock succeeded")
	} else {
		fmt.Println("First try-lock failed (unexpected)")
	}

	// Second attempt should fail as it's already locked
	if tryMu.TryLock() {
		fmt.Println("Second try-lock succeeded (unexpected)")
	} else {
		fmt.Println("Second try-lock failed as expected")
	}

	// Unlock and try again
	tryMu.Unlock()
	if tryMu.TryLock() {
		fmt.Println("Try-lock after unlock succeeded")
		tryMu.Unlock()
	} else {
		fmt.Println("Try-lock after unlock failed (unexpected)")
	}
}
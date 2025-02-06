// Example 38
// Resource hierarchy to prevent deadlocks
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

// Deadlock detection with timeouts
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

// Example: Try-lock pattern
type TryMutex struct {
    locked uint32
}

func (t *TryMutex) TryLock() bool {
    return atomic.CompareAndSwapUint32(&t.locked, 0, 1)
}

func (t *TryMutex) Unlock() {
    atomic.StoreUint32(&t.locked, 0)
}
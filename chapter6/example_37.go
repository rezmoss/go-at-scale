// Example 37
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

// Example: Lock-free counter
type LockFreeCounter struct {
    value uint64
}

func (c *LockFreeCounter) Increment() uint64 {
    return atomic.AddUint64(&c.value, 1)
}

func (c *LockFreeCounter) Value() uint64 {
    return atomic.LoadUint64(&c.value)
}

// Example: Double-checked locking
type Singleton struct {
    initialized uint32
    instance    *expensive
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



//sync.Once

var once sync.Once
var instance *expensive

func GetInstance() *expensive {
    once.Do(func() {
        instance = newExpensive()
    })
    return instance
}

// Explanation:
// sync.Once automatically ensures that the initialization code
// runs exactly once in a goroutine-safe manner. This is simpler
// and more idiomatic than double-checked locking in Go.
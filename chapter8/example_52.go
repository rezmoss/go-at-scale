// Example 52
type Observer interface {
    Update(message string)
}

type Subject interface {
    Register(observer Observer)
    Deregister(observer Observer)
    NotifyAll(message string)
}

type EventManager struct {
    observers map[Observer]struct{}
    mu        sync.RWMutex
}

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
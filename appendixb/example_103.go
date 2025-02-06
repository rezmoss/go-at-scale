// Example 103
// Pitfall 1: Concurrent map access
type UserCache struct {
    cache map[string]User
}

// Race condition!
func (uc *UserCache) GetUser(id string) User {
    return uc.cache[id]
}

// Solution: Use sync.Map or mutex
type SafeUserCache struct {
    cache map[string]User
    mu    sync.RWMutex
}

func (sc *SafeUserCache) GetUser(id string) User {
    sc.mu.RLock()
    defer sc.mu.RUnlock()
    return sc.cache[id]
}

// Pitfall 2: Memory leaks in long-lived maps
type Session struct {
    Data []byte  // Potentially large
}

var sessions = make(map[string]*Session)  // Never cleaned up!

// Solution: Implement cleanup
type SessionManager struct {
    sessions map[string]*Session
    mu       sync.RWMutex
}

func (sm *SessionManager) CleanupExpired() {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    for id, session := range sm.sessions {
        if session.isExpired() {
            delete(sm.sessions, id)
        }
    }
}
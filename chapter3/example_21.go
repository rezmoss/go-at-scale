// Example 21
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
    memoFib := Memoize(fibonacci)
    fmt.Println(memoFib(40))  // Much faster than unmemoized version
}
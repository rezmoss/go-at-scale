// Example 107
// Pitfall 1: Interface pollution
type DoEverything interface {
    DoA()
    DoB()
    DoC()
    // ... many more methods
}

// Solution: Small, focused interfaces
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// Pitfall 2: Interface{} abuse
func processPitfall(data interface{}) {
    // Type assertions everywhere
    switch v := data.(type) {
    case string:
        // handle string
    case int:
        // handle int
    // ... many cases
    }
}

// Solution: Use generics when possible
func process[T string | int](data T) {
    // Type-safe processing
}
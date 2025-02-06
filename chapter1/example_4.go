// Example 4
func createCounter() func() int {
    count := 0  // This variable is "closed over"
    return func() int {
        count++
        return count
    }
}

// Usage
counter := createCounter()
fmt.Println(counter())  // 1
fmt.Println(counter())  // 2
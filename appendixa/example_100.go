// Example 100
func main() {
    f, err := os.Create("trace.out")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    
    trace.Start(f)
    defer trace.Stop()
    
    // Your program here
}

// Analyze trace
// go tool trace trace.out
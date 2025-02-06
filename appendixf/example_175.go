// Example 175
// Anti-pattern: Unnecessary allocations
func processItems(items []string) []string {
    result := []string{}  // Grows with append
    for _, item := range items {
        result = append(result, process(item))
    }
    return result
}

// Proper pattern: Pre-allocate slice
func processItems(items []string) []string {
    result := make([]string, 0, len(items))
    for _, item := range items {
        result = append(result, process(item))
    }
    return result
}
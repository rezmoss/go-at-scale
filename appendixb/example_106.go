// Example 106
// Pitfall 1: String concatenation in loops
func buildStringPitfall(items []string) string {
    result := ""
    for _, item := range items {
        result += item  // Creates new string each time
    }
    return result
}

// Solution: Use strings.Builder
func buildStringSolution(items []string) string {
    var builder strings.Builder
    builder.Grow(len(items) * 8)  // Estimate capacity
    for _, item := range items {
        builder.WriteString(item)
    }
    return builder.String()
}

// Pitfall 2: Converting strings to bytes repeatedly
func processStringPitfall(s string) {
    for i := 0; i < len(s); i++ {
        b := []byte(s)  // Converts entire string each time
        process(b[i])
    }
}

// Solution: Convert once
func processStringSolution(s string) {
    b := []byte(s)
    for i := 0; i < len(b); i++ {
        process(b[i])
    }
}
// Example 176
// Anti-pattern: Inefficient string building
func buildReport(items []string) string {
    result := ""
    for _, item := range items {
        result += item + "\n"  // Creates new string each iteration
    }
    return result
}

// Proper pattern: Use strings.Builder
func buildReport(items []string) string {
    var builder strings.Builder
    builder.Grow(len(items) * 20)  // Estimate size
    for _, item := range items {
        builder.WriteString(item)
        builder.WriteByte('\n')
    }
    return builder.String()
}
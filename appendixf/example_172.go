// Example 172
// Anti-pattern: Using panic for error handling
func processConfig(path string) Config {
    data, err := readConfig(path)
    if err != nil {
        panic(err)  // Wrong: Using panic for normal errors
    }
    return parseConfig(data)
}

// Proper pattern: Return errors normally
func processConfig(path string) (Config, error) {
    data, err := readConfig(path)
    if err != nil {
        return Config{}, fmt.Errorf("reading config: %w", err)
    }
    return parseConfig(data)
}
// Example 174
// Anti-pattern: Resource leaks
func processFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    // No defer close, resource leak
    data, err := ioutil.ReadAll(file)
    if err != nil {
        return err
    }
    return processData(data)
}

// Proper pattern: Ensure cleanup
func processFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("opening file: %w", err)
    }
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if err != nil {
        return fmt.Errorf("reading file: %w", err)
    }
    return processData(data)
}
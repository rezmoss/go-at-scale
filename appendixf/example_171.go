// Example 171

// Anti-pattern: Losing error context
func processUser(id string) error {
    user, err := getUser(id)
    if err != nil {
        return err  // Context lost
    }

    return nil

}

// Proper pattern: Preserve error context

func processUser(id string) error {
    user, err := getUser(id)
    if err != nil {
        return fmt.Errorf("processing user %s: %w", id, err)
    }

    return nil

}
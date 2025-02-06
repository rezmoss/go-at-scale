// Example 5
// Incorrect
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i)  // Might print unexpected values
    }()
}

// Correct
for i := 0; i < 3; i++ {
    i := i  // Create new variable for closure
    go func() {
        fmt.Println(i)  // Prints expected values
    }()
}
// Example 108
// Pitfall 1: Not cleaning up test resources
func TestDatabase(t *testing.T) {
    db, _ := sql.Open("postgres", "connection-string")
    // No cleanup!
    runTest(db)
}

// Solution: Use cleanup function
func TestDatabaseCleanup(t *testing.T) {
    db, err := sql.Open("postgres", "connection-string")
    if err != nil {
        t.Fatal(err)
    }
    t.Cleanup(func() {
        db.Close()
    })
    runTest(db)
}

// Pitfall 2: Time-dependent tests
func TestTimeout(t *testing.T) {
    start := time.Now()
    time.Sleep(1 * time.Second)
    if time.Since(start) > time.Second {
        t.Error("took too long")
    }
}

// Solution: Use time.Timer or testing.Timer
func TestTimeoutSolution(t *testing.T) {
    timer := time.NewTimer(2 * time.Second)
    defer timer.Stop()
    
    done := make(chan struct{})
    go func() {
        operation()
        close(done)
    }()
    
    select {
    case <-done:
        // Success
    case <-timer.C:
        t.Error("timeout")
    }
}
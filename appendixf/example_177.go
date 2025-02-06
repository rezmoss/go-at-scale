// Example 177
// Anti-pattern: Monolithic tests
func TestUser(t *testing.T) {
    // Setup
    db := setupDB()
    user := createTestUser() 
    // Test everything in one function
    t.Run("test all user operations", func(t *testing.T) {
        // Test create
        err := db.CreateUser(user)
        if err != nil {
            t.Error(err)
        }
        // Test get
        retrieved, err := db.GetUser(user.ID)
        if err != nil {
            t.Error(err)
        }
        // Test update
        // Test delete
        // ... many more operations
    })
}

// Proper pattern: Focused test cases
func TestUser(t *testing.T) {
    tests := []struct {
        name string
        fn   func(*testing.T)
    }{
        {
            name: "create user",
            fn:   testCreateUser,
        },
        {
            name: "get user",
            fn:   testGetUser,
        },
        {
            name: "update user",
            fn:   testUpdateUser,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, tt.fn)
    }
}
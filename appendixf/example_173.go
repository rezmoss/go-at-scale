// Example 173
// Anti-pattern: Interface bloat
type Service interface {
    GetUser(id string) (*User, error)
    CreateUser(user *User) error
    UpdateUser(user *User) error
    DeleteUser(id string) error
    ListUsers() ([]*User, error)
    ValidateUser(user *User) error
    NotifyUser(id string, message string) error
    // ... many more methods
}

// Proper pattern: Interface segregation
type UserReader interface {
    GetUser(id string) (*User, error)
    ListUsers() ([]*User, error)
}

type UserWriter interface {
    CreateUser(user *User) error
    UpdateUser(user *User) error
    DeleteUser(id string) error
}

type UserNotifier interface {
    NotifyUser(id string, message string) error
}
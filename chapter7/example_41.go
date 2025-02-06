// Example 41
// Bad: Kitchen sink interface
type Service interface {
    CreateUser(ctx context.Context, user *User) error
    UpdateUser(ctx context.Context, user *User) error
    DeleteUser(ctx context.Context, id string) error
    GetUser(ctx context.Context, id string) (*User, error)
    ListUsers(ctx context.Context) ([]*User, error)
    ValidateUser(user *User) error
    NotifyUser(ctx context.Context, userID, message string) error
    // ... many more methods
}

// Good: Interface segregation
type UserReader interface {
    GetUser(ctx context.Context, id string) (*User, error)
    ListUsers(ctx context.Context) ([]*User, error)
}

type UserWriter interface {
    CreateUser(ctx context.Context, user *User) error
    UpdateUser(ctx context.Context, user *User) error
    DeleteUser(ctx context.Context, id string) error
}

type NotificationService interface {
    NotifyUser(ctx context.Context, userID, message string) error
}

// Composite service using small interfaces
type UserService struct {
    reader UserReader
    writer UserWriter
    notifier NotificationService
}
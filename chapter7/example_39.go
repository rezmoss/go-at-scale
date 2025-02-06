// Example 39
// project/
//   ├── cmd/
//   │   └── server/
//   │       └── main.go
//   ├── internal/
//   │   ├── auth/
//   │   │   ├── middleware.go
//   │   │   └── service.go
//   │   └── storage/
//   │       ├── memory.go
//   │       └── postgres.go
//   ├── pkg/
//   │   └── validator/
//   │       └── validator.go
//   └── api/
//   │   └── v1/
//   │       └── types.go
//   ├── integration/
//   │   ├── setup/
//   │   │   └── testcontainers.go
//   │   ├── api/
//   │   │   └── api_test.go
//   │   └── db/
//   │       └── db_test.go
//   └── test/
//       └── mocks/
// Example: Well-structured package
package user

// domain.go - Domain types
type User struct {
    ID        string
    Email     string
    CreatedAt time.Time
}

// repository.go - Data access interface
type Repository interface {
    Find(ctx context.Context, id string) (*User, error)
    Save(ctx context.Context, user *User) error
}

// service.go - Business logic
type Service struct {
    repo Repository
    auth Authenticator
}

// handler.go - HTTP handlers
type Handler struct {
    service *Service
    logger  Logger
}
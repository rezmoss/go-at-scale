// Example 42
// Generic repository pattern
type Reader[T any] interface {
    Find(ctx context.Context, id string) (T, error)
    List(ctx context.Context) ([]T, error)
}

type Writer[T any] interface {
    Create(ctx context.Context, item T) error
    Update(ctx context.Context, item T) error
    Delete(ctx context.Context, id string) error
}

type Repository[T any] interface {
    Reader[T]
    Writer[T]
}

// Implementation example
type PostgresRepository[T any] struct {
    db *sql.DB
}

func (r *PostgresRepository[T]) Find(ctx context.Context, id string) (T, error) {
    // Implementation
}

// Usage
type UserRepository interface {
    Repository[User]
}
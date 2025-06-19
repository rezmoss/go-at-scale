// Example 42
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

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

// User model for demonstration
type User struct {
	ID   string
	Name string
}

// Implementation example
type PostgresRepository[T any] struct {
	db *sql.DB
}

func (r *PostgresRepository[T]) Find(ctx context.Context, id string) (T, error) {
	var empty T
	// Implementation would normally query the database
	// For example purposes, we'll just return an empty value
	fmt.Printf("Finding item with ID: %s\n", id)
	return empty, nil
}

func (r *PostgresRepository[T]) List(ctx context.Context) ([]T, error) {
	// Implementation would normally query the database
	fmt.Println("Listing all items")
	return []T{}, nil
}

func (r *PostgresRepository[T]) Create(ctx context.Context, item T) error {
	// Implementation would normally insert into the database
	fmt.Println("Creating new item")
	return nil
}

func (r *PostgresRepository[T]) Update(ctx context.Context, item T) error {
	// Implementation would normally update the database
	fmt.Println("Updating existing item")
	return nil
}

func (r *PostgresRepository[T]) Delete(ctx context.Context, id string) error {
	// Implementation would normally delete from the database
	fmt.Printf("Deleting item with ID: %s\n", id)
	return nil
}

// Usage - defining a repository for User type
type UserRepository interface {
	Repository[User]
}

// Concrete implementation of UserRepository
type PostgresUserRepository struct {
	PostgresRepository[User]
}

// NewPostgresUserRepository creates a new repository with DB connection
func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		PostgresRepository: PostgresRepository[User]{db: db},
	}
}

func main() {
	// Set up a mock DB connection (in a real app, you'd use a real connection)
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create a repository
	userRepo := NewPostgresUserRepository(db)

	// Create a context
	ctx := context.Background()

	// Use the repository
	user := User{ID: "1", Name: "John Doe"}

	// Demonstrate the interface methods
	err = userRepo.Create(ctx, user)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	foundUser, err := userRepo.Find(ctx, "1")
	if err != nil {
		log.Fatalf("Failed to find user: %v", err)
	}
	fmt.Printf("Found user: %+v\n", foundUser)

	users, err := userRepo.List(ctx)
	if err != nil {
		log.Fatalf("Failed to list users: %v", err)
	}
	fmt.Printf("Found %d users\n", len(users))

	err = userRepo.Update(ctx, user)
	if err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}

	err = userRepo.Delete(ctx, "1")
	if err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}

	fmt.Println("All operations completed successfully")
}
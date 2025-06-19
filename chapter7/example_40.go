// Example 40
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"time"
)

func main() {
	fmt.Println("Example 40: Package-Level Guidelines")

	// Demonstrate the bad example
	fmt.Println("\nBad Example - Too many dependencies:")
	badUserService := createBadUserService()
	fmt.Printf("Bad UserService has direct dependencies: %v\n", badUserService != nil)

	// Demonstrate the good example
	fmt.Println("\nGood Example - Focused functionality:")
	goodUserService := createGoodUserService()
	fmt.Printf("Good UserService uses abstractions: %v\n", goodUserService != nil)
}

// Bad: Too many dependencies would be in package user
type BadUserService struct {
	db        *sql.DB
	cache     interface{} // Replacing redis.Client
	templates *template.Template
}

func createBadUserService() *BadUserService {
	// Just for demonstration, not actually connecting
	return &BadUserService{
		db:        nil, // Would be sql.DB in real code
		cache:     nil, // Would be redis.Client in real code
		templates: nil, // Would be template.Template in real code
	}
}

// Good: Clean abstraction would be in package user
type Repository interface {
	FindUser(ctx context.Context, id string) (interface{}, error)
}

type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, expiration time.Duration) error
}

type EventEmitter interface {
	Emit(event string, payload interface{}) error
}

type GoodUserService struct {
	store  Repository
	cache  Cache
	events EventEmitter
}

func createGoodUserService() *GoodUserService {
	return &GoodUserService{
		store:  &mockRepository{},
		cache:  &mockCache{},
		events: &mockEventEmitter{},
	}
}

// Mock implementations for interfaces
type mockRepository struct{}

func (m *mockRepository) FindUser(ctx context.Context, id string) (interface{}, error) {
	return nil, errors.New("not implemented")
}

type mockCache struct{}

func (m *mockCache) Get(key string) (interface{}, error) {
	return nil, errors.New("not implemented")
}
func (m *mockCache) Set(key string, value interface{}, expiration time.Duration) error {
	return errors.New("not implemented")
}

type mockEventEmitter struct{}

func (m *mockEventEmitter) Emit(event string, payload interface{}) error {
	return errors.New("not implemented")
}
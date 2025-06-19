// Example 63
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/graph-gophers/dataloader"
)

// Repository interface for data access
type Repository interface {
	GetUsersByIDs(ctx context.Context, ids []string) (map[string]User, error)
}

// User represents a user entity
type User struct {
	ID   string
	Name string
}

// Logger interface for logging
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// SimpleLogger is a basic implementation of Logger
type SimpleLogger struct{}

func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	log.Printf("[INFO] "+msg, args...)
}

func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	log.Printf("[ERROR] "+msg, args...)
}

type DataLoader struct {
	userLoader *dataloader.Loader
	postLoader *dataloader.Loader
	redis      *redis.Client
	logger     Logger
}

func NewDataLoader(ctx context.Context, repo Repository) *DataLoader {
	return &DataLoader{
		userLoader: dataloader.NewBatchedLoader(dataloader.BatchFunc(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
			// Convert dataloader.Keys to []string
			stringKeys := make([]string, 0, len(keys))
			for _, key := range keys {
				stringKeys = append(stringKeys, key.String())
			}

			// Batch load users
			users, err := repo.GetUsersByIDs(ctx, stringKeys)
			if err != nil {
				return makeBatchError(err, len(stringKeys))
			}

			// Map results to keys order
			return mapResultsToKeys(stringKeys, users)
		})),
		// ... other loaders
		logger: &SimpleLogger{},
	}
}

// makeBatchError creates error results for all keys
func makeBatchError(err error, count int) []*dataloader.Result {
	results := make([]*dataloader.Result, count)
	for i := 0; i < count; i++ {
		results[i] = &dataloader.Result{Error: err}
	}
	return results
}

// mapResultsToKeys maps results to the original key order
func mapResultsToKeys(keys []string, users map[string]User) []*dataloader.Result {
	results := make([]*dataloader.Result, len(keys))
	for i, key := range keys {
		user, exists := users[key]
		if !exists {
			results[i] = &dataloader.Result{Error: fmt.Errorf("user not found: %s", key)}
			continue
		}
		results[i] = &dataloader.Result{Data: user}
	}
	return results
}

// MockRepository implements Repository for testing
type MockRepository struct{}

func (r *MockRepository) GetUsersByIDs(ctx context.Context, ids []string) (map[string]User, error) {
	result := make(map[string]User)
	for _, id := range ids {
		// Create mock users
		result[id] = User{ID: id, Name: "User " + id}
	}
	return result, nil
}

func main() {
	ctx := context.Background()
	repo := &MockRepository{}

	// Create a new DataLoader
	loader := NewDataLoader(ctx, repo)

	thunk := loader.userLoader.Load(ctx, dataloader.StringKey("user1"))
	result, err := thunk()

	if err != nil {
		log.Fatalf("Error loading user: %v", err)
	}

	user := result.(User)
	fmt.Printf("Loaded user: ID=%s, Name=%s\n", user.ID, user.Name)
}
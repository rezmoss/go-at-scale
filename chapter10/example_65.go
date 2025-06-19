// Example 65
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// User represents a user in the system
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Resolver is the interface for resolving GraphQL queries
type Resolver interface {
	Query_user(ctx context.Context, id string) (*User, error)
}

// SimpleResolver is a basic implementation of the Resolver interface
type SimpleResolver struct{}

func (r *SimpleResolver) Query_user(ctx context.Context, id string) (*User, error) {
	// Simulate database lookup
	log.Printf("Getting user %s from database", id)
	return &User{
		ID:    id,
		Name:  fmt.Sprintf("User %s", id),
		Email: fmt.Sprintf("user%s@example.com", id),
	}, nil
}

// CacheResolver struct as per the example
type CacheResolver struct {
	underlying Resolver
	cache      *redis.Client
	ttl        time.Duration
}

func (r *CacheResolver) Query_user(ctx context.Context, id string) (*User, error) {
	cacheKey := fmt.Sprintf("user:%s", id)

	// Try cache first
	if cached, err := r.getFromCache(ctx, cacheKey); err == nil {
		log.Printf("Cache hit for user %s", id)
		return cached, nil
	}

	// Get from underlying resolver
	user, err := r.underlying.Query_user(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache result
	err = r.cacheResult(ctx, cacheKey, user)
	if err != nil {
		log.Printf("Failed to cache user: %v", err)
	}

	return user, nil
}

// getFromCache retrieves a user from the cache
func (r *CacheResolver) getFromCache(ctx context.Context, key string) (*User, error) {
	log.Printf("Attempting to get %s from cache", key)
	data, err := r.cache.Get(ctx, key).Bytes()
	if err != nil {
		log.Printf("Cache miss for %s: %v", key, err)
		return nil, err
	}

	log.Printf("Cache hit for %s", key)
	var user User
	if err := json.Unmarshal(data, &user); err != nil {
		log.Printf("Failed to unmarshal cached user: %v", err)
		return nil, err
	}

	return &user, nil
}

// cacheResult stores a user in the cache
func (r *CacheResolver) cacheResult(ctx context.Context, key string, user *User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	err = r.cache.Set(ctx, key, data, r.ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	log.Printf("Successfully cached user %s for %v", user.ID, r.ttl)
	return nil
}

func main() {
	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test connection
	ctx := context.Background()
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Printf("Redis connection successful: %s", pong)

	// Create resolvers
	simpleResolver := &SimpleResolver{}
	cacheResolver := &CacheResolver{
		underlying: simpleResolver,
		cache:      rdb,
		ttl:        time.Minute * 5,
	}

	// Test query
	// First query - should miss cache
	user1, err := cacheResolver.Query_user(ctx, "123")
	if err != nil {
		log.Fatalf("Failed to resolve user: %v", err)
	}
	log.Printf("First query result: %+v", user1)

	// Second query - should hit cache
	user2, err := cacheResolver.Query_user(ctx, "123")
	if err != nil {
		log.Fatalf("Failed to resolve user: %v", err)
	}
	log.Printf("Second query result: %+v", user2)

	// Different user - should miss cache
	user3, err := cacheResolver.Query_user(ctx, "456")
	if err != nil {
		log.Fatalf("Failed to resolve user: %v", err)
	}
	log.Printf("Different user query result: %+v", user3)

	log.Println("Example completed successfully")
}
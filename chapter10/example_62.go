// Example 62
package main

import (
	"context"
	"fmt"
	"time"
)

// User represents a user entity
type User struct {
	ID    string
	Name  string
	Email string
}

// Post represents a post entity
type Post struct {
	ID     string
	Title  string
	Body   string
	UserID string
}

// UserService provides user-related operations
type UserService interface {
	GetUser(ctx context.Context, id string) (*User, error)
}

// PostService provides post-related operations
type PostService interface {
	GetUserPosts(ctx context.Context, userID string) ([]*Post, error)
}

// Logger provides logging functionality
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// MetricsRecorder provides metrics recording functionality
type MetricsRecorder interface {
	ObserveLatency(metric string, duration time.Duration)
}

// DataLoader provides data loading with batching capability
type DataLoader struct {
	userService UserService
}

func NewDataLoader(us UserService) *DataLoader {
	return &DataLoader{userService: us}
}

func (dl *DataLoader) LoadUser(ctx context.Context, id string) (*User, error) {
	// In a real implementation, this would batch multiple requests
	return dl.userService.GetUser(ctx, id)
}

// Resolvers struct as shown in the original example
type Resolvers struct {
	userService UserService
	postService PostService
	dataLoader  *DataLoader
	metrics     MetricsRecorder
	logger      Logger
}

func (r *Resolvers) User(ctx context.Context, id string) (*User, error) {
	start := time.Now()
	defer func() {
		r.metrics.ObserveLatency("resolver.user", time.Since(start))
	}()

	// Use dataloader for batching
	return r.dataLoader.LoadUser(ctx, id)
}

// Field resolver pattern
func (r *Resolvers) User_posts(ctx context.Context, obj *User) ([]*Post, error) {
	return r.postService.GetUserPosts(ctx, obj.ID)
}

// Implementations for the interfaces

// SimpleUserService implements UserService
type SimpleUserService struct{}

func (s *SimpleUserService) GetUser(ctx context.Context, id string) (*User, error) {
	// Simulate database lookup
	user := &User{
		ID:    id,
		Name:  "User " + id,
		Email: "user" + id + "@example.com",
	}
	return user, nil
}

// SimplePostService implements PostService
type SimplePostService struct{}

func (s *SimplePostService) GetUserPosts(ctx context.Context, userID string) ([]*Post, error) {
	// Simulate fetching posts for a user
	posts := []*Post{
		{
			ID:     "post1",
			Title:  "First Post by User " + userID,
			Body:   "This is the first post content",
			UserID: userID,
		},
		{
			ID:     "post2",
			Title:  "Second Post by User " + userID,
			Body:   "This is the second post content",
			UserID: userID,
		},
	}
	return posts, nil
}

// SimpleMetricsRecorder implements MetricsRecorder
type SimpleMetricsRecorder struct{}

func (s *SimpleMetricsRecorder) ObserveLatency(metric string, duration time.Duration) {
	fmt.Printf("Metric: %s, Duration: %v\n", metric, duration)
}

// SimpleLogger implements Logger
type SimpleLogger struct{}

func (s *SimpleLogger) Info(msg string, args ...interface{}) {
	fmt.Printf("INFO: "+msg+"\n", args...)
}

func (s *SimpleLogger) Error(msg string, args ...interface{}) {
	fmt.Printf("ERROR: "+msg+"\n", args...)
}

func main() {
	// Initialize services
	userService := &SimpleUserService{}
	postService := &SimplePostService{}
	metrics := &SimpleMetricsRecorder{}
	logger := &SimpleLogger{}
	dataLoader := NewDataLoader(userService)

	// Initialize resolvers
	resolvers := &Resolvers{
		userService: userService,
		postService: postService,
		dataLoader:  dataLoader,
		metrics:     metrics,
		logger:      logger,
	}

	// Test the resolvers
	ctx := context.Background()

	// Test User resolver
	logger.Info("Testing User resolver")
	user, err := resolvers.User(ctx, "123")
	if err != nil {
		logger.Error("Failed to resolve user: %v", err)
	} else {
		logger.Info("Resolved user: ID=%s, Name=%s", user.ID, user.Name)
	}

	// Test User_posts resolver
	if user != nil {
		logger.Info("Testing User_posts resolver")
		posts, err := resolvers.User_posts(ctx, user)
		if err != nil {
			logger.Error("Failed to resolve user posts: %v", err)
		} else {
			logger.Info("Resolved %d posts for user %s", len(posts), user.ID)
			for i, post := range posts {
				logger.Info("Post %d: ID=%s, Title=%s", i+1, post.ID, post.Title)
			}
		}
	}
}
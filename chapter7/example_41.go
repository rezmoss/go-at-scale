// Example 41
package main

import (
	"context"
	"fmt"
)

// User represents a user in the system
type User struct {
	ID    string
	Name  string
	Email string
}

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
	reader   UserReader
	writer   UserWriter
	notifier NotificationService
}

// Implementation of the interfaces for demonstration
type InMemoryUserStore struct {
	users map[string]*User
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: make(map[string]*User),
	}
}

// Implement UserReader
func (s *InMemoryUserStore) GetUser(ctx context.Context, id string) (*User, error) {
	user, exists := s.users[id]
	if !exists {
		return nil, fmt.Errorf("user with ID %s not found", id)
	}
	return user, nil
}

func (s *InMemoryUserStore) ListUsers(ctx context.Context) ([]*User, error) {
	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users, nil
}

// Implement UserWriter
func (s *InMemoryUserStore) CreateUser(ctx context.Context, user *User) error {
	if _, exists := s.users[user.ID]; exists {
		return fmt.Errorf("user with ID %s already exists", user.ID)
	}
	s.users[user.ID] = user
	return nil
}

func (s *InMemoryUserStore) UpdateUser(ctx context.Context, user *User) error {
	if _, exists := s.users[user.ID]; !exists {
		return fmt.Errorf("user with ID %s not found", user.ID)
	}
	s.users[user.ID] = user
	return nil
}

func (s *InMemoryUserStore) DeleteUser(ctx context.Context, id string) error {
	if _, exists := s.users[id]; !exists {
		return fmt.Errorf("user with ID %s not found", id)
	}
	delete(s.users, id)
	return nil
}

// Simple notification service implementation
type ConsoleNotifier struct{}

func (n *ConsoleNotifier) NotifyUser(ctx context.Context, userID, message string) error {
	fmt.Printf("Notification to user %s: %s\n", userID, message)
	return nil
}

// Example usage of the segregated interfaces
func main() {
	// Create implementations
	store := NewInMemoryUserStore()
	notifier := &ConsoleNotifier{}

	// Create composite service
	service := &UserService{
		reader:   store,
		writer:   store,
		notifier: notifier,
	}

	// Example usage of the service
	ctx := context.Background()

	// Create a user
	user := &User{
		ID:    "user1",
		Name:  "John Doe",
		Email: "john@example.com",
	}

	// Use the writer interface through the composite service
	err := service.writer.CreateUser(ctx, user)
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		return
	}

	// Use the reader interface
	retrievedUser, err := service.reader.GetUser(ctx, "user1")
	if err != nil {
		fmt.Printf("Error getting user: %v\n", err)
		return
	}
	fmt.Printf("Retrieved user: %s (%s)\n", retrievedUser.Name, retrievedUser.Email)

	// Use the notifier interface
	err = service.notifier.NotifyUser(ctx, "user1", "Welcome to our service!")
	if err != nil {
		fmt.Printf("Error sending notification: %v\n", err)
		return
	}

	// List all users
	users, err := service.reader.ListUsers(ctx)
	if err != nil {
		fmt.Printf("Error listing users: %v\n", err)
		return
	}
	fmt.Printf("Total users: %d\n", len(users))
}
// Example 39
package user

import (
	"context"
	"errors"
	"time"
)

// Domain types
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

// Simple memory implementation of Repository
type MemoryRepository struct {
	users map[string]*User
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users: make(map[string]*User),
	}
}

func (r *MemoryRepository) Find(ctx context.Context, id string) (*User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *MemoryRepository) Save(ctx context.Context, user *User) error {
	r.users[user.ID] = user
	return nil
}

// Logger interface
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// Simple logger implementation
type SimpleLogger struct{}

func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	// Implementation omitted for brevity
}

func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	// Implementation omitted for brevity
}

// Authenticator interface
type Authenticator interface {
	Authenticate(token string) (string, error)
}

// Simple authenticator implementation
type SimpleAuthenticator struct{}

func (a *SimpleAuthenticator) Authenticate(token string) (string, error) {
	// Implementation omitted for brevity
	return "user-id", nil
}

// service.go - Business logic
type Service struct {
	repo Repository
	auth Authenticator
}

func NewService(repo Repository, auth Authenticator) *Service {
	return &Service{
		repo: repo,
		auth: auth,
	}
}

func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
	return s.repo.Find(ctx, id)
}

func (s *Service) CreateUser(ctx context.Context, email string) (*User, error) {
	user := &User{
		ID:        "user-" + email, // Simple ID generation
		Email:     email,
		CreatedAt: time.Now(),
	}

	err := s.repo.Save(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// handler.go - HTTP handlers
type Handler struct {
	service *Service
	logger  Logger
}

func NewHandler(service *Service, logger Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Example of a handler method
func (h *Handler) GetUserHandler(ctx context.Context, userID string) (*User, error) {
	h.logger.Info("Getting user", "id", userID)

	user, err := h.service.GetUser(ctx, userID)
	if err != nil {
		h.logger.Error("Failed to get user", "error", err)
		return nil, err
	}

	return user, nil
}
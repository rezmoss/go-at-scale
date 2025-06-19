// Example 45
package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
)

// User represents a user in the system
type User struct {
	ID       string
	Email    string
	Password string
}

// Repository interface for user data storage
type Repository interface {
	Find(ctx context.Context, id string) (*User, error)
	Save(ctx context.Context, user *User) error
}

// PasswordHasher handles password hashing
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(hashed, plain string) bool
}

// EmailSender handles sending emails
type EmailSender interface {
	Send(to, subject, body string) error
}

// UserService handles user-related operations
type UserService struct {
	repo   Repository
	hasher PasswordHasher
	mailer EmailSender
}

// NewUserService creates a new user service
func NewUserService(repo Repository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user User) error {
	return s.repo.Save(ctx, &user)
}

// Test doubles
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Find(ctx context.Context, id string) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) Save(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// Table-driven tests
func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name    string
		input   User
		mockFn  func(*MockRepository)
		wantErr bool
	}{
		{
			name:  "successful creation",
			input: User{Email: "test@example.com"},
			mockFn: func(repo *MockRepository) {
				repo.On("Save", mock.Anything, mock.AnythingOfType("*User")).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "repository error",
			input: User{Email: "test@example.com"},
			mockFn: func(repo *MockRepository) {
				repo.On("Save", mock.Anything, mock.AnythingOfType("*User")).
					Return(errDBConnection)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockRepository)
			tt.mockFn(repo)

			service := NewUserService(repo)
			err := service.CreateUser(context.Background(), tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			repo.AssertExpectations(t)
		})
	}
}

// Define a custom error for testing
var errDBConnection = ErrorDBConnection("failed to connect to database")

type ErrorDBConnection string

func (e ErrorDBConnection) Error() string {
	return string(e)
}

func main() {

}
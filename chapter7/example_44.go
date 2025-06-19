// Example 44
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// Domain-specific errors
type ErrorCode int

const (
	ErrNotFound ErrorCode = iota + 1
	ErrInvalidInput
	ErrUnauthorized
)

type DomainError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

// ErrorCollector to store errors in context
type ErrorCollector struct {
	Errors []error
}

func (ec *ErrorCollector) Add(err error) {
	ec.Errors = append(ec.Errors, err)
}

// Error handling middleware
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		ctx := context.WithValue(r.Context(), "errors", &ErrorCollector{})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Result type for operations that can fail
type Result[T any] struct {
	Value T
	Err   error
}

func (r Result[T]) Unwrap() (T, error) {
	return r.Value, r.Err
}

// Example usage
func getUserByID(id string) Result[User] {
	if id == "" {
		return Result[User]{
			Err: &DomainError{
				Code:    ErrInvalidInput,
				Message: "user ID cannot be empty",
			},
		}
	}

	// Simulate a database lookup
	if id != "123" {
		return Result[User]{
			Err: &DomainError{
				Code:    ErrNotFound,
				Message: "user not found",
			},
		}
	}

	return Result[User]{
		Value: User{
			ID:   id,
			Name: "John Doe",
		},
	}
}

// User represents a user in the system
type User struct {
	ID   string
	Name string
}

// Simple HTTP handler using our error handling
func userHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	userResult := getUserByID(id)
	user, err := userResult.Unwrap()

	if err != nil {
		var domainErr *DomainError
		if errors.As(err, &domainErr) {
			switch domainErr.Code {
			case ErrNotFound:
				http.Error(w, domainErr.Error(), http.StatusNotFound)
			case ErrInvalidInput:
				http.Error(w, domainErr.Error(), http.StatusBadRequest)
			default:
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintf(w, "User: %s - %s", user.ID, user.Name)
}

func main() {
	// Set up a basic HTTP server with our error handling middleware
	mux := http.NewServeMux()
	mux.HandleFunc("/user", userHandler)

	handler := ErrorHandler(mux)

	fmt.Println("Server started at http://localhost:8080")
	fmt.Println("Try: http://localhost:8080/user?id=123")
	fmt.Println("Or: http://localhost:8080/user?id=456 (not found)")
	fmt.Println("Or: http://localhost:8080/user (invalid input)")

	log.Fatal(http.ListenAndServe(":8080", handler))
}
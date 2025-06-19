// Example 8
package main

import (
	"fmt"
	"strings"
)

// ValidationError is a custom error type that includes field information
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (v *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s: %s", v.Field, v.Message)
}

// Example usage function that validates an email
func validateEmail(email string) error {
	if email == "" {
		return &ValidationError{
			Field:   "email",
			Message: "cannot be empty",
		}
	}

	// Simple check for @ symbol
	if !contains(email, "@") {
		return &ValidationError{
			Field:   "email",
			Message: "must contain @ symbol",
		}
	}

	return nil
}

// Simple helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func main() {
	// Test with valid email
	err := validateEmail("user@example.com")
	if err != nil {
		fmt.Println("Unexpected error:", err)
	} else {
		fmt.Println("Valid email!")
	}

	// Test with invalid email (empty)
	err = validateEmail("")
	if err != nil {
		fmt.Println(err)
	}

	// Test with invalid email (missing @ symbol)
	err = validateEmail("invalid-email")
	if err != nil {
		fmt.Println(err)
	}
}
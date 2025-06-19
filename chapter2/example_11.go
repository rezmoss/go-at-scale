// Example 11
package main

import (
	"errors"
	"fmt"
)

type Handler[T any] func(T) error

func Chain[T any](handlers ...Handler[T]) Handler[T] {
	return func(t T) error {
		for _, h := range handlers {
			if err := h(t); err != nil {
				return err
			}
		}
		return nil
	}
}

// Example: Request validation pipeline
type Request struct {
	UserID string
	Data   []byte
}

func validateUserID(req Request) error {
	if req.UserID == "" {
		return errors.New("empty user ID")
	}
	return nil
}

func validateData(req Request) error {
	if len(req.Data) == 0 {
		return errors.New("empty data")
	}
	return nil
}

func main() {
	// Create the validator chain
	validator := Chain(validateUserID, validateData)

	// Test with valid request
	validReq := Request{
		UserID: "user123",
		Data:   []byte("sample data"),
	}

	err := validator(validReq)
	if err != nil {
		fmt.Printf("Valid request validation failed: %v\n", err)
	} else {
		fmt.Println("Valid request passed validation successfully")
	}

	// Test with invalid user ID
	invalidUserReq := Request{
		UserID: "",
		Data:   []byte("sample data"),
	}

	err = validator(invalidUserReq)
	if err != nil {
		fmt.Printf("Invalid user request validation error: %v\n", err)
	} else {
		fmt.Println("Invalid user request passed validation (unexpected)")
	}

	// Test with invalid data
	invalidDataReq := Request{
		UserID: "user123",
		Data:   []byte{},
	}

	err = validator(invalidDataReq)
	if err != nil {
		fmt.Printf("Invalid data request validation error: %v\n", err)
	} else {
		fmt.Println("Invalid data request passed validation (unexpected)")
	}
}
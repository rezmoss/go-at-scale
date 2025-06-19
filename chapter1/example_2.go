// Example 2
package main

import (
	"fmt"
	"net/http"
	"strings"
)

// Basic function type
type StringMapper func(string) string

// More complex function type
type HTTPHandler func(w http.ResponseWriter, r *http.Request) error

// Function type with multiple returns
type Validator func(interface{}) (bool, error)

// Implementation examples
func upperMapper(s string) string {
	return strings.ToUpper(s)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "OK")
	return nil
}

func numberValidator(v interface{}) (bool, error) {
	_, ok := v.(int)
	if !ok {
		return false, fmt.Errorf("value is not an integer")
	}
	return true, nil
}

func main() {
	// StringMapper example
	var mapper StringMapper = upperMapper
	result := mapper("hello")
	fmt.Printf("Mapped string: %s\n", result)

	// Validator example
	var validator Validator = numberValidator
	valid, err := validator(42)
	fmt.Printf("Validation result: %v, error: %v\n", valid, err)

	// HTTPHandler example (not starting server, just showing usage)
	var handler HTTPHandler = healthCheckHandler
	fmt.Printf("Handler type: %T\n", handler)
}
// Example 66
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

// Logger interface for logging
type Logger interface {
	Error(msg string, keysAndValues ...interface{})
}

// SimpleLogger implements the Logger interface
type SimpleLogger struct{}

func (l *SimpleLogger) Error(msg string, keysAndValues ...interface{}) {
	log.Printf("ERROR: %s %v", msg, keysAndValues)
}

// Define custom error types
type ValidationError struct {
	Message string
	Field   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}

type AuthenticationError struct {
	Message string
}

func (e *AuthenticationError) Error() string {
	return e.Message
}

// Resolver interface
type Resolver interface {
	Handle(ctx context.Context, next func(ctx context.Context) error) error
}

// gqlerror represents a GraphQL error
type gqlerror struct {
	Message    string                 `json:"message"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

// ErrorResolver struct as shown in the example
type ErrorResolver struct {
	underlying Resolver
	logger     Logger
}

func (r *ErrorResolver) Handle(ctx context.Context, next func(ctx context.Context) error) error {
	err := next(ctx)
	if err == nil {
		return nil
	}

	// Handle different error types
	switch e := err.(type) {
	case *ValidationError:
		return fmt.Errorf("%s: %s", e.Message, fmt.Sprintf(`{"code": "VALIDATION_ERROR", "field": "%s"}`, e.Field))
	case *AuthenticationError:
		return fmt.Errorf("Authentication required: %s", `{"code": "UNAUTHENTICATED"}`)
	default:
		r.logger.Error("unexpected error", "error", err)
		return fmt.Errorf("Internal server error")
	}
}

// Simple schema definition
const schemaString = `
  type Query {
    hello(name: String!): String!
    secure: String!
  }
`

// Root resolver
type RootResolver struct {
	errorResolver *ErrorResolver
	logger        Logger
}

// Hello resolves the hello query
func (r *RootResolver) Hello(ctx context.Context, args struct{ Name string }) (string, error) {
	if args.Name == "error" {
		return "", &ValidationError{
			Message: "Invalid name",
			Field:   "name",
		}
	}
	return "Hello " + args.Name, nil
}

// Secure resolves the secure query
func (r *RootResolver) Secure(ctx context.Context) (string, error) {
	// Simulate authentication check
	return "", &AuthenticationError{
		Message: "User not authenticated",
	}
}

func main() {
	logger := &SimpleLogger{}

	// Create error resolver
	errorResolver := &ErrorResolver{
		logger: logger,
	}

	// Create root resolver
	rootResolver := &RootResolver{
		errorResolver: errorResolver,
		logger:        logger,
	}

	// Parse schema
	schema := graphql.MustParseSchema(schemaString, rootResolver)

	// Setup HTTP handler
	http.Handle("/query", &relay.Handler{Schema: schema})
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
				<head>
					<title>GraphQL Playground</title>
					<style>
						body { font-family: Arial, sans-serif; margin: 20px; }
						h1 { color: #333; }
						pre { background-color: #f5f5f5; padding: 10px; border-radius: 5px; }
						.query { margin-bottom: 20px; }
					</style>
				</head>
				<body>
					<h1>GraphQL Error Handling Example</h1>
					<p>Try these queries at <a href="/query">/query</a> (POST requests):</p>
					
					<div class="query">
						<h3>Valid Query:</h3>
						<pre>{ "query": "{ hello(name: \"world\") }" }</pre>
					</div>
					
					<div class="query">
						<h3>Validation Error:</h3>
						<pre>{ "query": "{ hello(name: \"error\") }" }</pre>
					</div>
					
					<div class="query">
						<h3>Authentication Error:</h3>
						<pre>{ "query": "{ secure }" }</pre>
					</div>
				</body>
			</html>
		`))
	}))

	// Start server
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
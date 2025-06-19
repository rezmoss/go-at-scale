// Example 3
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(middlewares ...Middleware) Middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(middlewares) - 1; i >= 0; i-- {
				last = middlewares[i](last)
			}
			last(w, r)
		}
	}
}

// Example middlewares
func loggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Request started: %s %s", r.Method, r.URL.Path)
		next(w, r)
		log.Printf("Request completed in %v", time.Since(start))
	}
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if auth := r.Header.Get("Authorization"); auth == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	// Chain middlewares
	chain := Chain(loggerMiddleware, authMiddleware)
	wrappedHandler := chain(handler)

	// Create a test request
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "test-token")

	// Create a response recorder
	w := &responseRecorder{}

	// Execute the chain
	wrappedHandler(w, req)
}

// Simple response recorder for testing
type responseRecorder struct {
	headers http.Header
	body    []byte
	status  int
}

func (r *responseRecorder) Header() http.Header {
	if r.headers == nil {
		r.headers = make(http.Header)
	}
	return r.headers
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return len(b), nil
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
}
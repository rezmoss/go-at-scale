// Example 60
package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
)

// User represents a user in the system
type User struct {
	ID   string
	Name string
}

// SecureCookieHandler provides middleware for setting secure cookies
func SecureCookieHandler(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Secure cookie settings
			w.Header().Add("Set-Cookie", "session=value; HttpOnly; Secure; SameSite=Strict")
			next.ServeHTTP(w, r)
		})
	}
}

// RateLimiter implements a per-key rate limiting
type RateLimiter struct {
	limit  rate.Limit
	burst  int
	limits map[string]*rate.Limiter
	mu     sync.RWMutex
}

func (l *RateLimiter) Allow(key string) bool {
	l.mu.Lock()
	limiter, exists := l.limits[key]
	if !exists {
		limiter = rate.NewLimiter(l.limit, l.burst)
		l.limits[key] = limiter
	}
	l.mu.Unlock()
	return limiter.Allow()
}

// SafeQuery prevents SQL injection by using parameterized queries
func SafeQuery(db *sql.DB, id string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, name FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Name)
	if err != nil {
		return nil, fmt.Errorf("querying user: %w", err)
	}
	return &user, nil
}

// HashPassword creates a bcrypt hash of a plain text password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hashing password: %w", err)
	}
	return string(bytes), nil
}

// CheckPassword compares a password against a bcrypt hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// NewRateLimiter creates a new rate limiter with the specified rate and burst
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limit:  r,
		burst:  b,
		limits: make(map[string]*rate.Limiter),
	}
}

func main() {
	// Example of secure cookies
	fmt.Println("Secure Cookies Example:")
	// AES-256 requires a 32-byte key
	hashKey := []byte("the-32-byte-key-has-to-be-exact!")
	// Block key should be 16, 24, or 32 bytes for AES-128, AES-192, or AES-256
	blockKey := []byte("16-bytes-block-k")
	s := securecookie.New(hashKey, blockKey)

	// Encode a cookie value
	encoded, err := s.Encode("session", "example-value")
	if err != nil {
		fmt.Println("Error encoding cookie:", err)
		return
	}
	fmt.Println("Encoded cookie value:", encoded)

	// Example of using rate limiter
	fmt.Println("\nRate Limiter Example:")
	limiter := NewRateLimiter(1, 3) // 1 request per second with burst of 3

	// Simulate requests from same IP
	userIP := "192.168.1.1"
	for i := 0; i < 5; i++ {
		allowed := limiter.Allow(userIP)
		fmt.Printf("Request %d allowed: %v\n", i+1, allowed)
		time.Sleep(300 * time.Millisecond)
	}

	// Example of password hashing
	fmt.Println("\nPassword Hashing Example:")
	password := "my-secure-password"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return
	}

	fmt.Println("Original password:", password)
	fmt.Println("Hashed password:", hashedPassword)
	fmt.Println("Password check (correct):", CheckPassword(password, hashedPassword))
	fmt.Println("Password check (incorrect):", CheckPassword("wrong-password", hashedPassword))

	// Note: SafeQuery example is not run as it requires a database connection
	fmt.Println("\nSafeQuery function demonstrates proper parameterized queries to prevent SQL injection.")
}
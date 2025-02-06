// Example 60
// Secure cookie handler
func SecureCookieHandler(secret []byte) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Secure cookie settings
            w.Header().Add("Set-Cookie", "session=value; HttpOnly; Secure; SameSite=Strict")
            next.ServeHTTP(w, r)
        })
    }
}

// Rate limiter
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

// SQL injection prevention
func SafeQuery(db *sql.DB, id string) (*User, error) {
    var user User
    err := db.QueryRow("SELECT id, name FROM users WHERE id = $1", id).
        Scan(&user.ID, &user.Name)
    if err != nil {
        return nil, fmt.Errorf("querying user: %w", err)
    }
    return &user, nil
}

// Password hashing
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", fmt.Errorf("hashing password: %w", err)
    }
    return string(bytes), nil
}

func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
// Example 44
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
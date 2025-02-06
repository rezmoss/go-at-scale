// Example 2
// Basic function type
type StringMapper func(string) string

// More complex function type
type HTTPHandler func(w http.ResponseWriter, r *http.Request) error

// Function type with multiple returns
type Validator func(interface{}) (bool, error)
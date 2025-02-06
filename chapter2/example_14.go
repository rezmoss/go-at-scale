// Example 14
func Compose[A, B, C any](f func(B) C, g func(A) B) func(A) C {
    return func(a A) C {
        return f(g(a))
    }
}

// Example: String processing with type safety
func parseInt(s string) (int, error) {
    return strconv.Atoi(s)
}

func double(n int) int {
    return n * 2
}

// Usage
processor := Compose(double, parseInt)
result, err := processor("21")  // 42
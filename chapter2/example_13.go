// Example 13
// Comparable type constraint
type Ordered interface {
    ~int | ~float64 | ~string
}

// Generic map function
func Map[T, U any](slice []T, f func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = f(v)
    }
    return result
}

// Generic filter function
func Filter[T any](slice []T, predicate func(T) bool) []T {
    var result []T
    for _, v := range slice {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}

// Example usage
numbers := []int{1, 2, 3, 4, 5}
doubled := Map(numbers, func(x int) int { return x * 2 })
evens := Filter(numbers, func(x int) bool { return x%2 == 0 })
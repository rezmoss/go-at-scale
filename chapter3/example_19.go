// Example 19
// Generic Map implementation
func Map[T, U any](slice []T, f func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = f(v)
    }
    return result
}

// Generic Reduce implementation
func Reduce[T, U any](slice []T, initial U, f func(U, T) U) U {
    result := initial
    for _, v := range slice {
        result = f(result, v)
    }
    return result
}

// Generic Filter implementation
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
func main() {
    numbers := []int{1, 2, 3, 4, 5}
    
    // Double all numbers
    doubled := Map(numbers, func(x int) int {
        return x * 2
    })
    
    // Sum all numbers
    sum := Reduce(numbers, 0, func(acc, x int) int {
        return acc + x
    })
    
    // Get even numbers
    evens := Filter(numbers, func(x int) bool {
        return x%2 == 0
    })
}
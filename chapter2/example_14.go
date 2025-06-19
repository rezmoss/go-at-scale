// Example 14
package main

import (
	"fmt"
	"strconv"
)

func Compose[A, B, C any](f func(B) C, g func(A) B) func(A) C {
	return func(a A) C {
		return f(g(a))
	}
}

// Example: String processing with type safety
func parseInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func double(n int) int {
	return n * 2
}

func main() {
	// Usage
	processor := Compose(double, parseInt)
	result := processor("21") // 42
	fmt.Println(result)
}
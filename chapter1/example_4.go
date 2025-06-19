// Example 4
package main

import "fmt"

func createCounter() func() int {
	count := 0 // This variable is "closed over"
	return func() int {
		count++
		return count
	}
}

func main() {
	// Usage
	counter := createCounter()
	fmt.Println(counter()) // 1
	fmt.Println(counter()) // 2

	// Show that each closure maintains its own state
	counter2 := createCounter()
	fmt.Println(counter2()) // 1
	fmt.Println(counter())  // 3
}
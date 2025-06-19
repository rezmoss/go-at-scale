// Example 5
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	fmt.Println("Example with loop variable trap:")
	i := 0
	for i < 3 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(i) // Will print the final value of i (3)
		}()
		i++
	}
	wg.Wait()

	fmt.Println("\nExample with correct closure handling:")
	for i := 0; i < 3; i++ {
		wg.Add(1)
		current := i // Create new variable for closure
		go func() {
			defer wg.Done()
			fmt.Println(current) // Will print 0, 1, 2
		}()
	}
	wg.Wait()
}
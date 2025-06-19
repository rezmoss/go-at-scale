// Example 16
package main

import (
	"errors"
	"fmt"
)

// Impure function - depends on external state
var taxRate = 0.2

func calculateTax(amount float64) float64 {
	return amount * taxRate // Depends on external taxRate
}

// Pure function - all dependencies are explicit
func calculateTaxPure(amount, rate float64) float64 {
	return amount * rate // Self-contained
}

// Pure function with multiple returns
func divideWithRemainder(a, b int) (quotient, remainder int, err error) {
	if b == 0 {
		return 0, 0, errors.New("division by zero")
	}
	return a / b, a % b, nil
}

func main() {
	// Demonstrate the impure function
	fmt.Println("Impure function:")
	amount := 100.0
	fmt.Printf("Tax on $%.2f with external rate %.2f: $%.2f\n", amount, taxRate, calculateTax(amount))

	// Change the external state
	taxRate = 0.3
	fmt.Printf("Tax on $%.2f after changing external rate to %.2f: $%.2f\n", amount, taxRate, calculateTax(amount))

	// Demonstrate the pure function
	fmt.Println("\nPure function:")
	fmt.Printf("Tax on $%.2f with rate 0.2: $%.2f\n", amount, calculateTaxPure(amount, 0.2))
	fmt.Printf("Tax on $%.2f with rate 0.3: $%.2f\n", amount, calculateTaxPure(amount, 0.3))

	// Demonstrate pure function with multiple returns
	fmt.Println("\nPure function with multiple returns:")
	q1, r1, err1 := divideWithRemainder(10, 3)
	if err1 != nil {
		fmt.Println("Error:", err1)
	} else {
		fmt.Printf("10 ÷ 3 = %d with remainder %d\n", q1, r1)
	}

	q2, r2, err2 := divideWithRemainder(10, 0)
	if err2 != nil {
		fmt.Printf("10 ÷ 0: %s\n", err2)
	} else {
		fmt.Printf("10 ÷ 0 = %d with remainder %d\n", q2, r2)
	}
}
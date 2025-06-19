// Example 0
package main

import "fmt"

// Function type declaration
type PaymentProcessor func(amount float64) error

// Function implementing the type
func processPayment(amount float64) error {
	// Implementation
	fmt.Printf("Processing payment of $%.2f\n", amount)
	return nil
}

func main() {
	// Usage as a value
	var processor PaymentProcessor = processPayment

	// Example usage
	err := processor(99.99)
	if err != nil {
		fmt.Println("Payment failed:", err)
		return
	}
}
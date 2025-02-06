// Example 0
// Function type declaration
type PaymentProcessor func(amount float64) error

// Function implementing the type
func processPayment(amount float64) error {
    // Implementation
    return nil
}

// Usage as a value
var processor PaymentProcessor = processPayment
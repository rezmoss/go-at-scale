// Example 16
// Impure function - depends on external state
var taxRate = 0.2
func calculateTax(amount float64) float64 {
    return amount * taxRate  // Depends on external taxRate
}

// Pure function - all dependencies are explicit
func calculateTaxPure(amount, rate float64) float64 {
    return amount * rate  // Self-contained
}

// Pure function with multiple returns
func divideWithRemainder(a, b int) (quotient, remainder int, err error) {
    if b == 0 {
        return 0, 0, errors.New("division by zero")
    }
    return a / b, a % b, nil
}
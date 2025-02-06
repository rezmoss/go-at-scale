// Example 1
type PaymentProcessor func(ctx context.Context, payment Payment) error

type Payment struct {
    Amount      float64
    Currency    string
    CustomerID  string
    PaymentType PaymentMethod
}

func processPayment(ctx context.Context, payment Payment) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // Process payment
        return nil
    }
}

// Different payment processors
func stripePayment(amount float64) error {
    fmt.Printf("Processing %v via Stripe\n", amount)
    return nil
}

func paypalPayment(amount float64) error {
    fmt.Printf("Processing %v via PayPal\n", amount)
    return nil
}

// Usage
func main() {
    // Choose processor at runtime
    var processor PaymentProcessor
    if useStripe {
        processor = stripePayment
    } else {
        processor = paypalPayment
    }
    
    err := processPayment(99.99, processor)
    if err != nil {
        log.Fatal(err)
    }
}
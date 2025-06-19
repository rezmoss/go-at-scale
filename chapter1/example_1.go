// Example 1
package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

type PaymentProcessor func(ctx context.Context, payment Payment) error

type PaymentMethod string

const (
	CreditCard   PaymentMethod = "credit_card"
	BankTransfer PaymentMethod = "bank_transfer"
)

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
func stripePayment(ctx context.Context, payment Payment) error {
	fmt.Printf("Processing %.2f %s via Stripe\n", payment.Amount, payment.Currency)
	return nil
}

func paypalPayment(ctx context.Context, payment Payment) error {
	fmt.Printf("Processing %.2f %s via PayPal\n", payment.Amount, payment.Currency)
	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	payment := Payment{
		Amount:      99.99,
		Currency:    "USD",
		CustomerID:  "cust_123",
		PaymentType: CreditCard,
	}

	// Choose processor at runtime
	useStripe := true
	var processor PaymentProcessor
	if useStripe {
		processor = stripePayment
	} else {
		processor = paypalPayment
	}

	err := processor(ctx, payment)
	if err != nil {
		log.Fatal(err)
	}
}
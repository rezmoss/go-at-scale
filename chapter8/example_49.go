// Example 49
package main

import "fmt"

// Third-party payment service
type PayPalAPI struct{}

func (p *PayPalAPI) MakePayment(amount float64, currency string) error {
	// PayPal specific implementation
	fmt.Printf("PayPal payment processed: %.2f %s\n", amount, currency)
	return nil
}

// Our payment interface
type PaymentProvider interface {
	Process(payment Payment) error
}

type Payment struct {
	Amount   float64
	Currency string
}

// Adapter for PayPal
type PayPalAdapter struct {
	api *PayPalAPI
}

func (a *PayPalAdapter) Process(payment Payment) error {
	return a.api.MakePayment(payment.Amount, payment.Currency)
}

// Main function to demonstrate usage
func main() {
	// Create the PayPal API instance
	paypalAPI := &PayPalAPI{}

	// Create the adapter
	paypalAdapter := &PayPalAdapter{
		api: paypalAPI,
	}

	// Create a payment
	payment := Payment{
		Amount:   99.99,
		Currency: "USD",
	}

	// Use the adapter to process payment
	var provider PaymentProvider = paypalAdapter
	err := provider.Process(payment)

	if err != nil {
		fmt.Println("Payment failed:", err)
	} else {
		fmt.Println("Payment successful through the adapter!")
	}
}
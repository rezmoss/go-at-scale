// Example 49
// Third-party payment service
type PayPalAPI struct{}

func (p *PayPalAPI) MakePayment(amount float64, currency string) error {
    // PayPal specific implementation
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
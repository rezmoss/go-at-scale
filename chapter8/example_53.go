// Example 53
type PaymentStrategy interface {
    Pay(amount float64) error
}

type CreditCardStrategy struct {
    cardNumber string
    cvv        string
}

func (s *CreditCardStrategy) Pay(amount float64) error {
    // Process credit card payment
    return nil
}

type PayPalStrategy struct {
    email    string
    password string
}

func (s *PayPalStrategy) Pay(amount float64) error {
    // Process PayPal payment
    return nil
}

type ShoppingCart struct {
    paymentStrategy PaymentStrategy
    items          []Item
}

func (c *ShoppingCart) SetPaymentStrategy(strategy PaymentStrategy) {
    c.paymentStrategy = strategy
}

func (c *ShoppingCart) Checkout() error {
    amount := c.calculateTotal()
    return c.paymentStrategy.Pay(amount)
}
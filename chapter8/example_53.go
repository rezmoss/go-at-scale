// Example 53
package main

import (
	"fmt"
)

// PaymentStrategy interface defines the Pay method
type PaymentStrategy interface {
	Pay(amount float64) error
}

// CreditCardStrategy implements PaymentStrategy for credit card payments
type CreditCardStrategy struct {
	cardNumber string
	cvv        string
}

func (s *CreditCardStrategy) Pay(amount float64) error {
	fmt.Printf("Paid %.2f using Credit Card\n", amount)
	// Process credit card payment
	return nil
}

// PayPalStrategy implements PaymentStrategy for PayPal payments
type PayPalStrategy struct {
	email    string
	password string
}

func (s *PayPalStrategy) Pay(amount float64) error {
	fmt.Printf("Paid %.2f using PayPal\n", amount)
	// Process PayPal payment
	return nil
}

// Item represents a product in the shopping cart
type Item struct {
	name  string
	price float64
}

// ShoppingCart holds items and the payment strategy
type ShoppingCart struct {
	paymentStrategy PaymentStrategy
	items           []Item
}

// SetPaymentStrategy sets the payment method
func (c *ShoppingCart) SetPaymentStrategy(strategy PaymentStrategy) {
	c.paymentStrategy = strategy
}

// calculateTotal calculates the sum of all item prices
func (c *ShoppingCart) calculateTotal() float64 {
	var sum float64
	for _, item := range c.items {
		sum += item.price
	}
	return sum
}

// Checkout processes the payment for all items in the cart
func (c *ShoppingCart) Checkout() error {
	amount := c.calculateTotal()
	return c.paymentStrategy.Pay(amount)
}

// AddItem adds an item to the shopping cart
func (c *ShoppingCart) AddItem(item Item) {
	c.items = append(c.items, item)
}

func main() {
	// Create some items
	book := Item{name: "Design Patterns Book", price: 49.99}
	headphones := Item{name: "Bluetooth Headphones", price: 119.95}

	// Create a shopping cart
	cart := &ShoppingCart{}

	// Add items to cart
	cart.AddItem(book)
	cart.AddItem(headphones)

	// Create payment strategies
	creditCard := &CreditCardStrategy{
		cardNumber: "1234-5678-9012-3456",
		cvv:        "123",
	}

	paypal := &PayPalStrategy{
		email:    "example@example.com",
		password: "password123",
	}

	// Use credit card payment
	fmt.Println("Checking out with Credit Card:")
	cart.SetPaymentStrategy(creditCard)
	err := cart.Checkout()
	if err != nil {
		fmt.Println("Payment failed:", err)
	}

	// Use PayPal payment
	fmt.Println("\nChecking out with PayPal:")
	cart.SetPaymentStrategy(paypal)
	err = cart.Checkout()
	if err != nil {
		fmt.Println("Payment failed:", err)
	}
}
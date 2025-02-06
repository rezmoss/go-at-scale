// Example 160
// internal/testing/contract/consumer.go
type ConsumerTest struct {
    provider *pact.Pact
    client   *APIClient
}

func NewConsumerTest() *ConsumerTest {
    return &ConsumerTest{
        provider: &pact.Pact{
            Consumer: "OrderService",
            Provider: "PaymentService",
        },
    }
}

func (c *ConsumerTest) TestCreatePayment(t *testing.T) {
    // Set up expected interactions
    c.provider.
        AddInteraction().
        Given("A payment request").
        UponReceiving("A POST request to create payment").
        WithRequest(pact.Request{
            Method: "POST",
            Path:   "/payments",
            Headers: map[string]string{
                "Content-Type": "application/json",
            },
            Body: map[string]interface{}{
                "order_id": "123",
                "amount":   100.50,
                "currency": "USD",
            },
        }).
        WillRespondWith(pact.Response{
            Status: 201,
            Headers: map[string]string{
                "Content-Type": "application/json",
            },
            Body: map[string]interface{}{
                "id":     pact.Like("pay_123"),
                "status": "processing",
            },
        })

    // Run test with mock server
    verify := func() error {
        payment, err := c.client.CreatePayment(context.Background(), PaymentRequest{
            OrderID:  "123",
            Amount:   100.50,
            Currency: "USD",
        })
        if err != nil {
            return err
        }
        if payment.Status != "processing" {
            return fmt.Errorf("unexpected status: %s", payment.Status)
        }
        return nil
    }

    err := c.provider.Verify(verify)
    if err != nil {
        t.Fatalf("Error on Verify: %v", err)
    }
}
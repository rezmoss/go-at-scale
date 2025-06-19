// Example 67
package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Version     int                    `json:"version"`
	Timestamp   time.Time              `json:"timestamp"`
	AggregateID string                 `json:"aggregate_id"`
	Data        json.RawMessage        `json:"data"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// Event versioning support
type EventVersioner interface {
	Version() int
	Upgrade(nextVersion int) (Event, error)
}

// OrderCreated is an example of a concrete event
type OrderCreated struct {
	Event
}

func (e OrderCreated) Version() int {
	return e.Event.Version
}

func (e OrderCreated) Upgrade(nextVersion int) (Event, error) {
	if nextVersion <= e.Version() {
		return e.Event, fmt.Errorf("cannot upgrade to same or lower version")
	}

	// In a real scenario, here you would transform the event data
	// according to the new schema/version
	upgradedEvent := e.Event
	upgradedEvent.Version = nextVersion

	return upgradedEvent, nil
}

func main() {
	// Create a sample event
	orderData := json.RawMessage(`{"order_id": "12345", "customer_id": "C789", "amount": 99.99}`)

	event := Event{
		ID:          "evt-001",
		Type:        "order.created",
		Version:     1,
		Timestamp:   time.Now(),
		AggregateID: "order-12345",
		Data:        orderData,
		Metadata: map[string]interface{}{
			"source":     "web",
			"user_agent": "Mozilla/5.0",
			"ip_address": "192.168.1.1",
		},
	}

	// Create a versioned event
	orderCreated := OrderCreated{Event: event}

	// Demonstrate event serialization
	eventJSON, _ := json.MarshalIndent(event, "", "  ")
	fmt.Println("Original Event (Version 1):")
	fmt.Println(string(eventJSON))

	// Demonstrate versioning/upgrade
	upgradedEvent, err := orderCreated.Upgrade(2)
	if err != nil {
		fmt.Println("Error upgrading event:", err)
		return
	}

	upgradedJSON, _ := json.MarshalIndent(upgradedEvent, "", "  ")
	fmt.Println("\nUpgraded Event (Version 2):")
	fmt.Println(string(upgradedJSON))
}
// Example 64
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

// User represents a user in our system
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Logger is a simple interface for logging
type Logger interface {
	Error(msg string, args ...interface{})
}

// SimpleLogger implements the Logger interface
type SimpleLogger struct{}

func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	log.Printf("ERROR: "+msg, args...)
}

// Subscription represents a subscription to a topic
type Subscription struct {
	channel chan string
	closed  bool
}

// Channel returns the channel for this subscription
func (s *Subscription) Channel() <-chan string {
	return s.channel
}

// Close closes the subscription
func (s *Subscription) Close() {
	if !s.closed {
		s.closed = true
		close(s.channel)
	}
}

// Message represents a message in the pubsub system
type Message struct {
	Topic   string
	Payload string
}

// PubSub defines the interface for publishing and subscribing
type PubSub interface {
	Publish(ctx context.Context, topic string, payload interface{}) error
	Subscribe(ctx context.Context, topic string) *Subscription
}

// RedisPubSub implements PubSub using Redis
type RedisPubSub struct {
	client *redis.Client
}

// NewRedisPubSub creates a new Redis-backed PubSub
func NewRedisPubSub(addr string) *RedisPubSub {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisPubSub{client: client}
}

// Publish publishes a message to a topic
func (p *RedisPubSub) Publish(ctx context.Context, topic string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return p.client.Publish(ctx, topic, data).Err()
}

// Subscribe creates a subscription to a topic
func (p *RedisPubSub) Subscribe(ctx context.Context, topic string) *Subscription {
	pubsub := p.client.Subscribe(ctx, topic)
	ch := pubsub.Channel()

	subscription := &Subscription{
		channel: make(chan string),
		closed:  false,
	}

	go func() {
		defer func() {
			pubsub.Close()
			if !subscription.closed {
				subscription.closed = true
				close(subscription.channel)
			}
		}()

		for {
			select {
			case msg := <-ch:
				if subscription.closed {
					return
				}
				subscription.channel <- msg.Payload
			case <-ctx.Done():
				return
			}
		}
	}()

	return subscription
}

// SubscriptionResolver handles GraphQL subscriptions
type SubscriptionResolver struct {
	pubsub   PubSub
	upgrader websocket.Upgrader
	logger   Logger
}

// UserUpdates creates a subscription for updates about a specific user
func (r *SubscriptionResolver) UserUpdates(ctx context.Context, userID string) (<-chan *User, error) {
	updates := make(chan *User, 1)

	// Handle subscription
	go func() {
		defer close(updates)

		sub := r.pubsub.Subscribe(ctx, fmt.Sprintf("user:%s", userID))
		defer sub.Close()

		for {
			select {
			case msg := <-sub.Channel():
				var user User
				if err := json.Unmarshal([]byte(msg), &user); err != nil {
					r.logger.Error("failed to unmarshal user update", "error", err)
					continue
				}
				updates <- &user
			case <-ctx.Done():
				return
			}
		}
	}()

	return updates, nil
}

func main() {
	// Create a new PubSub instance
	pubsub := NewRedisPubSub("localhost:6379")

	// Create a subscription resolver
	resolver := &SubscriptionResolver{
		pubsub: pubsub,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all connections in this example
			},
		},
		logger: &SimpleLogger{},
	}

	// Set up a test user
	testUser := &User{ID: "123", Name: "John Doe", Email: "john@example.com"}

	// Create a cancelable context for the subscription
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a subscription for updates about this user
	updates, err := resolver.UserUpdates(ctx, testUser.ID)
	if err != nil {
		log.Fatalf("Failed to create subscription: %v", err)
	}

	// In a real application, this would be triggered by some user update event
	// For this example, we'll just periodically publish updates
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for i := 0; i < 4; i++ {
			<-ticker.C
			testUser.Name = fmt.Sprintf("John Doe %d", i)
			err := pubsub.Publish(ctx, fmt.Sprintf("user:%s", testUser.ID), testUser)
			if err != nil {
				log.Printf("Failed to publish update: %v", err)
			}
		}
		// After 4 updates, cancel the subscription
		<-ticker.C
		cancel()
	}()

	// Consume updates
	for update := range updates {
		log.Printf("Received user update: %+v", update)
	}

	log.Println("Subscription ended")
	time.Sleep(1 * time.Second) // Give goroutines time to clean up before exit
}
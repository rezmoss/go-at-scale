// Example 73
package eventstore_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

// Event represents a domain event
type Event struct {
	ID          string
	Type        string
	AggregateID string
	Version     int
	Data        []byte
	Timestamp   time.Time
}

// EventStore interface defines methods for event persistence
type EventStore interface {
	Save(ctx context.Context, events []Event) error
	Load(ctx context.Context, aggregateID string) ([]Event, error)
}

// Implementation of EventStore for testing
type InMemoryEventStore struct {
	events map[string][]Event
}

func NewInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		events: make(map[string][]Event),
	}
}

func (s *InMemoryEventStore) Save(ctx context.Context, events []Event) error {
	for _, event := range events {
		s.events[event.AggregateID] = append(s.events[event.AggregateID], event)
	}
	return nil
}

func (s *InMemoryEventStore) Load(ctx context.Context, aggregateID string) ([]Event, error) {
	return s.events[aggregateID], nil
}

// Helper function to generate test events
func generateTestEvents(aggregateID string, count int) []Event {
	events := make([]Event, count)
	for i := 0; i < count; i++ {
		events[i] = Event{
			ID:          uuid.New().String(),
			Type:        "TestEvent",
			AggregateID: aggregateID,
			Version:     i + 1,
			Data:        []byte(`{"test":"data"}`),
			Timestamp:   time.Now(),
		}
	}
	return events
}

// EventStoreTestSuite is the original test suite
type EventStoreTestSuite struct {
	suite.Suite
	store  EventStore
	events []Event
	ctx    context.Context
}

func (s *EventStoreTestSuite) SetupTest() {
	s.store = NewInMemoryEventStore()
	s.ctx = context.Background()
}

func (s *EventStoreTestSuite) TestSaveAndLoadEvents() {
	// Arrange
	aggregateID := uuid.New().String()
	events := generateTestEvents(aggregateID, 5)

	// Act
	err := s.store.Save(s.ctx, events)
	s.Require().NoError(err)

	loaded, err := s.store.Load(s.ctx, aggregateID)
	s.Require().NoError(err)

	// Assert
	s.Equal(len(events), len(loaded))
	for i, event := range events {
		s.Equal(event.ID, loaded[i].ID)
		s.Equal(event.Type, loaded[i].Type)
		s.Equal(event.Version, loaded[i].Version)
	}
}

// Run the test suite
func TestEventStore(t *testing.T) {
	suite.Run(t, new(EventStoreTestSuite))
}
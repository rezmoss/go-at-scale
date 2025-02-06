// Example 73
type EventStoreTestSuite struct {
    suite.Suite
    store      EventStore
    events     []Event
    ctx        context.Context
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
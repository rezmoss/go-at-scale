// Example 67
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
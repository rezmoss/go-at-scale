// Example 122
// internal/infrastructure/database/event.go
type Event struct {
    ID        uuid.UUID
    Type      string
    AggregateID uuid.UUID
    Data      json.RawMessage
    Metadata  json.RawMessage
    Version   int
    CreatedAt time.Time
}

type EventStore interface {
    SaveEvents(ctx context.Context, aggregateID uuid.UUID, events []Event) error
    GetEvents(ctx context.Context, aggregateID uuid.UUID) ([]Event, error)
}

type PostgresEventStore struct {
    db *sql.DB
}

func (es *PostgresEventStore) SaveEvents(ctx context.Context, aggregateID uuid.UUID, events []Event) error {
    return es.db.WithTransaction(ctx, func(tx *sql.Tx) error {
        for _, event := range events {
            const query = `
                INSERT INTO events (id, type, aggregate_id, data, metadata, version, created_at)
                VALUES ($1, $2, $3, $4, $5, $6, $7)
            `
            
            _, err := tx.ExecContext(ctx, query,
                event.ID,
                event.Type,
                event.AggregateID,
                event.Data,
                event.Metadata,
                event.Version,
                event.CreatedAt,
            )
            
            if err != nil {
                return fmt.Errorf("saving event: %w", err)
            }
        }
        
        return nil
    })
}
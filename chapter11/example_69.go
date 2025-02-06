// Example 69
type Aggregate interface {
    ID() string
    Version() int
    ApplyEvent(event Event) error
}

type EventStore interface {
    Save(ctx context.Context, events []Event) error
    Load(ctx context.Context, aggregateID string) ([]Event, error)
}

type PostgresEventStore struct {
    db      *sql.DB
    metrics MetricsRecorder
    logger  Logger
}

func (s *PostgresEventStore) Save(ctx context.Context, events []Event) error {
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("beginning transaction: %w", err)
    }
    defer tx.Rollback()

    for _, event := range events {
        if err := s.saveEvent(ctx, tx, event); err != nil {
            return err
        }
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("committing transaction: %w", err)
    }

    return nil
}
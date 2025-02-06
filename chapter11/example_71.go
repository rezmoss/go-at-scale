// Example 71
type Projector interface {
    Project(ctx context.Context, event Event) error
    Rebuild(ctx context.Context) error
}

type UserProjector struct {
    db          *sql.DB
    eventStore  EventStore
    metrics     MetricsRecorder
    logger      Logger
}

func (p *UserProjector) Project(ctx context.Context, event Event) error {
    switch event.Type {
    case "UserCreated":
        return p.handleUserCreated(ctx, event)
    case "UserUpdated":
        return p.handleUserUpdated(ctx, event)
    case "UserDeleted":
        return p.handleUserDeleted(ctx, event)
    default:
        return fmt.Errorf("unknown event type: %s", event.Type)
    }
}

func (p *UserProjector) Rebuild(ctx context.Context) error {
    // Start transaction
    tx, err := p.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("beginning transaction: %w", err)
    }
    defer tx.Rollback()

    // Clear existing projections
    if err := p.clearProjections(ctx, tx); err != nil {
        return err
    }

    // Load all events
    events, err := p.eventStore.LoadAll(ctx)
    if err != nil {
        return fmt.Errorf("loading events: %w", err)
    }

    // Project all events
    for _, event := range events {
        if err := p.projectWithTx(ctx, tx, event); err != nil {
            return err
        }
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("committing transaction: %w", err)
    }

    return nil
}
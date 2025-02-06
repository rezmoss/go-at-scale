// Example 163
// internal/testing/data/manager.go
type TestDataManager struct {
    db      *sql.DB
    fixture map[string][]string
    logger  Logger
}

func (m *TestDataManager) LoadFixture(ctx context.Context, name string) error {
    tx, err := m.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("beginning transaction: %w", err)
    }
    defer tx.Rollback()

    // Execute fixture SQL
    statements, ok := m.fixture[name]
    if !ok {
        return fmt.Errorf("fixture %s not found", name)
    }

    for _, stmt := range statements {
        if _, err := tx.ExecContext(ctx, stmt); err != nil {
            return fmt.Errorf("executing fixture statement: %w", err)
        }
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("committing fixture: %w", err)
    }

    return nil
}

func (m *TestDataManager) Cleanup(ctx context.Context) error {
    // Truncate all tables
    tables, err := m.getTables(ctx)
    if err != nil {
        return fmt.Errorf("getting tables: %w", err)
    }

    tx, err := m.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("beginning transaction: %w", err)
    }
    defer tx.Rollback()

    for _, table := range tables {
        if _, err := tx.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)); err != nil {
            return fmt.Errorf("truncating table %s: %w", table, err)
        }
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("committing cleanup: %w", err)
    }

    return nil
}
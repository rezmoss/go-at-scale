// Example 124
// internal/infrastructure/database/zero_downtime.go
type ZeroDowntimeMigration struct {
    Old      string
    New      string
    Copy     string
    Validate string
    Cleanup  string
}

func (mm *MigrationManager) ApplyZeroDowntimeMigration(ctx context.Context, migration ZeroDowntimeMigration) error {
    steps := []struct {
        name string
        fn   func(context.Context) error
    }{
        {"create new structure", mm.executeStep(migration.New)},
        {"copy data", mm.executeStep(migration.Copy)},
        {"validate data", mm.executeStep(migration.Validate)},
        {"cleanup", mm.executeStep(migration.Cleanup)},
    }
    
    for _, step := range steps {
        if err := step.fn(ctx); err != nil {
            return fmt.Errorf("step %s failed: %w", step.name, err)
        }
    }
    
    return nil
}

// Example zero-downtime migration
var userMigration = ZeroDowntimeMigration{
    New: `
        CREATE TABLE users_new (
            id UUID PRIMARY KEY,
            email TEXT UNIQUE NOT NULL,
            name TEXT NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )
    `,
    Copy: `
        INSERT INTO users_new (id, email, name, created_at)
        SELECT id, email, name, created_at FROM users
    `,
    Validate: `
        SELECT COUNT(*) FROM users_new
        WHERE id NOT IN (SELECT id FROM users)
    `,
    Cleanup: `
        ALTER TABLE users RENAME TO users_old;
        ALTER TABLE users_new RENAME TO users;
        DROP TABLE users_old;
    `,
}
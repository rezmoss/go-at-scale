// Example 120
// internal/infrastructure/database/repository.go
type UnitOfWork struct {
    tx        *sql.Tx
    completed bool
}

func (uow *UnitOfWork) Complete() error {
    if uow.completed {
        return nil
    }
    
    if err := uow.tx.Commit(); err != nil {
        return fmt.Errorf("committing transaction: %w", err)
    }
    
    uow.completed = true
    return nil
}

func (uow *UnitOfWork) Rollback() error {
    if uow.completed {
        return nil
    }
    
    if err := uow.tx.Rollback(); err != nil {
        return fmt.Errorf("rolling back transaction: %w", err)
    }
    
    uow.completed = true
    return nil
}

// Generic repository
type Repository[T any] interface {
    Create(ctx context.Context, entity T) error
    Update(ctx context.Context, entity T) error
    Delete(ctx context.Context, id interface{}) error
    FindByID(ctx context.Context, id interface{}) (T, error)
    FindAll(ctx context.Context) ([]T, error)
}

// Base repository implementation
type BaseRepository[T any] struct {
    db        *sql.DB
    tableName string
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity T) error {
    query := fmt.Sprintf("INSERT INTO %s (...) VALUES (...)", r.tableName)
    _, err := r.db.ExecContext(ctx, query /* args */)
    return err
}
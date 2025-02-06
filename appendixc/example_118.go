// Example 118
// internal/infrastructure/database/transaction.go
type TxManager struct {
    db *sql.DB
}

type TransactionFunc func(*sql.Tx) error

func (tm *TxManager) WithTransaction(ctx context.Context, fn TransactionFunc) error {
    tx, err := tm.db.BeginTx(ctx, &sql.TxOptions{
        Isolation: sql.LevelSerializable,
    })
    if err != nil {
        return fmt.Errorf("beginning transaction: %w", err)
    }
    
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p) // Re-throw panic after rollback
        }
    }()
    
    if err := fn(tx); err != nil {
        if rbErr := tx.Rollback(); rbErr != nil {
            return fmt.Errorf("rolling back transaction: %v (original error: %w)", rbErr, err)
        }
        return err
    }
    
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("committing transaction: %w", err)
    }
    
    return nil
}
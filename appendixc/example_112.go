// Example 112
// internal/infrastructure/repository/postgres.go
type PostgresOrderRepository struct {
    db *sql.DB
}

func (r *PostgresOrderRepository) Create(ctx context.Context, order *domain.Order) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("beginning transaction: %w", err)
    }
    defer tx.Rollback()
    
    // Insert order
    const orderQuery = `
        INSERT INTO orders (id, user_id, status, total, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
    
    if _, err := tx.ExecContext(ctx, orderQuery,
        order.ID, order.UserID, order.Status, order.Total,
        order.CreatedAt, order.UpdatedAt); err != nil {
        return fmt.Errorf("inserting order: %w", err)
    }
    
    // Insert order products
    const productQuery = `
        INSERT INTO order_products (order_id, product_id, quantity, price)
        VALUES ($1, $2, $3, $4)
    `
    
    for _, product := range order.Products {
        if _, err := tx.ExecContext(ctx, productQuery,
            order.ID, product.ProductID, product.Quantity, product.Price); err != nil {
            return fmt.Errorf("inserting order product: %w", err)
        }
    }
    
    return tx.Commit()
}
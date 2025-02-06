// Example 121
// internal/infrastructure/database/cqrs.go
type CommandDB struct {
    master *sql.DB
}

type QueryDB struct {
    slaves []*sql.DB
    current uint64
}

type OrderCommand struct {
    ID        uuid.UUID
    UserID    uuid.UUID
    Products  []OrderProduct
    Status    OrderStatus
}

type OrderQuery struct {
    db *QueryDB
}

func (q *OrderQuery) GetOrderSummary(ctx context.Context, orderID uuid.UUID) (*OrderSummary, error) {
    const query = `
        SELECT o.id, o.status, COUNT(op.id) as total_items, SUM(op.price * op.quantity) as total_amount
        FROM orders o
        LEFT JOIN order_products op ON o.id = op.order_id
        WHERE o.id = $1
        GROUP BY o.id, o.status
    `
    
    var summary OrderSummary
    err := q.db.QueryRowContext(ctx, query, orderID).Scan(
        &summary.ID,
        &summary.Status,
        &summary.TotalItems,
        &summary.TotalAmount,
    )
    
    if err != nil {
        return nil, fmt.Errorf("querying order summary: %w", err)
    }
    
    return &summary, nil
}
// Example 110
// internal/domain/models.go
package domain

import (
    "time"
    "github.com/google/uuid"
)

type OrderStatus string

const (
    OrderStatusPending   OrderStatus = "pending"
    OrderStatusApproved OrderStatus = "approved"
    OrderStatusRejected OrderStatus = "rejected"
)

type Order struct {
    ID        uuid.UUID
    UserID    uuid.UUID
    Products  []OrderProduct
    Status    OrderStatus
    Total     float64
    CreatedAt time.Time
    UpdatedAt time.Time
}

type OrderProduct struct {
    ProductID uuid.UUID
    Quantity  int
    Price     float64
}

// Domain errors
type OrderError struct {
    Code    string
    Message string
    Err     error
}

func (e *OrderError) Error() string {
    if e.Err != nil {
        return e.Message + ": " + e.Err.Error()
    }
    return e.Message
}
// Example 111
// internal/application/interfaces.go
package application

import (
    "context"
    "github.com/google/uuid"
    "yourservice/internal/domain"
)

type OrderRepository interface {
    Create(ctx context.Context, order *domain.Order) error
    Get(ctx context.Context, id uuid.UUID) (*domain.Order, error)
    Update(ctx context.Context, order *domain.Order) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type OrderService interface {
    CreateOrder(ctx context.Context, order *domain.Order) error
    GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, error)
    UpdateOrderStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error
}

// internal/application/services.go
type orderService struct {
    repo      OrderRepository
    validator Validator
    publisher EventPublisher
    logger    Logger
    metrics   MetricsRecorder
}

func NewOrderService(
    repo OrderRepository,
    validator Validator,
    publisher EventPublisher,
    logger Logger,
    metrics MetricsRecorder,
) OrderService {
    return &orderService{
        repo:      repo,
        validator: validator,
        publisher: publisher,
        logger:    logger,
        metrics:   metrics,
    }
}

func (s *orderService) CreateOrder(ctx context.Context, order *domain.Order) error {
    // Start transaction span
    span, ctx := tracer.StartSpanFromContext(ctx, "CreateOrder")
    defer span.End()
    
    // Validate order
    if err := s.validator.Validate(order); err != nil {
        s.logger.Error("order validation failed", "error", err)
        s.metrics.IncCounter("order_validation_failures")
        return &domain.OrderError{
            Code:    "VALIDATION_ERROR",
            Message: "Invalid order data",
            Err:     err,
        }
    }
    
    // Create order
    if err := s.repo.Create(ctx, order); err != nil {
        s.logger.Error("failed to create order", "error", err)
        s.metrics.IncCounter("order_creation_failures")
        return err
    }
    
    // Publish event
    event := NewOrderCreatedEvent(order)
    if err := s.publisher.Publish(ctx, "orders.created", event); err != nil {
        s.logger.Error("failed to publish order created event", "error", err)
        // Don't return error as order is already created
    }
    
    s.metrics.IncCounter("orders_created")
    s.logger.Info("order created successfully", "order_id", order.ID)
    return nil
}
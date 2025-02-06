// Example 113
// internal/ports/http/handlers.go
type OrderHandler struct {
    service application.OrderService
    logger  Logger
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    // Extract correlation ID
    correlationID := r.Header.Get("X-Correlation-ID")
    if correlationID == "" {
        correlationID = uuid.New().String()
    }
    
    ctx := context.WithValue(r.Context(), "correlation_id", correlationID)
    
    // Parse request
    var req CreateOrderRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.logger.Error("failed to decode request", "error", err)
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    
    // Convert to domain model
    order := &domain.Order{
        ID:        uuid.New(),
        UserID:    req.UserID,
        Products:  make([]domain.OrderProduct, len(req.Products)),
        Status:    domain.OrderStatusPending,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    // Create order
    if err := h.service.CreateOrder(ctx, order); err != nil {
        h.logger.Error("failed to create order", 
            "error", err,
            "correlation_id", correlationID)
            
        switch e := err.(type) {
        case *domain.OrderError:
            respondWithError(w, http.StatusBadRequest, e.Message)
        default:
            respondWithError(w, http.StatusInternalServerError, 
                "Failed to process order")
        }
        return
    }
    
    // Respond with success
    respondWithJSON(w, http.StatusCreated, order)
}

// internal/ports/grpc/server.go
type OrderServer struct {
    pb.UnimplementedOrderServiceServer
    service application.OrderService
    logger  Logger
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
    // Extract metadata
    md, ok := metadata.FromIncomingContext(ctx)
    correlationID := "unknown"
    if ok {
        if values := md.Get("x-correlation-id"); len(values) > 0 {
            correlationID = values[0]
        }
    }
    
    ctx = context.WithValue(ctx, "correlation_id", correlationID)
    
    // Convert to domain model
    order := &domain.Order{
        ID:     uuid.New(),
        UserID: uuid.MustParse(req.UserId),
        // ... convert other fields
    }
    
    if err := s.service.CreateOrder(ctx, order); err != nil {
        s.logger.Error("failed to create order",
            "error", err,
            "correlation_id", correlationID)
            
        return nil, status.Errorf(codes.Internal,
            "Failed to create order: %v", err)
    }
    
    // Convert to protobuf response
    return convertOrderToProto(order), nil
}
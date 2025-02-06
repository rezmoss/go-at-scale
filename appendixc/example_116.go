// Example 116
// cmd/server/main.go
func main() {
    // Load configuration
    config := loadConfig()
    
    // Initialize logger
    logger := initLogger(config)
    
    // Initialize metrics
    metrics := initMetrics(config)
    
    // Initialize tracer
    tracer := initTracer(config)
    
    // Initialize database
    db := initDatabase(config)
    
    // Initialize repositories
    orderRepo := repository.NewPostgresOrderRepository(db)
    
    // Initialize event publisher
    publisher := messaging.NewRabbitMQPublisher(config)
    
    // Initialize service
    orderService := application.NewOrderService(
        orderRepo,
        validator.New(),
        publisher,
        logger,
        metrics,
    )
    
    // Initialize HTTP handlers
    orderHandler := http.NewOrderHandler(orderService, logger)
    
    // Initialize gRPC server
    grpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(
            grpc_middleware.ChainUnaryServer(
                grpc_recovery.UnaryServerInterceptor(),
                grpc_validator.UnaryServerInterceptor(),
                grpc_prometheus.UnaryServerInterceptor,
            ),
        ),
    )
    
    pb.RegisterOrderServiceServer(grpcServer, 
        grpchandlers.NewOrderServer(orderService, logger))
    
    // Start servers
    go startGRPCServer(grpcServer, config)
    startHTTPServer(orderHandler, config)
}
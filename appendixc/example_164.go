// Example 164
// internal/serverless/handler.go
type LambdaHandler struct {
    service    Service
    validator  Validator
    logger     Logger
    metrics    MetricsRecorder
}

func (h *LambdaHandler) HandleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    start := time.Now()
    defer func() {
        h.metrics.ObserveLatency("lambda_execution", time.Since(start))
    }()

    // Extract correlation ID
    correlationID := event.Headers["X-Correlation-ID"]
    if correlationID == "" {
        correlationID = uuid.New().String()
    }
    ctx = context.WithValue(ctx, "correlation_id", correlationID)

    // Parse request
    var req Request
    if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
        h.logger.Error("failed to parse request",
            "error", err,
            "correlation_id", correlationID)
        return createResponse(http.StatusBadRequest, ErrorResponse{
            Error: "Invalid request format",
        })
    }

    // Validate request
    if err := h.validator.Validate(req); err != nil {
        h.metrics.IncCounter("validation_errors")
        return createResponse(http.StatusBadRequest, ErrorResponse{
            Error: err.Error(),
        })
    }

    // Process request
    result, err := h.service.Process(ctx, req)
    if err != nil {
        h.logger.Error("processing failed",
            "error", err,
            "correlation_id", correlationID)
        return createResponse(http.StatusInternalServerError, ErrorResponse{
            Error: "Internal server error",
        })
    }

    h.metrics.IncCounter("successful_requests")
    return createResponse(http.StatusOK, result)
}

func createResponse(statusCode int, body interface{}) (events.APIGatewayProxyResponse, error) {
    jsonBody, err := json.Marshal(body)
    if err != nil {
        return events.APIGatewayProxyResponse{}, err
    }

    return events.APIGatewayProxyResponse{
        StatusCode: statusCode,
        Headers: map[string]string{
            "Content-Type": "application/json",
        },
        Body: string(jsonBody),
    }, nil
}
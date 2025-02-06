// Example 138
// internal/monitoring/logging/logger.go
type LoggerConfig struct {
    Level      string
    Format     string
    OutputPath string
}

type StructuredLogger struct {
    logger *zap.Logger
}

func NewStructuredLogger(config LoggerConfig) (*StructuredLogger, error) {
    cfg := zap.Config{
        Level:       parseLevel(config.Level),
        Development: false,
        Sampling: &zap.SamplingConfig{
            Initial:    100,
            Thereafter: 100,
        },
        Encoding:         config.Format,
        EncoderConfig:    getEncoderConfig(),
        OutputPaths:      []string{config.OutputPath},
        ErrorOutputPaths: []string{config.OutputPath},
    }
    
    logger, err := cfg.Build()
    if err != nil {
        return nil, fmt.Errorf("building logger: %w", err)
    }
    
    return &StructuredLogger{logger: logger}, nil
}

func (l *StructuredLogger) With(fields ...zap.Field) *StructuredLogger {
    return &StructuredLogger{
        logger: l.logger.With(fields...),
    }
}

// Request logging middleware
func (l *StructuredLogger) HTTPMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Extract trace ID from context or header
        traceID := r.Header.Get("X-Trace-ID")
        if traceID == "" {
            traceID = uuid.New().String()
        }
        
        // Add trace ID to response headers
        w.Header().Set("X-Trace-ID", traceID)
        
        // Create child logger with request context
        requestLogger := l.With(
            zap.String("trace_id", traceID),
            zap.String("method", r.Method),
            zap.String("path", r.URL.Path),
            zap.String("remote_addr", r.RemoteAddr),
            zap.String("user_agent", r.UserAgent()),
        )
        
        // Add logger to request context
        ctx := context.WithValue(r.Context(), "logger", requestLogger)
        
        // Use custom response writer to capture status code
        ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
        
        next.ServeHTTP(ww, r.WithContext(ctx))
        
        // Log request completion
        requestLogger.Info("request completed",
            zap.Int("status", ww.Status()),
            zap.Int("bytes", ww.BytesWritten()),
            zap.Duration("duration", time.Since(start)),
        )
    })
}
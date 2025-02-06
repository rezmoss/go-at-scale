// Example 59
type App struct {
    server   *http.Server
    db       *sql.DB
    cache    *redis.Client
    shutdown chan struct{}
}

func (a *App) Start() error {
    // Start HTTP server
    go func() {
        if err := a.server.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatalf("HTTP server error: %v", err)
        }
    }()
    
    // Handle shutdown signals
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
    
    <-stop
    return a.Shutdown()
}

func (a *App) Shutdown() error {
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Close HTTP server
    if err := a.server.Shutdown(ctx); err != nil {
        return fmt.Errorf("server shutdown: %w", err)
    }
    
    // Close database connections
    if err := a.db.Close(); err != nil {
        return fmt.Errorf("database shutdown: %w", err)
    }
    
    // Close cache connections
    if err := a.cache.Close(); err != nil {
        return fmt.Errorf("cache shutdown: %w", err)
    }
    
    close(a.shutdown)
    return nil
}
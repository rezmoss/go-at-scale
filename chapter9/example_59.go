// Example 59
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

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

func main() {
	// Setup a simple HTTP server with a basic handler
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second) // Simulate a long-running request
		fmt.Fprintf(w, "Hello, World!")
	})

	// Initialize PostgreSQL database connection
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Create the application
	app := &App{
		server: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
		db:       db,
		cache:    redisClient,
		shutdown: make(chan struct{}),
	}

	// Start the application
	log.Println("Server is starting on :8080...")
	if err := app.Start(); err != nil {
		log.Fatalf("Failed to shutdown gracefully: %v", err)
	}
	log.Println("Server stopped gracefully")
}
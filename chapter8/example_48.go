// Example 48
package main

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Database struct holds the database connection
type Database struct {
	conn *sql.DB
}

var (
	instance *Database
	once     sync.Once
)

// GetDatabase returns a singleton instance of Database
func GetDatabase() *Database {
	once.Do(func() {
		db, err := sql.Open("postgres", "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}
		instance = &Database{conn: db}
	})
	return instance
}

func main() {
	// Get the database instance
	db := GetDatabase()

	// Use the same instance again
	db2 := GetDatabase()

	// Verify that both variables reference the same instance
	log.Printf("Are db and db2 the same instance? %v\n", db == db2)

	// Verify connection works
	err := db.conn.Ping()
	if err != nil {
		log.Printf("Connection error: %v\n", err)
	} else {
		log.Println("Successfully connected to database")
	}
}
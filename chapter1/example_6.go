// Example 6
package main

import (
	"fmt"
	"time"
)

type Server struct {
	host    string
	port    int
	timeout time.Duration
	maxConn int
}

type ServerOption func(*Server)

func WithTimeout(t time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = t
	}
}

func WithMaxConnections(n int) ServerOption {
	return func(s *Server) {
		s.maxConn = n
	}
}

func NewServer(host string, options ...ServerOption) *Server {
	// Default values
	s := &Server{
		host:    host,
		port:    8080,
		timeout: 30 * time.Second,
		maxConn: 100,
	}

	// Apply options
	for _, option := range options {
		option(s)
	}

	return s
}

func main() {
	// Create a server with default options
	server1 := NewServer("localhost")
	fmt.Printf("Server 1: %+v\n", server1)

	// Create a server with custom timeout
	server2 := NewServer("localhost", WithTimeout(60*time.Second))
	fmt.Printf("Server 2: %+v\n", server2)

	// Create a server with custom max connections
	server3 := NewServer("localhost", WithMaxConnections(1000))
	fmt.Printf("Server 3: %+v\n", server3)

	// Create a server with multiple custom options
	server4 := NewServer("localhost",
		WithTimeout(120*time.Second),
		WithMaxConnections(500))
	fmt.Printf("Server 4: %+v\n", server4)
}
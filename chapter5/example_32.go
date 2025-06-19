// Example 32
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Process struct {
	done chan struct{}
	err  error
	once sync.Once
}

func NewProcess() *Process {
	return &Process{
		done: make(chan struct{}),
	}
}

func (p *Process) Stop(err error) {
	p.once.Do(func() {
		p.err = err
		close(p.done)
	})
}

func (p *Process) Wait() error {
	<-p.done
	return p.err
}

// Example: Graceful shutdown
type Server struct {
	processes []*Process
	done      chan struct{}
}

func NewServer() *Server {
	return &Server{
		processes: make([]*Process, 0),
		done:      make(chan struct{}),
	}
}

func (s *Server) AddProcess(p *Process) {
	s.processes = append(s.processes, p)
}

func (s *Server) Shutdown(ctx context.Context) error {
	// Signal shutdown to all processes
	close(s.done)

	// Wait for all processes with timeout
	for _, proc := range s.processes {
		select {
		case <-proc.done:
			// Process completed
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

// Simulate a process that does work
func runProcess(p *Process, id int) {
	defer p.Stop(nil)

	fmt.Printf("Process %d: Starting work\n", id)

	// Simulate some work
	for i := 0; i < 5; i++ {
		select {
		case <-p.done:
			fmt.Printf("Process %d: Received stop signal, cleaning up\n", id)
			return
		default:
			fmt.Printf("Process %d: Working... step %d\n", id, i+1)
			time.Sleep(300 * time.Millisecond)
		}
	}

	fmt.Printf("Process %d: Work completed successfully\n", id)
}

func main() {
	// Create a server
	server := NewServer()

	// Create some processes
	for i := 0; i < 3; i++ {
		proc := NewProcess()
		server.AddProcess(proc)

		// Start each process in its own goroutine
		go runProcess(proc, i)
	}

	// Let processes run for a while
	fmt.Println("Server running... press Enter to initiate graceful shutdown")
	fmt.Scanln() // Wait for user input

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Shutdown the server
	fmt.Println("Initiating graceful shutdown...")
	err := server.Shutdown(ctx)
	if err != nil {
		fmt.Printf("Error during shutdown: %v\n", err)
	} else {
		fmt.Println("Server shutdown completed successfully")
	}
}
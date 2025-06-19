// Example 46
package main

import (
	"fmt"
	"os"
	"time"
)

// Product interface
type Logger interface {
	Log(message string)
}

// Concrete products
type FileLogger struct {
	file *os.File
}

func (l *FileLogger) Log(message string) {
	fmt.Fprintf(l.file, "[%v] %s\n", time.Now(), message)
}

type ConsoleLogger struct{}

func (l *ConsoleLogger) Log(message string) {
	fmt.Printf("[%v] %s\n", time.Now(), message)
}

// Factory interface
type LoggerFactory interface {
	CreateLogger() (Logger, error)
}

// Concrete factories
type FileLoggerFactory struct {
	filepath string
}

func (f *FileLoggerFactory) CreateLogger() (Logger, error) {
	file, err := os.OpenFile(f.filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("creating file logger: %w", err)
	}
	return &FileLogger{file: file}, nil
}

type ConsoleLoggerFactory struct{}

func (f *ConsoleLoggerFactory) CreateLogger() (Logger, error) {
	return &ConsoleLogger{}, nil
}

func main() {
	// Create a console logger factory
	consoleFactory := &ConsoleLoggerFactory{}
	consoleLogger, err := consoleFactory.CreateLogger()
	if err != nil {
		fmt.Printf("Error creating console logger: %v\n", err)
		return
	}
	consoleLogger.Log("This is a log message to console")

	// Create a file logger factory
	fileFactory := &FileLoggerFactory{filepath: "log.txt"}
	fileLogger, err := fileFactory.CreateLogger()
	if err != nil {
		fmt.Printf("Error creating file logger: %v\n", err)
		return
	}
	fileLogger.Log("This is a log message to file")

	fmt.Println("Logging complete. Check log.txt for file logs.")
}
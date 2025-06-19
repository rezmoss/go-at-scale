// Example 12
package main

import (
	"context"
	"fmt"
)

// Logger interface for logging errors
type Logger interface {
	Error(msg string, args ...interface{})
}

// Simple logger implementation
type SimpleLogger struct{}

func (l SimpleLogger) Error(msg string, args ...interface{}) {
	fmt.Printf("%s: %v\n", msg, args)
}

type Pipeline[T any] struct {
	stages []PipelineStage[T, T]
	logger Logger
}

type PipelineStage[In, Out any] struct {
	Process func(context.Context, In) (Out, error)
	Cleanup func() error
}

func (p *Pipeline[T]) Execute(ctx context.Context, input T) (T, error) {
	var result T
	var err error

	result = input

	// Execute all stages
	for _, stage := range p.stages {
		result, err = stage.Process(ctx, result)
		if err != nil {
			return result, fmt.Errorf("pipeline execution error: %w", err)
		}
	}

	// Defer cleanup to ensure it runs after execution
	defer func() {
		for _, stage := range p.stages {
			if err := stage.Cleanup(); err != nil {
				p.logger.Error("cleanup error", "error", err)
			}
		}
	}()

	return result, nil
}

// Function to create a pipeline from stages
func NewPipeline[T any](logger Logger, stages ...PipelineStage[T, T]) *Pipeline[T] {
	return &Pipeline[T]{
		stages: stages,
		logger: logger,
	}
}

// Helper function to create a simple pipeline stage
func CreateStage[In, Out any](process func(context.Context, In) (Out, error)) PipelineStage[In, Out] {
	return PipelineStage[In, Out]{
		Process: process,
		Cleanup: func() error { return nil }, // Default no-op cleanup
	}
}

// Example: Text processing pipeline
func addPrefix(prefix string) PipelineStage[string, string] {
	return CreateStage(func(ctx context.Context, s string) (string, error) {
		return prefix + s, nil
	})
}

func addSuffix(suffix string) PipelineStage[string, string] {
	return CreateStage(func(ctx context.Context, s string) (string, error) {
		return s + suffix, nil
	})
}

func main() {
	// Create a logger
	logger := SimpleLogger{}

	// Create a pipeline
	processor := NewPipeline(
		logger,
		addPrefix("Hello, "),
		addSuffix("!"),
	)

	// Execute the pipeline
	ctx := context.Background()
	result, err := processor.Execute(ctx, "World")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Result: %s\n", result)
}
// Example 30
package main

import (
	"errors"
	"fmt"
	"log"
)

// Generic pipeline stage
type Stage[In, Out any] struct {
	Process func(in In) (Out, error)
}

// Pipeline definition
type Pipeline[T any] struct {
	stages []Stage[T, T]
	done   chan struct{}
}

func NewPipeline[T any](stages ...Stage[T, T]) *Pipeline[T] {
	return &Pipeline[T]{
		stages: stages,
		done:   make(chan struct{}),
	}
}

func (p *Pipeline[T]) Run(input <-chan T) <-chan T {
	output := make(chan T)

	go func() {
		defer close(output)
		defer close(p.done)

		for value := range input {
			result := value
			var err error

			// Process through each stage
			for _, stage := range p.stages {
				result, err = stage.Process(result)
				if err != nil {
					log.Printf("Pipeline error: %v", err)
					return
				}
			}

			output <- result
		}
	}()

	return output
}

func main() {
	// Create pipeline stages
	pipeline := NewPipeline(
		Stage[int, int]{
			Process: func(n int) (int, error) {
				return n * 2, nil
			},
		},
		Stage[int, int]{
			Process: func(n int) (int, error) {
				if n > 100 {
					return 0, errors.New("value too large")
				}
				return n + 1, nil
			},
		},
	)

	// Create input channel
	input := make(chan int)
	output := pipeline.Run(input)

	// Feed data and collect results
	go func() {
		defer close(input)
		for i := 0; i < 5; i++ {
			input <- i
		}
	}()

	for result := range output {
		fmt.Printf("Result: %d\n", result)
	}
}
// Example 12
type PipelineStage[In, Out any] struct {
    Process func(context.Context, In) (Out, error)
    Cleanup func() error
}

func (p *Pipeline[T]) Execute(ctx context.Context) error {
    defer func() {
        for _, stage := range p.stages {
            if err := stage.Cleanup(); err != nil {
                p.logger.Error("cleanup error", "error", err)
            }
        }
    }()
    // Pipeline execution
}

func Pipeline[T any](stages ...PipelineStage[T, T]) PipelineStage[T, T] {
    return func(input T) (T, error) {
        var err error
        result := input
        
        for _, stage := range stages {
            result, err = stage(result)
            if err != nil {
                return result, fmt.Errorf("pipeline error: %w", err)
            }
        }
        
        return result, nil
    }
}

// Example: Text processing pipeline
func addPrefix(prefix string) PipelineStage[string, string] {
    return func(s string) (string, error) {
        return prefix + s, nil
    }
}

func addSuffix(suffix string) PipelineStage[string, string] {
    return func(s string) (string, error) {
        return s + suffix, nil
    }
}

// Usage
processor := Pipeline(
    addPrefix("Hello, "),
    addSuffix("!"),
)
result, err := processor("World")  // "Hello, World!"
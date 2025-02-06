// Example 11
type Handler[T any] func(T) error

func Chain[T any](handlers ...Handler[T]) Handler[T] {
    return func(t T) error {
        for _, h := range handlers {
            if err := h(t); err != nil {
                return err
            }
        }
        return nil
    }
}

// Example: Request validation pipeline
type Request struct {
    UserID string
    Data   []byte
}

func validateUserID(req Request) error {
    if req.UserID == "" {
        return errors.New("empty user ID")
    }
    return nil
}

func validateData(req Request) error {
    if len(req.Data) == 0 {
        return errors.New("empty data")
    }
    return nil
}

// Usage
validator := Chain(validateUserID, validateData)
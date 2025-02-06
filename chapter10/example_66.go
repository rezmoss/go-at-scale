// Example 66
type ErrorResolver struct {
    underlying Resolver
    logger     Logger
}

func (r *ErrorResolver) Handle(ctx context.Context, next func(ctx context.Context) error) error {
    err := next(ctx)
    if err == nil {
        return nil
    }
    
    // Handle different error types
    switch e := err.(type) {
    case *ValidationError:
        return &gqlerror.Error{
            Message: e.Message,
            Extensions: map[string]interface{}{
                "code": "VALIDATION_ERROR",
                "field": e.Field,
            },
        }
    case *AuthenticationError:
        return &gqlerror.Error{
            Message: "Authentication required",
            Extensions: map[string]interface{}{
                "code": "UNAUTHENTICATED",
            },
        }
    default:
        r.logger.Error("unexpected error", "error", err)
        return &gqlerror.Error{
            Message: "Internal server error",
        }
    }
}
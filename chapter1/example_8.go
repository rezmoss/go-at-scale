// Example 8
type ValidationError struct {
    Field string
    Error string
}

func (v *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %s", v.Field, v.Error)
}
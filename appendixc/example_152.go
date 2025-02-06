// Example 152

// internal/infrastructure/messaging/schema.go
type SchemaRegistry struct {
    schemas map[string]Schema
    mu      sync.RWMutex
}


type Schema struct {
    Version     int
    Definition  string
    Validators  []SchemaValidator
    Migrations  []SchemaMigration
}


func (r *SchemaRegistry) ValidateMessage(eventType string, version int, data []byte) error {
    schema, err := r.GetSchema(eventType, version)
    if err != nil {
        return fmt.Errorf("getting schema: %w", err)
    }


    for _, validator := range schema.Validators {
        if err := validator.Validate(data); err != nil {
            return fmt.Errorf("validating against schema: %w", err)
        }
    }


    return nil

}
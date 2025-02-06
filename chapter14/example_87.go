// Example 87
type DataValidator struct {
    source      DataReader
    destination DataReader
    metrics     MetricsRecorder
    logger      Logger
}

func (v *DataValidator) ValidateMigration(ctx context.Context) error {
    // Get source and destination checksums
    sourceHash, err := v.calculateChecksum(ctx, v.source)
    if err != nil {
        return fmt.Errorf("calculating source checksum: %w", err)
    }

    destHash, err := v.calculateChecksum(ctx, v.destination)
    if err != nil {
        return fmt.Errorf("calculating destination checksum: %w", err)
    }

    // Compare checksums
    if sourceHash != destHash {
        v.metrics.IncCounter("validation_failures")
        return ErrChecksumMismatch
    }

    // Validate data integrity
    if err := v.validateDataIntegrity(ctx); err != nil {
        return fmt.Errorf("validating data integrity: %w", err)
    }

    // Validate business rules
    if err := v.validateBusinessRules(ctx); err != nil {
        return fmt.Errorf("validating business rules: %w", err)
    }

    return nil
}

func (v *DataValidator) validateDataIntegrity(ctx context.Context) error {
    stream := v.source.Stream(ctx)
    for record := range stream {
        if err := v.validateRecord(ctx, record); err != nil {
            v.metrics.IncCounter("record_validation_failures")
            return err
        }
    }
    return nil
}
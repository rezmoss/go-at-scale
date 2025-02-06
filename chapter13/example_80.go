// Example 80
// Interface for different backup strategies
type BackupStrategy interface {
    Backup(ctx context.Context, source Source) (*BackupMetadata, error)
    Restore(ctx context.Context, metadata *BackupMetadata) error
    Verify(ctx context.Context, metadata *BackupMetadata) error
}

// Backup metadata for tracking and verification
type BackupMetadata struct {
    ID            string
    CreatedAt     time.Time
    Size          int64
    Checksum      string
    BackupType    BackupType
    Dependencies  []string
    Configuration map[string]string
}

// Implementation of a database backup strategy
type DatabaseBackup struct {
    storage     StorageProvider
    encryption  EncryptionService
    compression CompressionService
    metrics     MetricsRecorder
    logger      Logger
}

type BackupResult struct {
    Metadata *BackupMetadata
    Errors   []BackupError
    Corrupt  bool
}

func (b *DatabaseBackup) Backup(ctx context.Context, source Source) (*BackupResult, error) {
    result := &BackupResult{
        Errors: make([]BackupError, 0),
    }
    
    checksum, err := b.createBackup(ctx, source, result)
    if err != nil {
        return result, err
    }
    
    if err := b.verifyBackup(ctx, checksum); err != nil {
        result.Corrupt = true
        result.Errors = append(result.Errors, BackupError{
            Type: BackupCorrupted,
            Err:  err,
        })
    }
    
    return result, nil
}
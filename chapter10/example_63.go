// Example 63
type DataLoader struct {
    userLoader *dataloader.Loader
    postLoader *dataloader.Loader
    redis      *redis.Client
    logger     Logger
}

func NewDataLoader(ctx context.Context, repo Repository) *DataLoader {
    return &DataLoader{
        userLoader: dataloader.NewBatchedLoader(func(keys []string) []*dataloader.Result {
            // Batch load users
            users, err := repo.GetUsersByIDs(ctx, keys)
            if err != nil {
                return makeBatchError(err, len(keys))
            }
            
            // Map results to keys order
            return mapResultsToKeys(keys, users)
        }),
        // ... other loaders
    }
}
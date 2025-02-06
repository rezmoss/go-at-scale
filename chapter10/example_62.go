// Example 62
type Resolvers struct {
    userService  UserService
    postService  PostService
    dataLoader   *DataLoader
    metrics      MetricsRecorder
    logger       Logger
}

func (r *Resolvers) User(ctx context.Context, id string) (*User, error) {
    start := time.Now()
    defer func() {
        r.metrics.ObserveLatency("resolver.user", time.Since(start))
    }()

    // Use dataloader for batching
    return r.dataLoader.LoadUser(ctx, id)
}

// Field resolver pattern
func (r *Resolvers) User_posts(ctx context.Context, obj *User) ([]*Post, error) {
    return r.postService.GetUserPosts(ctx, obj.ID)
}
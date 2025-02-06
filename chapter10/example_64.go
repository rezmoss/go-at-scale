// Example 64
type SubscriptionResolver struct {
    pubsub   PubSub
    upgrader websocket.Upgrader
    logger   Logger
}

func (r *SubscriptionResolver) UserUpdates(ctx context.Context, userID string) (<-chan *User, error) {
    updates := make(chan *User, 1)
    
    // Handle subscription
    go func() {
        defer close(updates)
        
        sub := r.pubsub.Subscribe(ctx, fmt.Sprintf("user:%s", userID))
        defer sub.Close()
        
        for {
            select {
            case msg := <-sub.Channel():
                var user User
                if err := json.Unmarshal([]byte(msg.Payload), &user); err != nil {
                    r.logger.Error("failed to unmarshal user update", "error", err)
                    continue
                }
                updates <- &user
            case <-ctx.Done():
                return
            }
        }
    }()
    
    return updates, nil
}
// Example 135
// internal/auth/session/manager.go
type SessionManager struct {
    store  SessionStore
    crypto SessionEncryptor
}

type Session struct {
    ID        string
    UserID    string
    Device    string
    ExpiresAt time.Time
    Metadata  map[string]interface{}
}

func (sm *SessionManager) CreateSession(ctx context.Context, userID, device string) (*Session, error) {
    session := &Session{
        ID:        uuid.New().String(),
        UserID:    userID,
        Device:    device,
        ExpiresAt: time.Now().Add(time.Hour * 24),
        Metadata:  make(map[string]interface{}),
    }
    
    // Encrypt sensitive data
    encrypted, err := sm.crypto.Encrypt(session)
    if err != nil {
        return nil, fmt.Errorf("encrypting session: %w", err)
    }
    

    // Store session
    if err := sm.store.SaveSession(ctx, encrypted); err != nil {
        return nil, fmt.Errorf("saving session: %w", err)
    }
    

    return session, nil
}


func (sm *SessionManager) GetSession(ctx context.Context, sessionID string) (*Session, error) {
    encrypted, err := sm.store.GetSession(ctx, sessionID)
    if err != nil {
        return nil, fmt.Errorf("getting session: %w", err)
    }
    

    session, err := sm.crypto.Decrypt(encrypted)
    if err != nil {
        return nil, fmt.Errorf("decrypting session: %w", err)
    }
    

    if time.Now().After(session.ExpiresAt) {
        sm.store.DeleteSession(ctx, sessionID)
        return nil, fmt.Errorf("session expired")
    }

    
    return session, nil

}
// Example 131
// internal/auth/refresh/manager.go
type RefreshTokenManager struct {
    store  TokenStore
    crypto TokenEncryptor
    ttl    time.Duration
}

type RefreshToken struct {
    Token     string    `json:"token"`
    UserID    string    `json:"user_id"`
    ExpiresAt time.Time `json:"expires_at"`
    Device    string    `json:"device"`
}

func (rm *RefreshTokenManager) GenerateRefreshToken(userID, device string) (*RefreshToken, error) {
    token := &RefreshToken{
        Token:     uuid.New().String(),
        UserID:    userID,
        ExpiresAt: time.Now().Add(rm.ttl),
        Device:    device,
    }
    
    // Encrypt token before storage
    encryptedToken, err := rm.crypto.Encrypt(token.Token)
    if err != nil {
        return nil, fmt.Errorf("encrypting token: %w", err)
    }
    
    token.Token = encryptedToken
    
    // Store refresh token
    if err := rm.store.SaveToken(context.Background(), token); err != nil {
        return nil, fmt.Errorf("saving refresh token: %w", err)
    }
    
    return token, nil
}

func (rm *RefreshTokenManager) ValidateRefreshToken(tokenString string) (*RefreshToken, error) {
    // Decrypt token
    decryptedToken, err := rm.crypto.Decrypt(tokenString)
    if err != nil {
        return nil, fmt.Errorf("decrypting token: %w", err)
    }
    
    // Get token from store
    token, err := rm.store.GetToken(context.Background(), decryptedToken)
    if err != nil {
        return nil, fmt.Errorf("getting refresh token: %w", err)
    }
    
    if time.Now().After(token.ExpiresAt) {
        return nil, fmt.Errorf("refresh token expired")
    }
    
    return token, nil
}
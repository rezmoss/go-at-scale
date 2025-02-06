// Example 130
// internal/auth/jwt/manager.go
type TokenManager struct {
    privateKey *rsa.PrivateKey
    publicKey  *rsa.PublicKey
    issuer    string
    duration  time.Duration
}

type TokenClaims struct {
    jwt.RegisteredClaims
    UserID    string   `json:"uid"`
    Roles     []string `json:"roles"`
    SessionID string   `json:"sid"`
}

func NewTokenManager(privateKeyPEM, publicKeyPEM []byte, issuer string, duration time.Duration) (*TokenManager, error) {
    privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyPEM)
    if err != nil {
        return nil, fmt.Errorf("parsing private key: %w", err)
    }
    
    publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyPEM)
    if err != nil {
        return nil, fmt.Errorf("parsing public key: %w", err)
    }
    
    return &TokenManager{
        privateKey: privateKey,
        publicKey:  publicKey,
        issuer:    issuer,
        duration:  duration,
    }, nil
}

func (tm *TokenManager) GenerateToken(userID string, roles []string) (string, error) {
    now := time.Now()
    claims := TokenClaims{
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    tm.issuer,
            IssuedAt:  jwt.NewNumericDate(now),
            ExpiresAt: jwt.NewNumericDate(now.Add(tm.duration)),
            NotBefore: jwt.NewNumericDate(now),
            ID:        uuid.New().String(),
        },
        UserID:    userID,
        Roles:     roles,
        SessionID: uuid.New().String(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
    return token.SignedString(tm.privateKey)
}

func (tm *TokenManager) ValidateToken(tokenString string) (*TokenClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return tm.publicKey, nil
    })
    
    if err != nil {
        return nil, fmt.Errorf("parsing token: %w", err)
    }
    
    if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, fmt.Errorf("invalid token claims")
}
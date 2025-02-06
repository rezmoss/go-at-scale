// Example 132
// internal/auth/oauth/manager.go
type OAuth2Manager struct {
    configs map[string]*oauth2.Config
    store   UserStore
}

type OAuth2Profile struct {
    ID            string
    Email         string
    Name          string
    Provider      string
    AccessToken   string
    RefreshToken  string
    ExpiresAt     time.Time
}

func (om *OAuth2Manager) AuthorizeURL(provider, state string) (string, error) {
    config, ok := om.configs[provider]
    if !ok {
        return "", fmt.Errorf("unknown provider: %s", provider)
    }
    
    return config.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

func (om *OAuth2Manager) Exchange(ctx context.Context, provider, code string) (*OAuth2Profile, error) {
    config, ok := om.configs[provider]
    if !ok {
        return nil, fmt.Errorf("unknown provider: %s", provider)
    }
    
    token, err := config.Exchange(ctx, code)
    if err != nil {
        return nil, fmt.Errorf("exchanging code: %w", err)
    }
    
    profile, err := om.fetchProfile(ctx, provider, token)
    if err != nil {
        return nil, fmt.Errorf("fetching profile: %w", err)
    }
    
    return profile, nil
}

func (om *OAuth2Manager) fetchProfile(ctx context.Context, provider string, token *oauth2.Token) (*OAuth2Profile, error) {
    switch provider {
    case "google":
        return om.fetchGoogleProfile(ctx, token)
    case "github":
        return om.fetchGithubProfile(ctx, token)
    default:
        return nil, fmt.Errorf("unsupported provider: %s", provider)
    }
}
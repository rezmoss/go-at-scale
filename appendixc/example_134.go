// Example 134
// internal/auth/rbac/manager.go
type RoleManager struct {
    store RoleStore
    cache *cache.Cache
}

type Role struct {
    Name        string
    Permissions []string
}

func (rm *RoleManager) AssignRole(ctx context.Context, userID, role string) error {
    if err := rm.store.AssignRole(ctx, userID, role); err != nil {
        return fmt.Errorf("assigning role: %w", err)
    }
    // Invalidate cache
    rm.cache.Delete(fmt.Sprintf("user_roles:%s", userID))
    return nil
}

func (rm *RoleManager) HasPermission(ctx context.Context, userID, permission string) bool {
    roles, err := rm.getUserRoles(ctx, userID)
    if err != nil {
        return false
    }
    
    for _, role := range roles {
        if containsPermission(role.Permissions, permission) {
            return true
        }
    }
    
    return false
}

func (rm *RoleManager) getUserRoles(ctx context.Context, userID string) ([]Role, error) {
    // Check cache
    if roles, ok := rm.cache.Get(fmt.Sprintf("user_roles:%s", userID)); ok {
        return roles.([]Role), nil
    }
    // Get from store
    roles, err := rm.store.GetUserRoles(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("getting user roles: %w", err)
    }
    // Cache roles
    rm.cache.Set(fmt.Sprintf("user_roles:%s", userID), roles, time.Minute*5)
    
    return roles, nil
}
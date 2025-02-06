// Example 119
// internal/infrastructure/database/query.go
type QueryBuilder struct {
    sb    strings.Builder
    args  []interface{}
    count int
}

func (qb *QueryBuilder) AddWhere(condition string, args ...interface{}) {
    if qb.count == 0 {
        qb.sb.WriteString(" WHERE ")
    } else {
        qb.sb.WriteString(" AND ")
    }
    
    qb.sb.WriteString(condition)
    qb.args = append(qb.args, args...)
    qb.count++
}

// Prepared statement manager
type PreparedStatements struct {
    statements map[string]*sql.Stmt
    mu         sync.RWMutex
}

func (ps *PreparedStatements) Get(ctx context.Context, db *sql.DB, query string) (*sql.Stmt, error) {
    ps.mu.RLock()
    stmt, exists := ps.statements[query]
    ps.mu.RUnlock()
    
    if exists {
        return stmt, nil
    }
    
    ps.mu.Lock()
    defer ps.mu.Unlock()
    
    // Double-check after acquiring write lock
    if stmt, exists := ps.statements[query]; exists {
        return stmt, nil
    }
    
    stmt, err := db.PrepareContext(ctx, query)
    if err != nil {
        return nil, err
    }
    
    ps.statements[query] = stmt
    return stmt, nil
}
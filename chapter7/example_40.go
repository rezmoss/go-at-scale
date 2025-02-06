// Example 40
// Bad: Too many dependencies
package user

import (
    "database/sql"
    "net/http"
    "encoding/json"
    "html/template"
    // ... many more imports
)

// Good: Focused functionality
package user

import (
    "context"
    "errors"
    "time"
)

// Bad: Mixed levels of abstraction
type UserService struct {
    db        *sql.DB
    cache     *redis.Client
    templates *template.Template
}

// Good: Clean abstraction
type UserService struct {
    store  Repository
    cache  Cache
    events EventEmitter
}
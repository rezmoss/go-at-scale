// Example 61
type Schema struct {
    resolvers *Resolvers
    loader    *DataLoader
}

// Common interfaces
type Node interface {
    ID() string
}

type Edge interface {
    Node() Node
    Cursor() string
}

// Connection pattern for pagination
type UserConnection struct {
    Edges    []*UserEdge
    PageInfo PageInfo
}

type UserEdge struct {
    Node   *User
    Cursor string
}

type PageInfo struct {
    HasNextPage     bool
    HasPreviousPage bool
    StartCursor     string
    EndCursor       string
}

// Input type patterns
type CreateUserInput struct {
    Email    string
    Username string
    Role     UserRole
}
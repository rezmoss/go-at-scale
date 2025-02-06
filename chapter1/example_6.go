// Example 6
type Server struct {
    host string
    port int
    timeout time.Duration
    maxConn int
}

type ServerOption func(*Server)

func WithTimeout(t time.Duration) ServerOption {
    return func(s *Server) {
        s.timeout = t
    }
}

func WithMaxConnections(n int) ServerOption {
    return func(s *Server) {
        s.maxConn = n
    }
}

func NewServer(host string, options ...ServerOption) *Server {
    // Default values
    s := &Server{
        host:    host,
        port:    8080,
        timeout: 30 * time.Second,
        maxConn: 100,
    }
    
    // Apply options
    for _, option := range options {
        option(s)
    }
    
    return s
}
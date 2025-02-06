// Example 47
// Server configuration builder
type ServerConfig struct {
    Host         string
    Port         int
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    TLS         *TLSConfig
}

type TLSConfig struct {
    CertFile string
    KeyFile  string
}

type ServerConfigBuilder struct {
    config *ServerConfig
}

func NewServerConfigBuilder() *ServerConfigBuilder {
    return &ServerConfigBuilder{
        config: &ServerConfig{
            Port:         8080,
            ReadTimeout:  5 * time.Second,
            WriteTimeout: 5 * time.Second,
        },
        validators: []ConfigValidator{
            validatePort,
            validateTimeouts,
            validateTLS,
        },
    }
}

func (b *ServerConfigBuilder) Host(host string) *ServerConfigBuilder {
    b.config.Host = host
    return b
}

func (b *ServerConfigBuilder) Port(port int) *ServerConfigBuilder {
    b.config.Port = port
    return b
}

func (b *ServerConfigBuilder) Timeouts(read, write time.Duration) *ServerConfigBuilder {
    b.config.ReadTimeout = read
    b.config.WriteTimeout = write
    return b
}

func (b *ServerConfigBuilder) WithTLS(certFile, keyFile string) *ServerConfigBuilder {
    b.config.TLS = &TLSConfig{
        CertFile: certFile,
        KeyFile:  keyFile,
    }
    return b
}

func (b *ServerConfigBuilder) Build() (*ServerConfig, error) {
    for _, validator := range b.validators {
        if err := validator(b.config); err != nil {
            return nil, fmt.Errorf("invalid configuration: %w", err)
        }
    }
    return b.config, nil
}

// Usage
config := NewServerConfigBuilder().
    Host("localhost").
    Port(9000).
    Timeouts(10*time.Second, 10*time.Second).
    WithTLS("cert.pem", "key.pem").
    Build()
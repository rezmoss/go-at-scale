// Example 47
package main

import (
	"fmt"
	"time"
)

// Server configuration builder
type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	TLS          *TLSConfig
}

type TLSConfig struct {
	CertFile string
	KeyFile  string
}

type ServerConfigBuilder struct {
	config     *ServerConfig
	validators []ConfigValidator
}

// ConfigValidator defines a function that validates a configuration
type ConfigValidator func(*ServerConfig) error

func validatePort(config *ServerConfig) error {
	if config.Port < 1 || config.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	return nil
}

func validateTimeouts(config *ServerConfig) error {
	if config.ReadTimeout <= 0 || config.WriteTimeout <= 0 {
		return fmt.Errorf("timeouts must be positive")
	}
	return nil
}

func validateTLS(config *ServerConfig) error {
	if config.TLS != nil {
		if config.TLS.CertFile == "" || config.TLS.KeyFile == "" {
			return fmt.Errorf("TLS configuration requires both cert and key files")
		}
	}
	return nil
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

func main() {
	config, err := NewServerConfigBuilder().
		Host("localhost").
		Port(9000).
		Timeouts(10*time.Second, 10*time.Second).
		WithTLS("cert.pem", "key.pem").
		Build()

	if err != nil {
		fmt.Printf("Error building config: %v\n", err)
		return
	}

	fmt.Println("Server Configuration:")
	fmt.Printf("  Host: %s\n", config.Host)
	fmt.Printf("  Port: %d\n", config.Port)
	fmt.Printf("  Read Timeout: %v\n", config.ReadTimeout)
	fmt.Printf("  Write Timeout: %v\n", config.WriteTimeout)
	if config.TLS != nil {
		fmt.Println("  TLS Enabled:")
		fmt.Printf("    Cert File: %s\n", config.TLS.CertFile)
		fmt.Printf("    Key File: %s\n", config.TLS.KeyFile)
	}
}
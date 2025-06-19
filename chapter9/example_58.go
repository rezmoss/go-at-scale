// Example 58
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

// Config struct holds all configuration values
type Config struct {
	Server struct {
		Host string `yaml:"host" envconfig:"SERVER_HOST"`
		Port int    `yaml:"port" envconfig:"SERVER_PORT"`
	} `yaml:"server"`

	Database struct {
		Host     string `yaml:"host" envconfig:"DB_HOST"`
		Port     int    `yaml:"port" envconfig:"DB_PORT"`
		User     string `yaml:"user" envconfig:"DB_USER"`
		Password string `yaml:"password" envconfig:"DB_PASSWORD"`
		Name     string `yaml:"name" envconfig:"DB_NAME"`
	} `yaml:"database"`

	Redis struct {
		Host     string        `yaml:"host" envconfig:"REDIS_HOST"`
		Port     int           `yaml:"port" envconfig:"REDIS_PORT"`
		Password string        `yaml:"password" envconfig:"REDIS_PASSWORD"`
		Timeout  time.Duration `yaml:"timeout" envconfig:"REDIS_TIMEOUT"`
	} `yaml:"redis"`
}

// LoadConfig loads configuration from a YAML file and overrides with environment variables
func LoadConfig(path string) (*Config, error) {
	var config Config

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	// Override with environment variables
	if err := envconfig.Process("APP", &config); err != nil {
		return nil, fmt.Errorf("processing env vars: %w", err)
	}

	return &config, nil
}

func main() {
	// Load the configuration
	config, err := LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Display the loaded configuration
	fmt.Println("Configuration loaded successfully:")
	fmt.Printf("Server: %s:%d\n", config.Server.Host, config.Server.Port)
	fmt.Printf("Database: %s:%d, User: %s, DB: %s\n",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Name)
	fmt.Printf("Redis: %s:%d, Timeout: %v\n",
		config.Redis.Host, config.Redis.Port, config.Redis.Timeout)
}
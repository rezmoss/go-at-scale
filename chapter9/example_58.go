// Example 58
type Config struct {
    Server struct {
        Host string `yaml:"host"`
        Port int    `yaml:"port"`
    } `yaml:"server"`
    
    Database struct {
        Host     string `yaml:"host"`
        Port     int    `yaml:"port"`
        User     string `yaml:"user"`
        Password string `yaml:"password"`
        Name     string `yaml:"name"`
    } `yaml:"database"`
    
    Redis struct {
        Host     string        `yaml:"host"`
        Port     int          `yaml:"port"`
        Password string       `yaml:"password"`
        Timeout  time.Duration `yaml:"timeout"`
    } `yaml:"redis"`
}

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
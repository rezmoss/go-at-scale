// Example 117
// internal/infrastructure/database/pool.go
type DBPool struct {
    master  *sql.DB
    slaves  []*sql.DB
    current uint64 // For round-robin slave selection
    mu      sync.RWMutex
}

func NewDBPool(config Config) (*DBPool, error) {
    master, err := sql.Open("postgres", config.MasterDSN)
    if err != nil {
        return nil, fmt.Errorf("connecting to master: %w", err)
    }
    
    // Configure master pool
    master.SetMaxOpenConns(config.MaxOpenConns)
    master.SetMaxIdleConns(config.MaxIdleConns)
    master.SetConnMaxLifetime(config.ConnMaxLifetime)
    
    // Initialize slave connections
    slaves := make([]*sql.DB, len(config.SlaveDSNs))
    for i, dsn := range config.SlaveDSNs {
        slave, err := sql.Open("postgres", dsn)
        if err != nil {
            return nil, fmt.Errorf("connecting to slave %d: %w", i, err)
        }
        
        // Configure slave pool
        slave.SetMaxOpenConns(config.MaxOpenConns)
        slave.SetMaxIdleConns(config.MaxIdleConns)
        slave.SetConnMaxLifetime(config.ConnMaxLifetime)
        
        slaves[i] = slave
    }
    
    return &DBPool{
        master: master,
        slaves: slaves,
    }, nil
}

func (p *DBPool) Master() *sql.DB {
    return p.master
}

func (p *DBPool) Slave() *sql.DB {
    if len(p.slaves) == 0 {
        return p.master
    }
    
    // Round-robin slave selection
    current := atomic.AddUint64(&p.current, 1)
    return p.slaves[current%uint64(len(p.slaves))]
}
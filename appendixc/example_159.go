// Example 159
// internal/testing/integration/framework.go
type TestFramework struct {
    containers map[string]*TestContainer
    network    *docker.Network
    cleanup    []func() error
    logger     Logger
}

type TestContainer struct {
    ID       string
    Name     string
    Host     string
    Ports    map[string]string
    Ready    chan bool
    cleanup  func() error
}

func NewTestFramework(ctx context.Context) (*TestFramework, error) {
    cli, err := client.NewClientWithOpts(client.FromEnv)
    if err != nil {
        return nil, fmt.Errorf("creating docker client: %w", err)
    }

    // Create network
    network, err := cli.NetworkCreate(ctx, "test-network", types.NetworkCreate{
        Driver: "bridge",
    })
    if err != nil {
        return nil, fmt.Errorf("creating network: %w", err)
    }

    return &TestFramework{
        containers: make(map[string]*TestContainer),
        network:    &network,
        cleanup:    make([]func() error, 0),
    }, nil
}

func (f *TestFramework) StartPostgres(ctx context.Context) (*TestContainer, error) {
    container, err := f.startContainer(ctx, ContainerConfig{
        Image: "postgres:13",
        Env: []string{
            "POSTGRES_PASSWORD=test",
            "POSTGRES_DB=testdb",
        },
        HealthCheck: func(host string, port string) error {
            db, err := sql.Open("postgres", 
                fmt.Sprintf("postgres://postgres:test@%s:%s/testdb?sslmode=disable", 
                    host, port))
            if err != nil {
                return err
            }
            return db.Ping()
        },
    })
    if err != nil {
        return nil, fmt.Errorf("starting postgres: %w", err)
    }

    return container, nil
}

func (f *TestFramework) Cleanup(ctx context.Context) error {
    for _, cleanup := range f.cleanup {
        if err := cleanup(); err != nil {
            f.logger.Error("cleanup failed", "error", err)
        }
    }
    return nil
}
// Example 48
type Database struct {
    conn *sql.DB
}

var (
    instance *Database
    once     sync.Once
)

func GetDatabase() *Database {
    once.Do(func() {
        db, err := sql.Open("postgres", "connection-string")
        if err != nil {
            log.Fatal(err)
        }
        instance = &Database{conn: db}
    })
    return instance
}
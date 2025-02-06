// Example 109
service/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── models.go
│   │   └── errors.go
│   ├── ports/
│   │   ├── http/
│   │   ├── grpc/
│   │   └── handlers.go
│   ├── application/
│   │   ├── services.go
│   │   └── interfaces.go
│   └── infrastructure/
│       ├── repository/
│       ├── cache/
│       └── messaging/
├── pkg/
│   ├── logger/
│   ├── metrics/
│   └── validator/
└── api/
    ├── grpc/
    └── http/
// Example 99
// Internal packages
project/
  ├── internal/     # Private to this module
  │   └── auth/
  ├── pkg/          # Public packages
  │   └── util/
  └── go.mod

// Versioned API
import "example.com/pkg/v2"

// Multiple major versions
project/
  ├── v1/
  │   └── go.mod
  └── v2/
      └── go.mod
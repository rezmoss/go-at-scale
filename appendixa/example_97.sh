// Example 97
# Initialize a new module
go mod init example.com/myproject

# Download dependencies
go mod download

# Update dependencies
go get -u ./...
go get -u=patch ./...  # Only patch updates

# Clean up dependencies
go mod tidy

# Verify dependencies
go mod verify

# List all dependencies
go list -m all
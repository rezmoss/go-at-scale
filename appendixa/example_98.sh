// Example 98
# Get specific version
go get example.com/pkg@v1.2.3

# Get latest version
go get -u example.com/pkg

# List available versions
go list -m -versions example.com/pkg

# Use local replacement
# go.mod
replace example.com/pkg => ../local/pkg
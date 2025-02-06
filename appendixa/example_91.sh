// Example 91
# CPU profiling
go test -cpuprofile cpu.prof -bench .

# Memory profiling
go test -memprofile mem.prof -bench .

# Analyze profiles
go tool pprof cpu.prof
go tool pprof -http=:8080 cpu.prof  # Web interface

# Common pprof commands
(pprof) top10           # Show top 10 functions
(pprof) list function   # Show source code
(pprof) web            # View in browser
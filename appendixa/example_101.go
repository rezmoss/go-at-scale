// Example 101
var memStats runtime.MemStats

func logMemStats() {
    runtime.ReadMemStats(&memStats)
    log.Printf("Alloc = %v MiB", memStats.Alloc/1024/1024)
    log.Printf("TotalAlloc = %v MiB", memStats.TotalAlloc/1024/1024)
    log.Printf("Sys = %v MiB", memStats.Sys/1024/1024)
    log.Printf("NumGC = %v", memStats.NumGC)
}
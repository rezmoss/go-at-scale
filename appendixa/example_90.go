// Example 90
// main_test.go
func BenchmarkComplexOperation(b *testing.B) {
    // Enable CPU profiling
    f, err := os.Create("cpu.prof")
    if err != nil {
        b.Fatal(err)
    }
    defer f.Close()
    
    if err := pprof.StartCPUProfile(f); err != nil {
        b.Fatal(err)
    }
    defer pprof.StopCPUProfile()
    
    // Run the benchmark
    for i := 0; i < b.N; i++ {
        ComplexOperation()
    }
}

// HTTP server profiling
func enableProfiling(mux *http.ServeMux) {
    mux.HandleFunc("/debug/pprof/", pprof.Index)
    mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
    mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
    mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
    mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
    mux.HandleFunc("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
    mux.HandleFunc("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
}
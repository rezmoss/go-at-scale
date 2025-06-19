// Example 90
package main

import (
	"net/http"
	httppprof "net/http/pprof"
	"testing"
)

// ComplexOperation is just a dummy function to profile
func ComplexOperation() {
	// Some work here
	for i := 0; i < 1000000; i++ {
		_ = i * i
	}
}

func BenchmarkComplexOperation(b *testing.B) {
	// Run the benchmark
	for i := 0; i < b.N; i++ {
		ComplexOperation()
	}
}

// HTTP server profiling
func enableProfiling(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", http.HandlerFunc(httppprof.Index))
	mux.HandleFunc("/debug/pprof/cmdline", http.HandlerFunc(httppprof.Cmdline))
	mux.HandleFunc("/debug/pprof/profile", http.HandlerFunc(httppprof.Profile))
	mux.HandleFunc("/debug/pprof/symbol", http.HandlerFunc(httppprof.Symbol))
	mux.HandleFunc("/debug/pprof/trace", http.HandlerFunc(httppprof.Trace))
	mux.Handle("/debug/pprof/heap", httppprof.Handler("heap"))
	mux.Handle("/debug/pprof/goroutine", httppprof.Handler("goroutine"))
}
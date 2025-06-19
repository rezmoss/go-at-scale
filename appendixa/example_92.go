// Example 92
// benchmark_test.go
package main

import (
	"testing"
)

// BenchmarkSliceOperations tests the performance of slice operations
// with different sizes of slices
func BenchmarkSliceOperations(b *testing.B) {
	benchmarks := []struct {
		name string
		size int
	}{
		{"Small", 100},
		{"Medium", 10000},
		{"Large", 1000000},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			slice := make([]int, 0, bm.size)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				slice = append(slice, i)
			}
		})
	}
}

// Parallel benchmark
func BenchmarkParallelOperation(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Operation to benchmark
			s := make([]int, 10)
			for i := 0; i < len(s); i++ {
				s[i] = i
			}
		}
	})
}
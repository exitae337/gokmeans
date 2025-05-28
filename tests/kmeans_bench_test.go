package kmeans_test

import (
	"testing"
	"time"

	"github.com/exitae337/gokmeans/lib/kmeans"
)

func BenchmarkKmeansGo_Classic(b *testing.B) {
	benchmarkKmeansGo(b, 8, 1000, 0.001, false, 0)
}

func BenchmarkKmeansGo_KmeansPP(b *testing.B) {
	benchmarkKmeansGo(b, 8, 1000, 0.001, true, 0) // k=8
}

func BenchmarkKmeansGo_KmeansMiniBatch(b *testing.B) {
	benchmarkKmeansGo(b, 8, 1000, 0.001, true, 1000)
}

// B / op : bytes for operation
// allocs / op: allocations
func benchmarkKmeansGo(b *testing.B, k int, max_iterations int, threshold float64, kmeansPP bool, batch_size int) {
	// Засекаем общее время выполнения всех итераций
	startTime := time.Now()
	for i := 0; i < b.N; i++ {
		_, _ = kmeans.KmeansGo("../examples/clustering_datasets.xlsx", "Blobs", k, max_iterations, threshold, kmeansPP, batch_size)
	}
	total := time.Since(startTime)
	b.ReportMetric(float64(total.Milliseconds())/float64(b.N), "ms/op(total)")
	b.ReportMetric(float64(total.Seconds()), "total_sec")
}

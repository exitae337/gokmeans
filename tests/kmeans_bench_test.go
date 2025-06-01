package kmeans_test

import (
	"testing"
	"time"

	"github.com/exitae337/gokmeans/lib/kmeans"
)

func BenchmarkKmeans_Classic(b *testing.B) {
	benchmarkKmeans(b, false, 0)
}

func BenchmarkKmeans_KmeansPP(b *testing.B) {
	benchmarkKmeans(b, true, 0)
}

func BenchmarkKmeans_MiniBatch(b *testing.B) {
	benchmarkKmeans(b, true, 500)
}

func benchmarkKmeans(b *testing.B, kmeansPP bool, batchSize int) {
	path := "../examples/clustering_datasets.xlsx"
	sheet := "Blobs"
	k := 3
	maxIter := 1000
	threshold := 0.001

	// Предварительный "прогрев" (не учитывается в результатах)
	_, _ = kmeans.KmeansGo(path, sheet, k, maxIter, threshold, kmeansPP, batchSize)

	b.ResetTimer()

	// Замер общего времени выполнения всех итераций
	start := time.Now()

	for i := 0; i < b.N; i++ {
		_, _ = kmeans.KmeansGo(path, sheet, k, maxIter, threshold, kmeansPP, batchSize)
	}

	// Добавляем кастомные метрики
	elapsed := time.Since(start)
	b.ReportMetric(float64(elapsed.Nanoseconds())/float64(b.N), "ns/op(total)")
	b.ReportMetric(float64(elapsed.Seconds())/float64(b.N), "s/op")
	b.ReportMetric(float64(elapsed.Milliseconds())/float64(b.N), "ms/op")
}

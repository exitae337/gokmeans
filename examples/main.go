package main

import (
	"fmt"
	"math/rand"

	"github.com/xuri/excelize/v2"

	kmeans "github.com/exitae337/gokmeans/lib/kmeans"
	metric "github.com/exitae337/gokmeans/lib/metrics"
)

func main() {
	demoKmeans()
	demoKmeansWithInitPlus()
	demoKmeansWithInitPlusAndBatches()
}

// Classic K-means example
func demoKmeans() {
	moduleName := "GoKmeans: "
	clusters, err := kmeans.KmeansGo("clustering_datasets.xlsx", "Blobs", 3, 10000, 0.0001, false, 0)
	if err != nil {
		fmt.Println(moduleName, " : ", err)
	}
	for i, cluster := range clusters {
		fmt.Printf("Cluster %d:\n", i+1)
		fmt.Printf("Centroid: %v\n", cluster.Centroid)
	}
	fmt.Println(metric.DaviesBouldinIndex(clusters))
	fmt.Println(metric.SilhouetteScore(clusters))
}

// Kmeans with kmeans++ init example
func demoKmeansWithInitPlus() {
	moduleName := "GoKmeans: "
	clusters, err := kmeans.KmeansGo("clustering_datasets.xlsx", "Blobs", 3, 10000, 0.0001, true, 0)
	if err != nil {
		fmt.Println(moduleName, " : ", err)
	}
	for i, cluster := range clusters {
		fmt.Printf("Cluster %d:\n", i+1)
		fmt.Printf("Centroid: %v\n", cluster.Centroid)
	}
	fmt.Println(metric.DaviesBouldinIndex(clusters))
	fmt.Println(metric.SilhouetteScore(clusters))
}

// Mini-batch K-means with k-means++ example
func demoKmeansWithInitPlusAndBatches() {
	moduleName := "GoKmeans: "
	clusters, err := kmeans.KmeansGo("clustering_datasets.xlsx", "Blobs", 3, 10000, 0.0001, true, 250)
	if err != nil {
		fmt.Println(moduleName, " : ", err)
	}
	for i, cluster := range clusters {
		fmt.Printf("Cluster %d:\n", i+1)
		fmt.Printf("Centroid: %v\n", cluster.Centroid)
	}
	fmt.Println(metric.DaviesBouldinIndex(clusters))
	fmt.Println(metric.SilhouetteScore(clusters))
}

// Creating test "Example File" .xslx for testing and working example. Full random points.
func createTestFile() {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Sheet1"

	numPoints := 400 // Number of points in data for clastering
	for row := 1; row <= numPoints+1; row++ {
		x := rand.Float64() * 1000
		y := rand.Float64() * 1000
		z := rand.Float64() * 1000

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), x)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), y)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), z)
	}

	if err := f.SaveAs("points.xlsx"); err != nil {
		fmt.Println("Failed to save test file:", err)
		return
	}

	fmt.Println("File created: points.xlsx")
}

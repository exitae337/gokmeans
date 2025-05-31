package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/xuri/excelize/v2"

	kmeans "github.com/exitae337/gokmeans/lib/kmeans"
	metric "github.com/exitae337/gokmeans/lib/metrics"
)

func main() {
	fmt.Println("K-means module: Exitae337")
	demoKmeans()
}

// Classic K-means example
func demoKmeans() {
	moduleName := "GoKmeans: "
	// kmeans.KmeansGo("clustering_datasets.xlsx", "Blobs", 3, 10000, 0.0001, true, 250) - for mini-batch
	clusters, err := kmeans.KmeansGo("clustering_datasets.xlsx", "Moons", 3, 10000, 0.0001, true, 0)
	if err != nil {
		fmt.Println(moduleName, " : ", err)
	}
	for i, cluster := range clusters {
		fmt.Printf("Cluster %d:\n", i+1)
		fmt.Printf("Centroid: %v\n", cluster.Centroid)
	}
	fmt.Println("DBI", metric.DaviesBouldinIndex(clusters))
	fmt.Println("Sihoulete:", metric.SilhouetteScore(clusters))

	// ARI
	points, err := kmeans.TakePointsFromExel("clustering_datasets.xlsx", "Blobs")
	if err != nil {
		log.Panic("Failed to PARSE th file with DATA")
	}
	y_pred := metric.GetPredictedLabels(clusters, points)
	y_true, err := metric.ReadTrueLabels("clustering_datasets.xlsx", "Blobs")
	if err != nil {
		log.Panic("Failed to PARSE th file with DATA")
	}
	fmt.Println("ARI", metric.AdjustedRandIndex(y_true, y_pred))
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

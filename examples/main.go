package main

import (
	"fmt"
	"log"

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
	clusters, err := kmeans.KmeansGo("clustering_datasets.xlsx", "Circles", 2, 10000, 0.0001, true, 500)
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
	points, err := kmeans.TakePointsFromExel("clustering_datasets.xlsx", "Circles")
	if err != nil {
		log.Panic("Failed to PARSE th file with DATA")
	}
	y_pred := metric.GetPredictedLabels(clusters, points)
	y_true, err := metric.ReadTrueLabels("clustering_datasets.xlsx", "Circles")
	if err != nil {
		log.Panic("Failed to PARSE th file with DATA")
	}
	fmt.Println("ARI", metric.AdjustedRandIndex(y_true, y_pred))
}

package main

import (
	"fmt"

	gokmeans "github.com/exitae337/gokmeans/lib/kmeans"
)

func main() {
	moduleName := "GoKmeans: "
	clusters, err := gokmeans.KmeansGo("points.xlsx", "Sheet1", 4, 100, 0.001, true, 6)
	if err != nil {
		fmt.Println(moduleName, " : ", err)
	}
	for i, cluster := range clusters {
		fmt.Printf("Cluster %d:\n", i+1)
		fmt.Printf("Centroid: %v\n", cluster.Centroid)
		fmt.Printf("Points: %v\n\n", cluster.ClusterPoints)
	}
}

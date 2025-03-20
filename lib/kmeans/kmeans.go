package kmeans

import (
	"fmt"
)

// Point struct
type Point struct {
	PointValues  []float64
	MetkaKlaster int
}

// Function
func KmeansGo(pathToFile string, k, measurements int) (bool, error) {
	if measurements > 3 || measurements <= 0 {
		return false, fmt.Errorf("count of measurements must be in [1..3]")
	}
	fmt.Printf("Path to file: %v, K: %v, Measurements count for point: %v\n", pathToFile, k, measurements)
	return true, nil
}

func distanceBetween(firstPoint, secondPoint Point) float64 {
	return 0
}

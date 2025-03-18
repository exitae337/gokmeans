package kmeans

import (
	"fmt"
	"reflect"
)

// Generics for distsnce between Points function
type Point interface {
	pointForCube | pointForSquare | pointForSegment
}

// Three measurements
type pointForCube struct {
	CoordX float64
	CoordY float64
	CoordZ float64
	Metka  int
}

// Two measurements
type pointForSquare struct {
	CoordX float64
	CoordY float64
	Metka  int
}

// One measurement
type pointForSegment struct {
	CoordX float64
	Metka  int
}

// Function
func KmeansGo(pathToFile string, k, measurements int) (bool, error) {
	if measurements > 3 || measurements <= 0 {
		return false, fmt.Errorf("count of measurements must be in [1..3]")
	}
	fmt.Printf("Path to file: %v, K: %v, Measurements count for point: %v\n", pathToFile, k, measurements)
	distanceBetween(pointForCube{}, pointForCube{})
	return true, nil
}

func distanceBetween[T Point](firstPoint, secondPoint T) float64 {
	pointType := reflect.TypeOf(firstPoint)
	_ = secondPoint
	switch pointType {
	case reflect.TypeOf(pointForCube{}):
		fmt.Println("cube value")
	case reflect.TypeOf(pointForSquare{}):
		fmt.Println("square value")
	case reflect.TypeOf(pointForSegment{}):
		fmt.Println("segment value")
	default:
		fmt.Println("unknown value")
	}
	return 0
}

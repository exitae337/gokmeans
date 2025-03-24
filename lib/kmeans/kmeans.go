package kmeans

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strconv"

	"github.com/xuri/excelize/v2"
)

// Point struct
type Point struct {
	PointValues       []float64
	DistanceToKlaster float64
	MetkaKlaster      int
}

// Klaster struct
type Klaster struct {
	KlasterValues      []float64
	PointsKlasterArray []Point
	KlasterNumber      int
}

// Distance
type Distance struct {
	PointPointer *Point
	Distance     float64
	Klaster      *Klaster
}

// Function. Return -> map[int]float64, error
func KmeansGo(pathToFile, sheetName string, k, measurements int) (map[*Klaster][]Point, error) {
	var pointsArray []Point
	var klastersArray []Klaster
	// Working with Excel file
	// Open file
	file, err := excelize.OpenFile(pathToFile)
	if err != nil {
		return nil, err
	}
	// Reading and working with rows
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		point := Point{}
		for _, colCel := range row {
			pointKoord, err := strconv.ParseFloat(colCel, 64)
			if err != nil {
				return nil, err
			}
			point.PointValues = append(point.PointValues, pointKoord)
		}
		pointsArray = append(pointsArray, point)
	}

	// Length of PointsArray
	N := len(pointsArray)

	// Init Klasters for kmeans
	for i := 0; i < k; i++ {
		n := rand.IntN(N - 1)
		klastersArray = append(klastersArray, Klaster{
			KlasterValues: pointsArray[n].PointValues,
			KlasterNumber: i,
		})
	}
	// Init distance map
	distances := make(map[int][]Distance)

	// Count distances from Klaster to all Points
	for _, klaster := range klastersArray {
		distances[klaster.KlasterNumber] = distanceBetween(&klaster, &pointsArray)
	}

	minDistances := make([]Distance, N)

	// Points to Klasters
	for _, distToPoints := range distances {
		for i, val := range distToPoints {
			if minDistances[i].Distance > val.Distance {
				minDistances[i] = val
			}
		}
	}

	// Points to -> Klaster(PointsArray)
	for _, klaster := range klastersArray {
		for _, dist := range minDistances {
			if klaster.KlasterNumber == dist.Klaster.KlasterNumber {
				klaster.PointsKlasterArray = append(klaster.PointsKlasterArray, *dist.PointPointer)
			}
		}
	}

	for _, klaster := range klastersArray {
		for _, point := range klaster.PointsKlasterArray {
			fmt.Printf("Klaster Number: %v, Point: %v", klaster.KlasterNumber, point)
		}
	}

	// TODO
	return nil, nil
}

// Euclidean distance
func distanceBetween(klaster *Klaster, points *[]Point) []Distance {
	distances := make([]Distance, len(klaster.KlasterValues))
	for _, point := range *points {
		sumKvadrCoord := 0.0
		for i, coord := range point.PointValues {
			sumKvadrCoord += math.Pow(klaster.KlasterValues[i]-coord, 2)
		}
		distances = append(distances, Distance{
			Klaster:      klaster,
			Distance:     math.Sqrt(sumKvadrCoord),
			PointPointer: &point,
		})
	}
	return distances
}

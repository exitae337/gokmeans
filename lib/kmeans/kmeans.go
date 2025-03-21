package kmeans

import (
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

// Function. Return -> map[int]float64, error
func KmeansGo(pathToFile, sheetName string, k, measurements int) (map[int]float64, error) {
	var pointsArray []Point
	var klastersArray []Klaster
	// Working with Excel file
	// Open file
	file, err := excelize.OpenFile(pathToFile)
	if err != nil {
		return nil, err
	}
	// Reading and working rows
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
	// Init Klasters for kmeans
	for i := 0; i < k; i++ {
		n := rand.IntN(len(pointsArray) - 1)
		klastersArray = append(klastersArray, Klaster{
			KlasterValues: pointsArray[n].PointValues,
			KlasterNumber: i,
		})
	}
	// Init distance map
	distances := make(map[int][]float64)
	// Count distances from Klaster to all Points
	for _, klaster := range klastersArray {
		distances[klaster.KlasterNumber] = distanceBetween(&klaster, &pointsArray)
	}

	// TODO
	return nil, nil
}

// Euclidean distance
func distanceBetween(klaster *Klaster, points *[]Point) []float64 {
	distances := make([]float64, len(klaster.KlasterValues))
	for _, point := range *points {
		sumKvadrCoord := 0.0
		for i, coord := range point.PointValues {
			sumKvadrCoord += math.Pow(klaster.KlasterValues[i]-coord, 2)
		}
		distances = append(distances, math.Sqrt(sumKvadrCoord))
	}
	return distances
}

package kmeans

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

// Point type
type Point []float64

// Cluster struct
type Cluster struct {
	Centroid      Point
	ClusterPoints []Point
}

// Main function. Return -> []Cluster, error
func KmeansGo(pathToFile, sheetName string, k, maxIterations int, threshold float64, kmeans_plus bool, batchSize int) ([]Cluster, error) {
	points, err := takePointsFromExel(pathToFile, sheetName)
	if err != nil {
		return nil, err
	}
	// Algorythm
	if batchSize > 0 && batchSize < len(points) {
		return miniBatchKmeans(points, k, batchSize, maxIterations, threshold)
	}
	return ClassicKMeans(points, k, maxIterations, kmeans_plus, threshold)
}

// Classic K-means
func ClassicKMeans(points []Point, k int, maxIterations int, kmeans_plus bool, threshold float64) ([]Cluster, error) {
	if k <= 0 || len(points) <= k {
		return nil, fmt.Errorf("value of 'k' parameter is invalid: zero or bigger than points count -> k=%d", k)
	}

	var centroids []Point

	// K-maens++ or K-means
	if kmeans_plus {
		centroids = centroidsInitPP(points, k)
	} else {
		centroids = centroidsInit(points, k)
	}

	var clusters []Cluster

	for i := 0; i < maxIterations; i++ {
		clusters = assignPoints(points, centroids)
		newCentroids := updateCenrtoids(clusters)

		if !centroidsChanged(centroids, newCentroids, threshold) {
			break
		}

		centroids = newCentroids
	}
	return clusters, nil
}

// Take Points From exel file
func takePointsFromExel(pathToFile, sheetName string) ([]Point, error) {
	// Working with Excel file
	// Open file
	f, err := excelize.OpenFile(pathToFile)
	if err != nil {
		return nil, err
	}
	// Close file
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Reading and working with rows (Current Points array from exel file)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	currentPoints := []Point{}
	for _, row := range rows {
		currentPoint := Point{}
		for _, colCell := range row {
			// Convert to float64
			floatValue, err := strconv.ParseFloat(colCell, 64)
			if err != nil {
				return nil, err
			}
			currentPoint = append(currentPoint, floatValue)
		}
		currentPoints = append(currentPoints, currentPoint)
	}
	return currentPoints, nil
}

// Mini-Batch K-means (with batches)
func miniBatchKmeans(points []Point, k int, batchSize int, maxIterations int, threshold float64) ([]Cluster, error) {
	if k <= 0 || k >= len(points) {
		return nil, fmt.Errorf("value of 'k' parameter is invalid: zero or bigger than points count -> k=%d", k)
	}
	randSeed := rand.NewSource(time.Now().UnixNano())
	centroids := centroidsInitPP(points, k)
	randForMiniBatch := rand.New(randSeed)

	var clusters []Cluster
	for i := 0; i < maxIterations; i++ {
		batch := make([]Point, batchSize)
		perm := randForMiniBatch.Perm(len(points))[:batchSize]
		for j, idx := range perm {
			batch[j] = points[idx]
		}
		clusters = assignPoints(batch, centroids)
		newCentroids := updateCentroidsWithMiniBatch(clusters, centroids, 1.0/(float64(i+1)))
		if !centroidsChanged(centroids, newCentroids, threshold) {
			break
		}
		centroids = newCentroids
	}
	return assignPoints(points, centroids), nil
}

// Update for Mini-Batch
func updateCentroidsWithMiniBatch(clusters []Cluster, oldCentroids []Point, learningRate float64) []Point {
	newCentroids := make([]Point, len(clusters))
	for i, cluster := range clusters {
		newCentroids[i] = make(Point, len(cluster.Centroid))
		copy(newCentroids[i], oldCentroids[i])
		if len(cluster.ClusterPoints) == 0 {
			continue
		}
		// Middle for mini-batch
		batchMean := make(Point, len(cluster.Centroid))
		for _, p := range cluster.ClusterPoints {
			for j := range p {
				batchMean[j] += p[j]
			}
		}
		for j := range batchMean {
			batchMean[j] /= float64(len(cluster.ClusterPoints))
		}

		for j := range newCentroids[i] {
			newCentroids[i][j] = (1-learningRate)*newCentroids[i][j] + learningRate*batchMean[j]
		}
	}
	return newCentroids
}

// Centroids Init -> random choice
func centroidsInit(points []Point, k int) []Point {
	seedInit := rand.NewSource(time.Now().UnixNano())
	randInit := rand.New(seedInit)
	centroids := make([]Point, k)
	perm := randInit.Perm(len(points))[:k]
	for i, idx := range perm {
		centroids[i] = make(Point, len(points[idx]))
		copy(centroids[i], points[idx])
	}
	return centroids
}

// Centroids init PP
func centroidsInitPP(points []Point, k int) []Point {
	seedInit := rand.NewSource(time.Now().UnixNano())
	randInit := rand.New(seedInit)

	centroids := make([]Point, k)

	firstIdx := randInit.Intn(len(points))
	centroids[0] = make(Point, len(points[firstIdx]))
	copy(centroids[0], points[firstIdx])

	for i := 1; i < k; i++ {
		distances := make([]float64, len(points))
		sum := 0.0

		for j, p := range points {
			minDist := math.MaxFloat64
			for _, c := range centroids[:i] {
				if c != nil {
					dist := p.distanceBetween(c)
					if dist < minDist {
						minDist = dist
					}
				}
			}
			distances[j] = minDist * minDist
			sum += distances[j]
		}

		r := randInit.Float64() * sum
		cumSum := 0.0
		selectedIdx := 0

		for j, d := range distances {
			cumSum += d
			if cumSum >= r {
				selectedIdx = j
				break
			}
		}

		centroids[i] = make(Point, len(points[selectedIdx]))
		copy(centroids[i], points[selectedIdx])
	}

	return centroids
}

// Assign Points to Klasters
func assignPoints(points []Point, centroids []Point) []Cluster {
	clusters := make([]Cluster, len(centroids))
	for i := range clusters {
		clusters[i].Centroid = make(Point, len(centroids[i]))
		copy(clusters[i].Centroid, centroids[i])
	}

	for _, p := range points {
		minDistance := math.MaxFloat64
		clusterIdx := 0

		for i, c := range centroids {
			currentDistance := p.distanceBetween(c)
			if currentDistance < minDistance {
				minDistance = currentDistance
				clusterIdx = i
			}
		}

		clusters[clusterIdx].ClusterPoints = append(clusters[clusterIdx].ClusterPoints, p)
	}
	return clusters
}

// Update Centroids (Classic K-means)
func updateCenrtoids(clusters []Cluster) []Point {
	newCentroids := make([]Point, len(clusters))
	for i, cluster := range clusters {
		if len(cluster.ClusterPoints) == 0 {
			newCentroids[i] = make(Point, len(cluster.ClusterPoints))
			copy(newCentroids[i], cluster.Centroid)
			continue
		}
		newCentroid := make(Point, len(cluster.Centroid))
		for _, p := range cluster.ClusterPoints {
			for j := range p {
				newCentroid[j] += p[j]
			}
		}
		for j := range newCentroid {
			newCentroid[j] /= float64(len(cluster.ClusterPoints))
		}
		newCentroids[i] = newCentroid
	}
	return newCentroids
}

// Convergence check
func centroidsChanged(oldCentroids, newCentroids []Point, threshold float64) bool {
	if len(oldCentroids) != len(newCentroids) {
		return true
	}
	for i := range oldCentroids {
		if oldCentroids[i].distanceBetween(newCentroids[i]) > threshold {
			return true
		}
	}
	return false
}

// Euclidean distance
func (p Point) distanceBetween(other Point) float64 {
	sum := 0.0
	for i := range p {
		diff := p[i] - other[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

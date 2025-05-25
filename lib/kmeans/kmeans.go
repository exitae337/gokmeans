// Package gokmeans implements a convenient tool for clustering
// Author: Chernyshev Maxim <exitae337@gmail.com>
// License: Apache License, Version 2.0
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
	return classicKMeans(points, k, maxIterations, kmeans_plus, threshold)
}

// Classic K-means. Return -> Cluster[], error
func classicKMeans(points []Point, k int, maxIterations int, kmeans_plus bool, threshold float64) ([]Cluster, error) {
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

// Mini-batch K-means (K-means++ init). Return -> Cluster[], error
func miniBatchKmeans(points []Point, k int, batchSize int, maxIterations int, threshold float64) ([]Cluster, error) {
	if k <= 0 || k >= len(points) {
		return nil, fmt.Errorf("value of 'k' parameter is invalid: zero or bigger than points count -> k=%d", k)
	}

	randSeed := rand.NewSource(time.Now().UnixNano())
	randForMiniBatch := rand.New(randSeed)
	centroids := centroidsInitPP(points, k)

	clusterCounts := make([]int, k) // points count in clusters

	for i := 0; i < maxIterations; i++ {

		batch := make([]Point, batchSize)
		perm := randForMiniBatch.Perm(len(points))[:batchSize]
		for j, idx := range perm {
			batch[j] = points[idx]
		}

		clusters := assignPoints(batch, centroids) // batch points to clusters

		newCentroids := make([]Point, len(centroids))
		for j := range newCentroids {
			newCentroids[j] = make(Point, len(centroids[j]))
			copy(newCentroids[j], centroids[j])

			if len(clusters[j].ClusterPoints) == 0 {
				continue
			}

			// Middle
			batchMean := make(Point, len(clusters[j].Centroid))
			for _, p := range clusters[j].ClusterPoints {
				for dim := range p {
					batchMean[dim] += p[dim]
				}
			}
			for dim := range batchMean {
				batchMean[dim] /= float64(len(clusters[j].ClusterPoints))
			}

			n := float64(clusterCounts[j])               // before
			m := float64(len(clusters[j].ClusterPoints)) // after
			for dim := range newCentroids[j] {
				newCentroids[j][dim] = (n*newCentroids[j][dim] + m*batchMean[dim]) / (n + m)
			}

			// Обновляем счетчик точек в кластере
			clusterCounts[j] += len(clusters[j].ClusterPoints)
		}

		// 4. Проверяем сходимость
		if !centroidsChanged(centroids, newCentroids, threshold) {
			break
		}
		centroids = newCentroids
	}

	// Возвращаем финальные кластеры для всех точек
	return assignPoints(points, centroids), nil
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

// Centroids init PP -> Kmeans++ init
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
					// MinDist to centroids
					dist := p.DistanceBetween(c)
					if dist < minDist {
						minDist = dist
					}
				}
			}
			distances[j] = minDist * minDist // D(x)^2
			sum += distances[j]              // for probability choice
		}

		// Normalize distances into probabilities
		probs := make([]float64, len(points))
		cumSum := 0.0
		for j, d := range distances {
			cumSum += d / sum
			probs[j] = cumSum
		}

		r := randInit.Float64()
		selectedIdx := 0
		for j, prob := range probs {
			if r <= prob {
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
			currentDistance := p.DistanceBetween(c)
			if currentDistance < minDistance {
				minDistance = currentDistance
				clusterIdx = i
			}
		}

		clusters[clusterIdx].ClusterPoints = append(clusters[clusterIdx].ClusterPoints, p)
	}
	return clusters
}

// Update Centroids (Classic or K-means++)
func updateCenrtoids(clusters []Cluster) []Point {
	newCentroids := make([]Point, len(clusters))
	for i, cluster := range clusters {
		if len(cluster.ClusterPoints) == 0 {
			newCentroids[i] = make(Point, len(cluster.Centroid))
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

// Convergence check: > threshold
func centroidsChanged(oldCentroids, newCentroids []Point, threshold float64) bool {
	if len(oldCentroids) != len(newCentroids) {
		return true
	}
	for i := range oldCentroids {
		if oldCentroids[i].DistanceBetween(newCentroids[i]) > threshold {
			return true
		}
	}
	return false
}

// Helper: Euclidean distance
func (p Point) DistanceBetween(other Point) float64 {
	sum := 0.0
	for i := range p {
		diff := p[i] - other[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

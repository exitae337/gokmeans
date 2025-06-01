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
	points, err := TakePointsFromExel(pathToFile, sheetName)
	if err != nil {
		return nil, err
	}
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
func TakePointsFromExel(pathToFile, sheetName string) ([]Point, error) {
	f, err := excelize.OpenFile(pathToFile)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}
	currentPoints := []Point{}
	for _, row := range rows {
		currentPoint := Point{}
		for i, colCell := range row {
			if i == len(row)-1 {
				break
			}
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

func miniBatchKmeans(points []Point, k int, batchSize int, maxIterations int, threshold float64) ([]Cluster, error) {
	if k <= 0 || k >= len(points) {
		return nil, fmt.Errorf("invalid k value: %d", k)
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	centroids := centroidsInitPP(points, k)

	clusterCounts := make([]int, k)
	sums := make([]Point, k)
	counts := make([]int, k)

	for i := 0; i < maxIterations; i++ {

		start := rnd.Intn(len(points) - batchSize)
		batch := points[start : start+batchSize]

		for i := range sums {
			if sums[i] == nil {
				sums[i] = make(Point, len(centroids[i]))
			} else {
				for j := range sums[i] {
					sums[i][j] = 0
				}
			}
			counts[i] = 0
		}

		for _, p := range batch {
			minDist := math.MaxFloat64
			closest := 0

			for i, c := range centroids {
				dist := p.DistanceBetween(c)
				if dist < minDist {
					minDist = dist
					closest = i
				}
			}

			for dim := range p {
				sums[closest][dim] += p[dim]
			}
			counts[closest]++
		}

		changed := false
		for j := range centroids {
			if counts[j] == 0 {
				continue
			}

			newCentroid := make(Point, len(centroids[j]))
			for dim := range newCentroid {
				batchMean := sums[j][dim] / float64(counts[j])
				total := float64(clusterCounts[j])
				newCentroid[dim] = (total*centroids[j][dim] + float64(counts[j])*batchMean) / (total + float64(counts[j]))
			}

			if centroids[j].DistanceBetween(newCentroid) > threshold {
				changed = true
			}
			centroids[j] = newCentroid
			clusterCounts[j] += counts[j]
		}

		if !changed {
			break
		}
	}

	clusters := make([]Cluster, k)
	for i := range clusters {
		clusters[i].Centroid = centroids[i]
		clusters[i].ClusterPoints = make([]Point, 0)
	}

	for _, p := range points {
		minDist := math.MaxFloat64
		closest := 0
		for i, c := range centroids {
			dist := p.DistanceBetween(c)
			if dist < minDist {
				minDist = dist
				closest = i
			}
		}
		clusters[closest].ClusterPoints = append(clusters[closest].ClusterPoints, p)
	}

	return clusters, nil
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
					dist := p.DistanceBetween(c)
					if dist < minDist {
						minDist = dist
					}
				}
			}
			distances[j] = minDist * minDist
			sum += distances[j]
		}

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

// Copy point for speed
func (p Point) Copy() Point {
	newPoint := make(Point, len(p))
	copy(newPoint, p)
	return newPoint
}

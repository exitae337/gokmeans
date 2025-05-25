package metric

import (
	"math"

	"github.com/exitae337/gokmeans/lib/kmeans"
)

// Metrics: DBI & Sihoulette Score

// DaviesBouldinIndex
func DaviesBouldinIndex(clusters []kmeans.Cluster) float64 {
	k := len(clusters)
	if k <= 1 {
		return 0.0
	}

	s := make([]float64, k)
	for i, cluster := range clusters {
		sumDist := 0.0
		for _, point := range cluster.ClusterPoints {
			sumDist += point.DistanceBetween(cluster.Centroid)
		}
		if len(cluster.ClusterPoints) > 0 {
			s[i] = sumDist / float64(len(cluster.ClusterPoints))
		}
	}

	dbi := 0.0
	for i := 0; i < k; i++ {
		maxRatio := -1.0
		for j := 0; j < k; j++ {
			if i == j {
				continue
			}
			distance := clusters[i].Centroid.DistanceBetween(clusters[j].Centroid)
			ratio := (s[i] + s[j]) / distance
			if ratio > maxRatio {
				maxRatio = ratio
			}
		}
		dbi += maxRatio
	}
	dbi /= float64(k)

	return dbi
}

// SilhouetteScore
func SilhouetteScore(clusters []kmeans.Cluster) float64 {
	var allPoints []kmeans.Point
	labels := make([]int, 0)

	for clusterID, cluster := range clusters {
		for _, point := range cluster.ClusterPoints {
			allPoints = append(allPoints, point)
			labels = append(labels, clusterID)
		}
	}

	n := len(allPoints)
	if n <= 1 {
		return 0.0
	}

	if len(labels) != n {
		return 0.0
	}

	silScores := make([]float64, n)
	for i := 0; i < n; i++ {
		if i >= len(labels) {
			continue
		}

		cluster := labels[i]
		a := averageIntraClusterDistance(allPoints, labels, i, cluster)
		minB := minInterClusterDistance(allPoints, labels, i, cluster)

		denominator := math.Max(a, minB)
		if denominator == 0 {
			silScores[i] = 0
		} else {
			silScores[i] = (minB - a) / denominator
		}
	}

	sum := 0.0
	for _, s := range silScores {
		sum += s
	}
	return sum / float64(n)
}

// In cluster distance
func averageIntraClusterDistance(points []kmeans.Point, labels []int, pointIdx int, cluster int) float64 {
	if len(points) == 0 || len(labels) == 0 || pointIdx >= len(points) || pointIdx >= len(labels) {
		return 0.0
	}

	sum := 0.0
	count := 0
	for i, label := range labels {
		if i >= len(points) {
			continue
		}
		if label == cluster && i != pointIdx {
			dist := points[pointIdx].DistanceBetween(points[i])
			sum += dist
			count++
		}
	}

	if count == 0 {
		return 0.0
	}
	return sum / float64(count)
}

// Min to other cluster distance
func minInterClusterDistance(points []kmeans.Point, labels []int, pointIdx int, currentCluster int) float64 {
	if pointIdx >= len(points) {
		return 0.0
	}

	uniqueClusters := make(map[int]bool)
	for _, label := range labels {
		if label != currentCluster {
			uniqueClusters[label] = true
		}
	}

	minB := math.MaxFloat64
	for cluster := range uniqueClusters {
		sum := 0.0
		count := 0
		for i, label := range labels {
			if label == cluster {
				dist := points[pointIdx].DistanceBetween(points[i])
				sum += dist
				count++
			}
		}
		if count > 0 {
			avgDist := sum / float64(count)
			if avgDist < minB {
				minB = avgDist
			}
		}
	}

	if minB == math.MaxFloat64 {
		return 0.0
	}
	return minB
}

# Gokmeans library
K-means library on *golang*.

This module implements popular K-means variants of efficient clustering:

- **K-means**: Classic Lloyds's algorythm partitioning data into k clustersby minimizing within-cluster variance
- **K-maens++**: Smart centroid initialization for faster convergence and better accurancy
- **Mini-Batch K-means**: Optimized for large datasets using small random *batches*, trading precision for speed 

*Designed for scalability and simplicity*
**Contributions are welcome!**

# Module installation

``` PowerShell
go get github.com/exitae337/gokmeans 
```

# Using

First, you need to prepare an *excel table* with points and their **measurements** (coordinates)

The main function of this module called **KmeansGo**:
``` Go
gokmeans.KmeansGo("name.xlsx", "sheet_name", k, count_iter, threshold, bool_kmeans_init, batch_size)
```

*Where:*
- **"name.xslx"**: file with prepared data (points with coordinates)
- **"sheet name"**: xslx **sheet** name
- **k**: number of clusters (**K**-means)
- **count_iter**: number of algorithm iterations
- **threshold**: parameter that determines when the K-Means algorithm should terminate (algorithm convergence)
- **bool_kmeans_init**: K-means++ initialization
- **batch_size**: number of points in *mini-batch* (Mini-batch K-means)

The output of the algorithm is a slice of the formed clusters.

**Cluster** and **Point** type:

``` Go
type Cluster struct {
	Centroid Point
	ClusterPoints []Point
}

type Point []float64
```


**Classic K-means:**

For Classic K-means:
- bool_kmeans_init: **False**
- batch_size: **Zero or below**

*Example:*

``` Go
// Importing
import (
	gokmeans "github.com/exitae337/gokmeans/lib/kmeans"
)

func main() {
    // Using
    clusters, err := gokmeans.KmeansGo("points.xlsx", "Sheet1", 4, 100, 0.001, false, 0)
	// Errors handling
    if err != nil {
		// Logic of working with error
	}
	// View
	for i, cluster := range clusters {
		fmt.Printf("Cluster %d:\n", i+1)
		fmt.Printf("Centroid: %v\n", cluster.Centroid) // Centroid
		fmt.Printf("Points: %v\n\n", cluster.ClusterPoints) // Points of Cluster
	}
}
```

**K-means++ initialization:**

For K-means++ (only):
- bool_kmeans_init: **False**
- batch_size: **Zero or below**

*Example:*

``` Go
// Importing
import (
	gokmeans "github.com/exitae337/gokmeans/lib/kmeans"
)

func main() {
    // Using
    clusters, err := gokmeans.KmeansGo("points.xlsx", "Sheet1", 4, 100, 0.001, false, 0)
	// Errors handling
    if err != nil {
		// Logic of working with error
	}
}
```

**Mini-batch initialization:**

For Mini-batch (only):
- bool_kmeans_init: **True**
- batch_size: **0 <= size < points_count**

*Example:*

``` Go
// Importing
import (
	gokmeans "github.com/exitae337/gokmeans/lib/kmeans"
)

func main() {
    // Using
    clusters, err := gokmeans.KmeansGo("points.xlsx", "Sheet1", 4, 100, 0.001, false, 0)
	// Errors handling
    if err != nil {
		// Logic of working with error
	}
}
```
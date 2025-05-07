# Gokmeans module
K-means library on *golang*.

This module implements popular K-means variants of efficient clustering:

- **K-means**: Classic Lloyds's algorythm partitioning data into k clustersby minimizing within-cluster variance
- **K-means++**: Smart centroid initialization for faster convergence and better accurancy
- **Mini-Batch K-means**: Optimized for large datasets using small random *batches*, trading precision for speed 

*Designed for scalability and simplicity*

# Module installation

``` PowerShell
go get github.com/exitae337/gokmeans 
```

# Using

First, you need to prepare an *excel table* with points and their **measurements** (coordinates)

*For example*

**xlsx** table:

```
       A       B       C
1  0.123   0.456   0.786 -> Point's coordinates (three measurements for example)
```



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

**Cluster** and **Point** struct type:

``` Go
type Cluster struct {
	Centroid Point
	ClusterPoints []Point
}

type Point []float64
```

Main function have helper func for testing (creating .xlsx file with data set):

- **numPoints** - number of points

``` Go
func createTestFile() {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Sheet1"

	numPoints := 100000
	for row := 2; row <= numPoints+1; row++ {
		x := rand.Float64() * 1000
		y := rand.Float64() * 1000
		z := rand.Float64() * 1000

		// Записываем в колонки A, B, C (X, Y, Z)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), x)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), y)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), z)
	}

	// Сохраняем файл
	if err := f.SaveAs("points.xlsx"); err != nil {
		fmt.Println("Failed to save test file:", err)
		return
	}

	fmt.Println("File created: points.xlsx")
}
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
- bool_kmeans_init: **True**
- batch_size: **Zero or below**

*Example:*

``` Go
// Importing
import (
	gokmeans "github.com/exitae337/gokmeans/lib/kmeans"
)

func main() {
    // Using
    clusters, err := gokmeans.KmeansGo("points.xlsx", "Sheet1", 4, 100, 0.001, true, 0) // Mini-batch size = 0
	// Errors handling
    if err != nil {
		// Logic of working with error
	}
}
```

**Mini-batch initialization:**

For Mini-batch (only):
- bool_kmeans_init: **False**
- batch_size: **0 <= size < points_count**

*Example:*

``` Go
// Importing
import (
	gokmeans "github.com/exitae337/gokmeans/lib/kmeans"
)

func main() {
    // Using
    clusters, err := gokmeans.KmeansGo("points.xlsx", "Sheet1", 4, 100, 0.001, false, 6) // No K-means++ init
	// Errors handling
    if err != nil {
		// Logic of working with error
	}
}
```

But you can **combine** initializations (k-means++ & mini-batch):

```Go
	clusters, err := gokmeans.KmeansGo("points.xlsx", "Sheet1", 4, 100, 0.001, true, 6) // K-means++ init && Mini-batch init
	// Errors handling
    if err != nil {
		// Logic of working with error
	}
```

🔧 Contributions welcome! Report issues or submit PRs.
📜 License: Apache License, Version 2.0

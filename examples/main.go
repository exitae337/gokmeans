package main

import (
	"fmt"
	"math/rand"

	"github.com/xuri/excelize/v2"

	gokmeans "github.com/exitae337/gokmeans/lib/kmeans"
)

func main() {
	createTestFile()
	demoKmeans()
}

func demoKmeans() {
	moduleName := "GoKmeans: "
	clusters, err := gokmeans.KmeansGo("points.xlsx", "Sheet1", 8, 100, 0.001, false, 0)
	if err != nil {
		fmt.Println(moduleName, " : ", err)
	}
	for i, cluster := range clusters {
		fmt.Printf("Cluster %d:\n", i+1)
		fmt.Printf("Centroid: %v\n", cluster.Centroid)
		fmt.Printf("Points: %v\n\n", cluster.ClusterPoints)
	}
}

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

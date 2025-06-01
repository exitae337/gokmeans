package datagen

import (
	"fmt"
	"math/rand"

	"github.com/xuri/excelize/v2"
)

// Creating test "Example File" .xslx for testing and working example. Full random points.
func CreateTestFile() {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Sheet1"

	numPoints := 400 // Number of points in data for clastering
	for row := 1; row <= numPoints+1; row++ {
		x := rand.Float64() * 1000
		y := rand.Float64() * 1000
		z := rand.Float64() * 1000

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), x)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), y)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), z)
	}

	if err := f.SaveAs("points.xlsx"); err != nil {
		fmt.Println("Failed to save test file:", err)
		return
	}

	fmt.Println("File created: points.xlsx")
}

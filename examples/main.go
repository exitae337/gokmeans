package main

import (
	"fmt"
	"log"

	gokmeans "github.com/exitae337/gokmeans/lib/kmeans"
)

func main() {
	moduleName := "GoKmeans: "
	if ok, err := gokmeans.KmeansGo("points.xlsx", "Sheet1", 4, 100, 0.001); err != nil {
		log.Fatal(moduleName, err)
	} else {
		fmt.Println(ok)
	}
}

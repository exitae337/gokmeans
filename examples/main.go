package main

import (
	"fmt"
	"log"

	gokmeans "github.com/exitae337/gokmeans/lib/kmeans"
)

func main() {
	moduleName := "GoKmeans: "
	if ok, err := gokmeans.KmeansGo("../../storage/file.exel", 6, 4); err != nil {
		log.Fatal(moduleName, err)
	} else {
		fmt.Println(ok)
	}
}

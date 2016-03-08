package main

import (
	"fmt"
	"os"

	"github.com/czxichen/AutoWork/tools/zip"
)

func main() {
	File, _ := os.Create("1.zip")
	defer File.Close()
	tgz := zip.NewZipWriter(File)
	defer tgz.Close()

	err := tgz.Walk("./server/")
	fmt.Println(err)
}

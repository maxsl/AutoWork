package main

import (
	"fmt"

	"github.com/czxichen/AutoWork/tools/zip"
)

func main() {
	zip.Gzip("./", "1.tar.gz", []string{"*.go"}, fmt.Printf)
}

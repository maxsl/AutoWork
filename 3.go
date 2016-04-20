package main

import (
	"os"
)

func main() {
	F, _ := os.Create("ce")
	F.Write([]byte{0})
	F.Close()
}

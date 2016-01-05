package main

import (
	"fmt"
	"net/http"
)

func main() {
	err := http.ListenAndServe("0.0.0.0:80", http.FileServer(http.Dir(".")))
	if err != nil {
		fmt.Println(err)
	}
}

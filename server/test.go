package main

import (
	"net/http"

	"github.com/AutoWork/server/route"
)

func main() {
	http.HandleFunc("/", route.Router)

	http.ListenAndServe(":1789", nil)
}

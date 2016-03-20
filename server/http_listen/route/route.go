package route

import (
	"net/http"
)

func Router(w http.ResponseWriter, r *http.Request) {
	if !auth(r) {
		http.Error(w, "Authentication failed!", StatusAuthFaild)
		return
	}
	switch r.URL.Path {
	case "/":
		Index(w)
	case "/receive":
		Receive(w, r)
	case "/run":
		Run(w, r)
	case "/commond":
		Commond(w, r)
	default:
		http.Error(w, "Page NotFound.", http.StatusNotFound)
	}
}

func auth(r *http.Request) bool {
	return true
}

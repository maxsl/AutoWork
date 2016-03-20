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
		index(w)
	case "/receive":
		receive(w, r)
	case "/run":
		run(w, r)
	case "/commond":
		commond(w, r)
	default:
		http.Error(w, "Page NotFound.", http.StatusNotFound)
	}
}

func auth(r *http.Request) bool {
	return true
}

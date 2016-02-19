package route

import (
	"net/http"
)

func Router(w http.ResponseWriter, r *http.Request) {
	if !auth(r) {
		http.Error(w, "Authentication failed!", 600)
		return
	}
	switch r.URL.Path {
	case "/":
		index(w)
	case "/commond":
		commond(w, r)
	default:
		http.Error(w, "Page NotFound.", 404)
	}

}

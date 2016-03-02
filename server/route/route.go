package route

import (
	"net/http"

	"github.com/AutoWork/server/route/handler"
)

func Router(w http.ResponseWriter, r *http.Request) {
	if !auth(r) {
		http.Error(w, "Authentication failed!", 600)
		return
	}
	switch r.URL.Path {
	case "/":
		handler.Index(w)
	case "/commond":
		handler.Commond(w, r)
	default:
		http.Error(w, "Page NotFound.", 404)
	}
}

package route

import (
	"net/http"
)

func index(w http.ResponseWriter) {
	w.Write([]byte("Welcome to AutoWork!"))
}

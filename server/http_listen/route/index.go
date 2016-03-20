package route

import (
	"net/http"
)

func Index(w http.ResponseWriter) {
	w.Write([]byte("Welcome to AutoWork!"))
}

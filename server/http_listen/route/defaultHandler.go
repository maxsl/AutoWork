package route

import (
	"net/http"
)

func Default(w http.ResponseWriter, errCode int) {
	http.Error(w, StatusText(errCode), errCode)
}

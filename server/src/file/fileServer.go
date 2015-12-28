package file

import (
	"alltype"
	"net/http"
)

func HttpServer() {
	url()
	http.ListenAndServe(":1789", nil)
}

func url() {
	http.HandleFunc("/file/", route)
}

func route(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.URL.Path == "/file/" {
		http.Error(w, "Request error.", 404)
		return
	}
	k, ok := r.URL.Query()["key"]
	if ok && len(k) == 1 {
		if _, ok := alltype.AgentMap[k[0]]; ok {
			http.FileServer(http.Dir("./")).ServeHTTP(w, r)
			return
		}
	}
	http.Error(w, "Request error.", 404)
}

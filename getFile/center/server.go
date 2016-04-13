package center

import (
	"net/http"
	"strings"
)

func Server() {
	Init()
	http.HandleFunc("/", route)
	http.ListenAndServe(config.ListenIP, nil)
}

func route(w http.ResponseWriter, r *http.Request) {
	Log.Println("RemoteAddr:", r.RemoteAddr, "RequestURI:", r.RequestURI)
	if !strings.Contains(config.WhiteList, strings.Split(r.RemoteAddr, ":")[0]) {
		http.NotFound(w, r)
		return
	}
	switch r.URL.Path {
	case "/getFile/index":
		index(w, r)
	case "/getFile/run":
		run(w, r)
	case "/getFile/finished":
		finished(w, r)
	case "/getFile/download":
		download(w, r)
	default:
		return
	}
}

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
	if !strings.Contains(config.WhiteList, strings.Split(r.RemoteAddr, ":")[0]) {
		Log.Printf("%s Not in white list\n", r.RemoteAddr)
		http.Error(w, "Not in white list", 403)
		return
	}
	Log.Println("RemoteAddr:", r.RemoteAddr, "RequestURI:", r.RequestURI)
	switch r.URL.Path {
	case "/getFile/index":
		index(w, r)
	case "/getFile/run":
		run(w, r)
	case "/getFile/finished":
		finished(w, r)
	case "/getFile/download":
		download(w, r)
	case "/getFile/contacts":
		contacts(w, r)
	default:
		return
	}
}

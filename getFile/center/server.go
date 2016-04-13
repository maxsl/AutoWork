package center

import "net/http"

func Server() {
	Init()
	http.HandleFunc("/", route)
	http.ListenAndServe(config.ListenIP, nil)
}

func route(w http.ResponseWriter, r *http.Request) {
	Log.Println("RemoteAddr:", r.RemoteAddr, "RequestURI:", r.RequestURI)
	switch r.URL.Path {
	case "/getFile/index":
		index(w, r)
	case "/getFile/run":
		run(w, r)
	case "/getFile/finished":
		finished(w, r)
	default:
		return
	}
}

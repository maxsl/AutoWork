package getFile

import "net/http"

func Server() {
	Init()
	go exec_run()
	http.HandleFunc("/", route)
	err := http.ListenAndServe(config.ListenIP+":"+config.Port, nil)
	if err != nil {
		Log.Printf("Listen server faild,%v\n", err)
	}
}

func route(w http.ResponseWriter, r *http.Request) {
	if config.Debug {
		Log.Println("RemoteAddr:", r.RemoteAddr, "URI:", r.RequestURI)
	}
	if !ipIsLanIP(r.RemoteAddr) {
		http.Error(w, ErrorCode(WhitelistError), WhitelistError)
		return
	}
	switch r.URL.Path {
	case "/getFile/index":
		index(w, r)
	case "/getFile/run":
		run(w, r)
	case "/getFile/download":
		download(w, r)
	default:
		http.NotFound(w, r)
	}
}

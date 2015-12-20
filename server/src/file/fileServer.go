package file

import "net/http"

func FileServer() {
	http.HandleFunc("/file/", route)
	http.ListenAndServe(":1789", nil)
}
func route(w http.ResponseWriter, r *http.Request) {
	k, ok := r.URL.Query()["key"]
	defer r.Body.Close()
	if ok && len(k) == 1 {
		if k[0] == "ABCDEF" {
			http.FileServer(http.Dir("./")).ServeHTTP(w, r)
		}
		w.Write([]byte("Sorry Have no."))
	}
}

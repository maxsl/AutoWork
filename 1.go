package main

import (
	"io"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/getFile/download", download)
	http.ListenAndServe(":1789", nil)
}

func download(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id := r.FormValue("JobId")
	file := r.FormValue("file")
	if len(id) > 0 {
		var list []string
		if len(file) <= 0 {
			l, err := ioutil.ReadDir("tmp/" + id)
			if err != nil {
				http.NotFound(w, r)
				return
			}
			for _, v := range l {
				list = append(list, v.Name())
			}
			writeIndex(w, id, list)
		} else {
			http.ServeFile(w, r, "tmp/"+id+"/"+file)
		}
	} else {
		http.NotFound(w, r)
	}
}

func writeIndex(w io.Writer, JobId string, list []string) {
	w.Write([]byte("<html><pre>"))
	for _, v := range list {
		str := "<a href=download?JobId=" + JobId + "&file=" + v + ">" + v + "</a><br>"
		w.Write([]byte(str))
	}
	w.Write([]byte("</pre></html>"))
}

package getFile

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

var exec_Chan chan Exec = make(chan Exec, 5)

func exec_run() {
	for {
		e := <-exec_Chan
		_, err := e.Copy()
		if err != nil {
			if config.Debug {
				Log.Printf("zip %s faild,%s\n", e.JobId, err)
			}
			continue
		}
		client(e.JobId)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Log.Println(err)
		http.Error(w, ErrorCode(readBodyError), readBodyError)
		return
	}
	var j Job
	err = json.Unmarshal(buf, &j)
	if err != nil {
		if config.Debug {
			Log.Println(err)
		}
		http.Error(w, ErrorCode(unmarshalError), unmarshalError)
		return
	}
	if config.Debug {
		Log.Println("Recive job:", j)
	}
	e := j.Start().GetFilesInfo()
	if e.Size <= 0 {
		http.Error(w, ErrorCode(notFoundFile), notFoundFile)
		return
	}
	if config.Debug {
		Log.Println("Search finished:", e)
	}
	buf, err = json.Marshal(e)
	if err != nil {
		Log.Println(err)
		http.Error(w, ErrorCode(marshalError), marshalError)
		return
	}
	w.Write(buf)
}

func run(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Log.Println(err)
		http.Error(w, ErrorCode(readBodyError), readBodyError)
		return
	}
	var l Exec
	err = json.Unmarshal(buf, &l)
	if err != nil {
		Log.Println(err)
		http.Error(w, ErrorCode(unmarshalError), unmarshalError)
		return
	}
	exec_Chan <- l
	if config.Debug {
		Log.Println("Run Exec_job:", l)
	}
}

func download(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	Value := r.FormValue("JobId")
	file := config.TempPath + Value + ".zip"
	_, err := os.Lstat(file)
	if err != nil {
		Log.Println("download file: ", err)
		http.Error(w, ErrorCode(notFoundFile), notFoundFile)
		return
	}
	if config.Debug {
		Log.Println("download file: ", file)
	}
	http.ServeFile(w, r, file)
}

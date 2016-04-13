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
				Log.Println(err)
			}
			return
		}
		code, err := client(e.JobId)
		if err != nil || code != 200 {
			if config.Debug {
				Log.Println(err)
			}
			return
		}
	}
}

func Server() error {
	go exec_run()
	http.HandleFunc("/", route)
	err := http.ListenAndServe(config.ListenIP+":"+config.Port, nil)
	if err != nil {
		return err
	}
	return nil
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
	if config.Debug {
		Log.Println("Receive exec job:", l)
	}
	exec_Chan <- l
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

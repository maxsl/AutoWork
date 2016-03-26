package fileagent

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
)

var FileInfo chan *FilesInfo = make(chan *FilesInfo, 5)

func backresolve() {
	for {
		info := <-FileInfo
		go func(info *FilesInfo) {
			if info.Host == "" {
				if debug {
					currentLine("callBack host can't null")
				}
				return
			}
			list := info.Run()
			buf, err := json.Marshal(list)
			if err != nil {
				if debug {
					currentLine(err)
				}
				return
			}
			b := bytes.NewReader(buf)
			for i := 0; i < 3; i++ {
				resp, err := http.Post(info.Host, BodyType, b)
				if err == nil && resp.StatusCode == 200 {
					break
				}
				if debug {
					currentLine(err)
				}
			}
		}(info)
	}
}

func Server(ip string) {
	go backresolve()
	http.HandleFunc("/download", download)
	http.HandleFunc("/ack", ack)
	http.HandleFunc("/", index)
	http.ListenAndServe(ip, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var request *Job = new(Job)
	err := resolve(w, r, request)
	if err != nil {
		return
	}
	if request.Tag == "printfiles" {
		result := request.PrintFiles()
		buf, err := json.Marshal(result)
		if err != nil {
			http.Error(w, "MarshalError", MarshalError)
			return
		}
		w.Write(buf)
		return
	}
	result := request.GetFilesInfo()
	buf, err := json.Marshal(*result)
	if err != nil {
		http.Error(w, "MarshalError", MarshalError)
		return
	}
	w.Write(buf)
	return
}

func ack(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var t *FilesInfo = new(FilesInfo)
	err := resolve(w, r, t)
	if err != nil {
		return
	}
	FileInfo <- t
}

func download(w http.ResponseWriter, r *http.Request) {

}

func resolve(w http.ResponseWriter, r *http.Request, t interface{}) error {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil || !ipIsLanIP(ip) {
		http.Error(w, "ForbiddenAccess", ForbiddenAccess)
		return err
	}
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "ReadBodyError", ReadBodyError)
		return err
	}
	err = json.Unmarshal(buf, t)
	if err != nil {
		http.Error(w, "UnmarshalError", UnmarshalError)
		return err
	}
	return nil
}

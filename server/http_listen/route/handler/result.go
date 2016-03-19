package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	Wait = iota
	Run
	Finish
)

type JobResult struct {
	JobID  string
	Action string
	User   string
	Result string
	Tag    string
	Status int
}

func receive(w http.ResponseWriter, r *http.Response) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println(err.Error())
		w.Write([]byte("unknowError"))
		return
	}
	var rst JobResult
	err = json.Unmarshal(buf, &rst)
	if err != nil {
		println(err.Error())
		w.Write([]byte("unknowError"))
		return
	}
	fmt.Println(rst)
	w.Write([]byte("OK"))
}

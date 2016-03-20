package route

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	JobWait = iota
	JobRun
	JobFinish
)

type JobResult struct {
	JobID  string `json:jobid`
	Action string `json:action`
	User   string `json:user`
	Result string `json:result`
	Tag    string `json:tag`
	Status int    `json:status`
}

func Receive(w http.ResponseWriter, r *http.Request) {
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

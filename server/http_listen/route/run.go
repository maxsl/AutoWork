package route

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/AutoWork/server"
)

func run(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, StatusText(StatusReadBodyError), StatusReadBodyError)
		return
	}
	var msg server.ClientMsg
	err = json.Unmarshal(buf, &msg)
	if err != nil {
		http.Error(w, StatusText(StatusUnmarshalError), StatusUnmarshalError)
		return
	}
	if !checkmsg(&msg) {
		http.Error(w, StatusText(StatusPostArgsError), StatusPostArgsError)
		return
	}
	job := server.CreateJob(&msg)
	var event *server.Event = &server.Event{job, msg.ServerList}
	event.Put()
}

func checkmsg(msg *server.ClientMsg) bool {
	if msg.Action == "" || msg.Body == "" || len(msg.ServerList) == 0 || msg.Action == "" {
		return false
	}
	return true
}

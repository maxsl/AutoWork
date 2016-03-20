package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/AutoWork/server"
)

type clientMsg struct {
	User       string   `json:user`
	Action     string   `json:action`
	Body       string   `json:body`
	ServerList []string `json:serverlist`
}

func Run(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	var msg server.ClientMsg
	err = json.Unmarshal(buf, &msg)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	if !checkmsg(&msg) {
		w.Write(buf)
		return
	}
	job := server.CreateJob(&msg)
	var event *server.Event = &server.Event{job, msg.ServerList}
	go event.Put()
	w.Write([]byte("ok"))
}
func checkmsg(msg *server.ClientMsg) bool {
	if msg.Action == "" || msg.Body == "" || len(msg.ServerList) == 0 || msg.Action == "" {
		return false
	}
	return true
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/AutoWork/server"
	"github.com/AutoWork/server/http_listen"
)

const bodyType = "application/x-www-form-urlencoded"

func main() {
	go http_listen.RunHttpServer(":1789")
	go server.StartServer()
	select {}
}

type clientMsg struct {
	User       string   `json:user`
	Action     string   `json:action`
	Body       string   `json:body`
	ServerList []string `json:serverlist`
}

func cmain() {
	var x clientMsg = clientMsg{User: "root", Action: "cmd",
		Body: "ipconfig", ServerList: []string{"127.0.0.1", "172.18.80.247"}}
	b, _ := json.Marshal(x)
	buf := bytes.NewReader(b)
	resp, err := http.Post("http://127.0.0.1:1789/run", bodyType, buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Request.Method)
}

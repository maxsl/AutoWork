package main

import (
	"github.com/czixchen/Log_server/client"
)

func main() {
	client.Server(":1789")
}

//package main

//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	"net/http"

//	"github.com/czixchen/Log_server/client"
//)

//func main() {
//	/*var info client.Job = client.Job{JobId: "cp_19216823", Path: "C:/test/server",
//	Name: []string{"game/logs/log", "gate/logs/log", "mail/logs/log", "room-manage/logs/log"},
//	Tag:  "compress"}*/
//	var info client.FilesInfo = client.FilesInfo{Path: "C:/test/server/",
//		Files:    []string{"game/logs/log", "gate/logs/log", "mail/logs/log", "room-manage/logs/log"},
//		AbsFiles: []string{"C:/oem8.log"},
//		Size:     737,
//		JobId:    "cp_19216823",
//		Tag:      "copyfiles", Host: "http://172.18.80.247:2789"}
//	b, _ := json.Marshal(info)
//	buf := bytes.NewReader(b)
//	resp, err := http.Post("http://172.18.80.247:1789/ack", client.BodyType, buf)
//	if err != nil {
//		fmt.Println(err)
//	}
//	defer resp.Body.Close()
//	b, _ = ioutil.ReadAll(resp.Body)
//	fmt.Println(string(b))
//}

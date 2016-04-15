package center

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
)

func sendResult(id, ip string) {
	m, ok := Contacts.del(id)
	if !ok {
		return
	}
	p := make(map[string][]string)
	p["action"] = []string{m.mode}
	p["tos"] = []string{m.contacts}
	p["content"] = []string{"http://" + ip + "/getFile/download?JobId=" + id}
	if config.Debug {
		Log.Println("send result to ", config.MsgApi, p)
	}
	http.PostForm(config.MsgApi, p)
}

func forwardrequst(iplist []string, buf []byte) (m map[string]Exec) {
	m = make(map[string]Exec)
	var w *sync.WaitGroup = new(sync.WaitGroup)
	for _, ip := range iplist {
		w.Add(1)
		go func(w *sync.WaitGroup, ip string) {
			defer w.Done()
			r := bytes.NewReader(buf)
			resp, err := http.Post("http://"+ip+"/getFile/index", bodyType, r)
			if err != nil {
				Log.Println("send request faild:", err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				Log.Printf("%v error code:%v\n", ip, resp.StatusCode)
				return
			}
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				Log.Println("Read result faild:", err)
				return
			}
			var e Exec
			err = json.Unmarshal(b, &e)
			if err != nil {
				Log.Println("Unmarshal result faild:", err)
				return
			}
			Log.Printf("recive %v msg:%v\n", ip, e)
			m[ip] = e
		}(w, ip)
	}
	w.Wait()
	return
}

func request(m map[string]Exec) {
	if config.Debug {
		Log.Println("run msg:", m)
	}
	var w *sync.WaitGroup = new(sync.WaitGroup)
	for ip, msg := range m {
		w.Add(1)
		go func(w *sync.WaitGroup, ip string, msg Exec) {
			defer w.Done()
			buf, err := json.Marshal(msg)
			if err != nil {
				Log.Println("marshal msg faild:", err)
				return
			}
			r := bytes.NewReader(buf)
			u := "http://" + ip + "/getFile/run"
			resp, err := http.Post(u, bodyType, r)
			if err != nil {
				Log.Println("send request faild:", err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				Log.Printf("%v error code:%v\n", ip, resp.StatusCode)
				return
			} else {
				if config.Debug {
					Log.Println("Request successful ->", u)
				}
			}
		}(w, ip, msg)
	}
	w.Wait()
}

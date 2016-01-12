package main

import (
	"alltype/server"
	"encoding/json"
	"fmt"
	"listen"
	"net"
	"public_tool/rwMsg"
	"sync"
	"time"
)

func smain() {
	listen.Listen()
}
func main() {
	con, err := net.Dial("tcp", cfg.IP)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	buf := rwMsg.GetrwcPool(con)
	defer rwMsg.PutrwcPool(buf)
	b, _, err := buf.ReadLine()
	fmt.Println(string(b))
	if err != nil || string(b) == "AlreadyExist" {
		fmt.Println("连接失败.")
		return
	}
	fmt.Fprint(buf, "Auth")
	go testPing(buf)
	go handerClientMsg(buf)
	select {}
}

type config struct {
	IP      string
	Inteval int
}

var flushtime int64 = 0
var cfg config = config{"127.0.0.1:1789", 10}
var lock *sync.Mutex = new(sync.Mutex)

func handerClientMsg(buf *rwMsg.Rwc) {
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			fmt.Printf("读取数据失败,错误信息:%s\n", err)
			break
		}
		msg, err := rwMsg.Decode(line)
		if err != nil {
			fmt.Printf("解析数据错误,错误信息:%s\n", err)
			continue
		}
		handmsg(msg)
		rwMsg.PutJsonPool(msg)
	}
}
func handmsg(m *server.Msg) {
	fmt.Println(*m)
}
func testPing(buf *rwMsg.Rwc) {
	for {
		if time.Now().Unix()-flushtime > int64(cfg.Inteval) {
			ping(buf)
			flushtime = time.Now().Unix()
			time.Sleep(time.Duration(cfg.Inteval) * time.Second)
		}
	}
}

func ping(buf *rwMsg.Rwc) {
	lock.Lock()
	defer lock.Unlock()
	s := server.Msg{Action: "Test.Ping"}
	b, _ := json.Marshal(s)
	_, err := buf.Write(b)
	if err != nil {
		reconnection(buf)
		fmt.Println("重连成功.")
	}
}
func reconnection(buf *rwMsg.Rwc) {
	var err error
	for {
		buf.Conn, err = net.Dial("tcp", cfg.IP)
		if err != nil {
			fmt.Println("重连失败.")
			time.Sleep(1e9)
			continue
		}
		break
	}
	line, _, err := buf.ReadLine()
	if err != nil {
		reconnection(buf)
	}
	fmt.Println(string(line))
	fmt.Fprint(buf, "Auth")
	buf.Buf.Reset(buf.Conn)
}

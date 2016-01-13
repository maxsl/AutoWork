package main

import (
	"alltype/server"
	"encoding/json"
	"fmt"
	"listen"
	"log"
	"net"
	"public_tool/rwMsg"
	"sync"
	"time"
)

func main() {
	listen.Listen()
}
func cmain() {
	con, err := net.Dial("tcp", cfg.IP)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	buf := rwMsg.GetrwcPool(con)
	defer rwMsg.PutrwcPool(buf)
	b, _, err := buf.ReadLine()
	log.Println(string(b))
	if err != nil || string(b) == "AlreadyExist" {
		log.Println("连接失败.")
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
			log.Printf("读取数据失败,错误信息:%s\n", err)
			break
		}
		msg, err := rwMsg.Decode(line)
		if err != nil {
			log.Printf("解析数据错误,错误信息:%s\n", err)
			continue
		}
		handmsg(msg)
		rwMsg.PutJsonPool(msg)
	}
}
func handmsg(m *server.Msg) {
	log.Println(*m)
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
	}
}
func reconnection(buf *rwMsg.Rwc) {
	var err error
	for {
		buf.Conn, err = net.Dial("tcp", cfg.IP)
		if err != nil {
			log.Println("重连失败.")
			time.Sleep(1e9)
			continue
		}
		break
	}
	buf.Buf.Reset(buf.Conn)
	line, _, err := buf.ReadLine()
	if err != nil {
		reconnection(buf)
		return
	}
	log.Println(string(line))
	fmt.Fprint(buf, "Auth")
}

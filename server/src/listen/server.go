package listen

import (
	"alltype/server"
	"fmt"
	"net"
	"public_tool/rwMsg"
	"strings"
	"sync"
	"time"
)

func Listen() {
	lis, err := net.Listen("tcp", "0.0.0.0:1789")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer lis.Close()
	for {
		con, err := lis.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go Hander(con)
	}
}

var exitchan chan bool = make(chan bool, 1)

func Hander(con net.Conn) {
	defer con.Close()
	ip := strings.Split(con.RemoteAddr().String(), ":")[0]
	defer statusMap.Del(ip)
	if !statusMap.Add(ip) {
		fmt.Fprint(con, "AlreadyExist\r\n")
		return
	}
	buf := rwMsg.GetrwcPool(con)
	defer rwMsg.PutrwcPool(buf)
	line, _, err := buf.ReadLine()
	//验证源主机是否授权,规则自己写.
	if err != nil || string(line) != "Auth" {
		return
	}
	go read(buf, ip)
	go write(buf, ip)
	<-exitchan
}
func write(buf *rwMsg.Rwc, ip string) {
	for {
		msg := <-statusMap[ip]
		b, err := handlerMsg(msg, ip)
		if err != nil {
			fmt.Printf("主机%s编码消息%s失败,错误信息:%s\n", ip, *msg, err)
			continue
		}
		_, err = buf.Write(b)
		if err != nil {
			fmt.Printf("往%s发送消息%s是被,错误信息:%s\n", ip, *msg, err)
			break
		}
	}
	exitchan <- true
}
func read(buf *rwMsg.Rwc, ip string) {
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			fmt.Printf("读取主机%s数据失败,错误信息:%s\n", ip, err)
			break
		}
		msg, err := rwMsg.Decode(line)
		if err != nil {
			fmt.Printf("解析主机%s返回数据错误,错误信息:%s\n", ip, err)
			rwMsg.PutJsonPool(msg)
			continue
		}
		fmt.Printf("接收到%s返回结果:\n", ip)
		fmt.Println(*msg)
		rwMsg.PutJsonPool(msg)
	}
	exitchan <- true
}

func handlerMsg(msg *server.Msg, ip string) ([]byte, error) {
	handmsg(msg, ip)
	b, err := rwMsg.Encode(msg)
	if err != nil {
		return empty, err
	}
	return b, nil
}

func handmsg(m *server.Msg, ip string) {
	m.Address = cfg.DownLoadAddress
	m.Remark = cfg.ResultAddress
}

type config struct {
	DownLoadAddress string
	ResultAddress   string
}

var cfg config = config{"http://127.0.0.1:8888/download", "http://127.0.0.1:8888/result"}

type statusM map[string]chan *server.Msg

var empty []byte
var statusMap statusM = make(statusM)

func (statusM) Del(ip string) {
	putclientStatusPool(statusMap[ip])
	delete(statusMap, ip)
}

func (statusM) Add(ip string) bool {
	_, ok := statusMap[ip]
	if ok {
		return false
	}
	statusMap[ip] = getclientStatusPool()
	return true
}

var GlobalLock sync.Mutex
var id int64 = getNow()

func NewID() int64 {
	GlobalLock.Lock()
	defer GlobalLock.Unlock()
	id = id + 1
	return id
}

func getNow() int64 {
	return time.Now().Unix()
}

var clientStatusPool sync.Pool

func getclientStatusPool() chan *server.Msg {
	c := clientStatusPool.Get()
	if c != nil {
		return c.(chan *server.Msg)
	}
	channal := make(chan *server.Msg, 1)
	return channal
}

func putclientStatusPool(c chan *server.Msg) {
	c = nil
	clientStatusPool.Put(c)
}

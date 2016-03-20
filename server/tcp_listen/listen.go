package tcp_listen

import (
	"encoding/base64"
	"strings"
	"sync"
)

type ClientMap struct {
	lock   *sync.RWMutex
	Client map[string]NewConnection
	auth   func(b []byte) bool
}

var Clients *ClientMap = &ClientMap{new(sync.RWMutex), make(map[string]NewConnection), auth}
var encode *base64.Encoding = base64.RawStdEncoding

func (self *ClientMap) RegisterAuthFunc(f func(b []byte) bool) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.auth = f
}

func (self *ClientMap) IsExist(ip string) bool {
	self.lock.RLock()
	defer self.lock.RUnlock()
	_, ok := self.Client[ip]
	if ok {
		return true
	}
	return false
}

func (self *ClientMap) Get(ip string) NewConnection {
	self.lock.RLock()
	defer self.lock.RUnlock()
	return self.Client[ip]
}

func (self *ClientMap) Put(con NewConnection) bool {
	ip := strings.Split(con.RemoteAddr().String(), ":")[0]
	if self.IsExist(ip) {
		self.Close(con)
		return false
	}
	authMsg, err := con.DelimRead()
	if err != nil {
		self.Close(con)
		println(err.Error())
		return false
	}
	msg, err := encode.DecodeString(string(authMsg))
	if err != nil {
		self.Close(con)
		println(err.Error())
		return false
	}
	if !self.auth(msg) {
		self.Close(con)
		return false
	}
	self.lock.Lock()
	defer self.lock.Unlock()
	self.Client[ip] = con
	return true
}

func (self *ClientMap) GetClients() map[string]NewConnection {
	self.lock.Lock()
	defer self.lock.Unlock()

	return self.Client
}

func (self *ClientMap) Close(con NewConnection) error {
	ip := strings.Split(con.RemoteAddr().String(), ":")[0]
	self.lock.Lock()
	defer self.lock.Unlock()
	delete(self.Client, ip)
	return con.Close()
}

func RunServer(ip string) error {
	lis, err := ListenServer(ip)
	if err != nil {
		return err
	}
	defer lis.Close()
	for {
		con, err := lis.Accept()
		if err != nil {
			continue
		}
		go Clients.Put(con)
	}
}

func auth(authmsg []byte) bool {
	if string(authmsg) != "Hello World" {
		return false
	}
	return true
}

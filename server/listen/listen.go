package listen

import (
	"encoding/base64"
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

func (self *ClientMap) Put(con NewConnection) bool {
	ip := con.RemoteAddr().String()
	if self.IsExist(ip) {
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
	self.lock.Lock()
	defer self.lock.Unlock()
	delete(self.Client, con.RemoteAddr().String())
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

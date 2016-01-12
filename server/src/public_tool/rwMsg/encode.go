package rwMsg

import (
	"alltype/server"
	"bufio"
	"encoding/json"
	"net"
	"sync"
)

func Decode(b []byte) (*server.Msg, error) {
	j := GetJsonPool()
	err := json.Unmarshal(b, j)
	if err != nil {
		PutJsonPool(j)
		return nil, err
	}
	return j, nil
}

func Encode(msg *server.Msg) ([]byte, error) {
	return json.Marshal(msg)
}

var JsonPool sync.Pool

func GetJsonPool() *server.Msg {
	buf := JsonPool.Get()
	if buf != nil {
		b := buf.(*server.Msg)
		return b
	}
	return new(server.Msg)
}

func PutJsonPool(m *server.Msg) {
	*m = server.Msg{}
	JsonPool.Put(m)
}

type Rwc struct {
	Conn net.Conn
	Buf  *bufio.Reader
}

func (self *Rwc) ReadLine() ([]byte, bool, error) {
	return self.Buf.ReadLine()
}
func (self *Rwc) Write(b []byte) (int, error) {
	b = append(b, []byte("\r\n")...)
	return self.Conn.Write(b)
}

var rwcPool sync.Pool

func GetrwcPool(con net.Conn) *Rwc {
	buf := rwcPool.Get()
	if buf != nil {
		b := buf.(*Rwc)
		b.Buf.Reset(con)
		b.Conn = con
		return b
	}
	rw := new(Rwc)
	rw.Conn = con
	rw.Buf = bufio.NewReader(con)
	return rw
}

func PutrwcPool(b *Rwc) {
	b.Buf.Reset(nil)
	b.Conn = nil
	rwcPool.Put(b)
}

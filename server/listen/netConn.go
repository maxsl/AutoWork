package listen

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	headlen = 5
	delim   = '\n'
)

var (
	SendHeadError   = fmt.Errorf("SendMsgError")
	ReceiveMsgError = fmt.Errorf("MegLengthError")
	empty           = []byte{}
)

type NewConnection interface {
	net.Conn
	DelimRead() ([]byte, error)
	DelimWrite(b []byte) (n int, err error)
}

func ListenServer(ip string) (*newListener, error) {
	lis, err := net.Listen("tcp", ip)
	if err != nil {
		return nil, err
	}
	return &newListener{lis}, nil
}

type newListener struct {
	lis net.Listener
}

func (self *newListener) Accept() (NewConnection, error) {
	conn, err := self.lis.Accept()
	if err != nil {
		return nil, err
	}
	return NewConn(conn), nil
}

func (self *newListener) Close() error {
	return self.lis.Close()
}

func (self *newListener) Addr() net.Addr {
	return self.lis.Addr()
}

var connctionPool sync.Pool

func getConnection(conn net.Conn) *Connection {
	c := connctionPool.Get()
	if c != nil {
		con, ok := c.(*Connection)
		if ok {
			con.con = conn
			con.buf = bufio.NewReader(conn)
			return con
		}
	}
	return &Connection{conn, bufio.NewReader(conn)}
}

func putConnection(conn *Connection) {
	conn.buf.Reset(nil)
	connctionPool.Put(conn)
}

func NewConn(conn net.Conn) *Connection {
	return getConnection(conn)
}

type Connection struct {
	con net.Conn
	buf *bufio.Reader
}

func (self *Connection) DelimRead() ([]byte, error) {
	head := make([]byte, headlen)
	n, err := self.buf.Read(head)
	if err != nil && n != 5 {
		if n == 0 {
			return empty, err
		}
		return empty, ReceiveMsgError
	}
	bodylen, _ := binary.Uvarint(head)
	if bodylen <= 0 {
		return empty, err
	}
	line, err := self.buf.ReadSlice(delim)
	if err != nil {
		return line, err
	}
	if len(line) != int(bodylen+1) {
		return line, ReceiveMsgError
	}
	return line[:int(bodylen)], nil
}

func (self *Connection) DelimWrite(b []byte) (n int, err error) {
	msglen := uint64(len(b))
	head := make([]byte, headlen)
	binary.PutUvarint(head, msglen)
	n, err = self.con.Write(head)
	if n != headlen {
		if err != nil {
			return 0, err
		}
		return 0, SendHeadError
	}
	n, err = self.con.Write(b)
	self.con.Write([]byte{delim})
	if err != nil {
		return n, err
	}
	return n, nil
}

func (self *Connection) Read(b []byte) (int, error) {
	return self.con.Read(b)
}

func (self *Connection) Write(b []byte) (int, error) {
	return self.con.Write(b)
}

func (self *Connection) Close() error {
	err := self.con.Close()
	putConnection(self)
	return err
}

func (self *Connection) LocalAddr() net.Addr {
	return self.con.LocalAddr()
}

func (self *Connection) RemoteAddr() net.Addr {
	return self.con.RemoteAddr()
}

func (self *Connection) SetDeadline(t time.Time) error {
	return self.con.SetDeadline(t)
}

func (self *Connection) SetReadDeadline(t time.Time) error {
	return self.con.SetReadDeadline(t)
}

func (self *Connection) SetWriteDeadline(t time.Time) error {
	return self.con.SetWriteDeadline(t)
}

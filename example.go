package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

var empty []byte

func main() {
	lis, _ := net.Listen("tcp", ":1789")
	for {
		con, err := lis.Accept()
		if err != nil {
			continue
		}
		go registerConn(con)
	}
}

type connectionInfo struct {
	conn net.Conn
	buf  *bufio.Reader
}

func (self *connectionInfo) onceRead() ([]byte, error) {
	return reader(self.buf)
}

func (self *connectionInfo) StillRead() {
	defer self.conn.Close()
	for {
		msg, err := reader(self.buf)
		if err != nil {
			if err == io.EOF {
				continue
			}
			return
		}
		fmt.Printf("收到的信息: %s\n", string(msg))
	}
}

func registerConn(conn net.Conn) {
	reg := &connectionInfo{conn, bufio.NewReader(conn)}
	msg, err := reg.onceRead()
	if err != nil {
		reg.conn.Close()
		return
	}
	fmt.Printf("注册信息: %s\n", string(msg))
	reg.StillRead()
}

func reader(buf *bufio.Reader) ([]byte, error) {
	head := make([]byte, 5)
	_, err := buf.Read(head)
	if err != nil {
		return empty, err
	}
	bodyLen, _ := binary.Uvarint(head)
	line, err := buf.ReadSlice('\n')
	if err != nil {
		return empty, err
	}
	if uint64(len(line)-1) != bodyLen {
		return empty, io.EOF
	}
	return line[:bodyLen], nil
}

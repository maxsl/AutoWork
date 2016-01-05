package listen

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	IntToByte "public_tool/intTobyte"
	ftime "public_tool/time"
	"sync"
)

var empty []byte = []byte{}

func ListenServer(IP string) {
	lis, err := net.Listen("tcp", IP)
	if err != nil {
		fmt.Println(ftime.LineTime(), " ", err)
		os.Exit(1)
	}
	for {
		con, err := lis.Accept()
		if err != nil {
			continue
		}
		go HandleMessage(con)
	}
}

func HandleMessage(con net.Conn) {
	buf := GetReadPool(con)
	defer func(con net.Conn, b *bufio.Reader) {
		PutReadPool(r)
		con.Close()
	}(con, buf)
}

var bufRead sync.Pool

func ReadMsg(b *bufio.Reader) ([]byte, error) {
	by, _, err := b.ReadLine()
	if err != nil {
		return empty, err
	}
	return by, nil
}

func GetReadPool(r io.Reader) *bufio.Reader {
	if v := bufRead.Get(); v != nil {
		br := v.(*bufio.Reader)
		br.Reset(r)
		return br
	}
	return bufio.NewReader(r)
}

func PutReadPool(br *bufio.Reader) {
	br.Reset(nil)
	bufioReaderPool.Put(br)
}

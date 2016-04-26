package tcp_listen

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func Test_msg(t *testing.T) {
	go func() {
		lis, err := Server(":2789")
		if err != nil {
			t.Logf("%v\n", err)
			return
		}
		defer lis.Close()
		for {
			conn, err := lis.Accept()
			if err != nil {
				t.Logf("%v\n", err)
				continue
			}
			Handler(NewConnction(conn))
		}
	}()

	conn, err := Client("127.0.0.1:2789")
	if err != nil {
		t.Logf("%v\n", err)
		return
	}
	h := NewConnction(conn)
	defer h.Close()
	buf := []byte("Hello World")

	head := NewMsgBytes(h.MsgTag, int64(len(buf)))
	if head == nil {
		return
	}
	copy(head[HeadLenght:], buf)
	_, err = h.Write(head)
	if err != nil {
		fmt.Println(err)
		return
	}

	buf, err = ioutil.ReadAll(h)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(buf))
}

func Handler(conn *HeadConnection) {
	defer conn.Close()
	buf, err := ioutil.ReadAll(conn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(buf))
	head := NewMsgBytes(conn.MsgTag, int64(len(buf)))
	if head == nil {
		fmt.Println(headLenghtError)
	}
	copy(head[HeadLenght:], buf)
	conn.Write(head)
}

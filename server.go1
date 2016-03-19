package main

import (
	"fmt"
	"io"

	"github.com/czxichen/AutoWork/server/listen"
)

func main() {
	lis, err := listen.ListenServer(":1789")
	if err != nil {
		fmt.Println(err)
	}
	defer lis.Close()
	for {
		con, err := lis.Accept()
		if err != nil {
			continue
		}
		go handler(con)
	}
}

func handler(con listen.NewConnection) {
	defer con.Close()
	msg := make(chan string, 5)
	go func(con listen.NewConnection, msg chan string) {
		for {
			buf, err := con.DelimRead()
			if err != nil {
				if err == listen.ReceiveMsgError || err == io.EOF {
					continue
				}
				break
			}
			fmt.Println("收到的消息: ", string(buf))
			msg <- string(buf)
		}
		msg <- "exit"
	}(con, msg)
	for {
		buf := <-msg
	ResendMsg:
		_, err := con.DelimWrite([]byte(buf))
		if err != nil {
			if err == listen.SendHeadError {
				goto ResendMsg
			}
			break
		}
	}
}

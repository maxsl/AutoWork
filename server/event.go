package server

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/czxichen/AutoWork/server/tcp_listen"
)

type Event struct {
	Action *Job
	Remote []string
	Faild  []string
}

const eventChanLen = 50

var receive chan *Event = make(chan *Event, eventChanLen)

func (self *Event) Put() {
	receive <- self
}

func (self *Event) String() string {
	return fmt.Sprintf("{{%s %s %s %s} %s %s", self.Action.JobID, self.Action.Action,
		self.Action.User, self.Action.Body, self.Remote, self.Remote)
}

func GetEventChanLen() int {
	return len(receive)
}

func StartServer() error {
	go func() {
		err := tcp_listen.RunServer(":2789")
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}()

	for {
		event := <-receive
		fmt.Println("收到消息:", event.Action.String())
		go func(e *Event) {
			defer e.Action.Close()
			unsend := sendMgs(e)
			if len(unsend) > 0 {
				//	第二次发送没有发送成功的消息.
				unsend = sendMgs(&Event{Action: e.Action, Remote: unsend})
			}
			e.Faild = unsend
			if len(e.Faild) != 0 {
				fmt.Println("未发送: ", e.Faild)
			}
			//时间记录到数据库.
		}(event)
	}
}

func sendMgs(e *Event) []string {
	msg := []byte(e.Action.Base64EncodeString())
	var unsend []string
	var wait *sync.WaitGroup = new(sync.WaitGroup)
	for _, v := range e.Remote {
		con := tcp_listen.Clients.Get(v)
		if con == nil {
			unsend = append(unsend, v)
			continue
		}
		wait.Add(1)
		go func(con tcp_listen.NewConnection, wait *sync.WaitGroup) {
			defer wait.Done()
			ip := strings.Split(con.RemoteAddr().String(), ":")[0]
			_, err := con.DelimWrite(msg)
			if err != nil {
				unsend = append(unsend, ip)
				tcp_listen.Clients.Close(con)
				return
			}
			buf, err := con.DelimRead()
			if err != nil || string(buf) != "ok" {
				unsend = append(unsend, ip)
				tcp_listen.Clients.Close(con)
				return
			}
		}(con, wait)
	}
	wait.Wait()
	return unsend
}

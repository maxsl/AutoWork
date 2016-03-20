package server

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/AutoWork/server/tcp_listen"
)

type Event struct {
	Action *Job
	Remote []string
}

func (self *Event) Put() {
	receive <- self
}

const eventChanLen = 50

var receive chan *Event = make(chan *Event, eventChanLen)

func GetEventChanLen() int {
	return len(receive)
}
func GetEvent() *Event {
	return <-receive
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
		event := GetEvent()
		//		fmt.Println("收到消息:", event.Action.String())
		//		fmt.Println(tcp_listen.Clients.Client)
		go func(e *Event) {
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
			//			fmt.Println("未发送: ", unsend)
		}(event)
	}
}

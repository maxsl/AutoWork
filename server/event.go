package server

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

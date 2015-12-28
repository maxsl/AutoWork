package alltype

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

type AgentInfo struct {
	IP  string
	MAC string
}

func (self AgentInfo) Md5() string {
	m := md5.New()
	str := self.IP + "1" + self.MAC
	if len(str) <= 0 {
		return ""
	}
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum([]byte{}))
}

type postMessage struct {
	Body []byte
}

func (self *postMessage) Read(buf []byte) (int, error) {
	if len(self.Body) <= 0 {
		return 0, io.EOF
	}
	n := copy(buf, self.Body)
	self.Body = self.Body[n:]
	return n, nil
}

type StartConfig struct {
	Port string
	IP   string
}

var Config StartConfig

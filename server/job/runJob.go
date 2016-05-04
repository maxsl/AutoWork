package job

import (
	"bytes"
	"errors"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

const MaxBufferSize = 1 << 24 //默认的缓存大小16M

var (
	timeOutError    = errors.New("Run is timeout")
	manualExitError = errors.New("Manual Exit")
)

func GetCommond(user string, cmd string) []string {
	switch runtime.GOOS {
	case "windows":
		return []string{"cmd", "/C", cmd}
	case "linux":
		return []string{"su", "-", user, "-c", cmd}
	}
	return []string{}
}

type jobInfo struct {
	commond  []string
	timeout  int
	ExitChan chan bool
	Result   *bytes.Buffer // make([]byte,0,MaxBufferSize)
	ErrInfo  error
}

func (self *jobInfo) Run() {
	cmd := exec.Command(self.commond[0], self.commond[1:]...)

	go func() {
		cmd.Stdout = self.Result
		self.ErrInfo = cmd.Run()
		self.ExitChan <- true
	}()

	time.AfterFunc(time.Duration(self.timeout)*time.Second, func() {
		self.ErrInfo = timeOutError
		cmd.Process.Kill()
		self.ExitChan <- false
	})

	<-self.ExitChan
}

func (self *jobInfo) Stop() {
	self.ErrInfo = manualExitError
	self.ExitChan <- false
}

var jobInfoPool sync.Pool

func getjobInfo() *jobInfo {
	j := jobInfoPool.Get()
	if j != nil {
		J, ok := j.(*jobInfo)
		if ok {
			return J
		}
	}
	return &jobInfo{Result: bytes.NewBuffer(make([]byte, 0, MaxBufferSize))}
}

func putjobInfo(j *jobInfo) {
	j.Result.Reset()
	jobInfoPool.Put(j)
}

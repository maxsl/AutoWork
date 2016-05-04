package job

import (
	"os"
	"sync"
)

type job struct {
	Id         string
	Date       int64
	User       string
	RunEnv     map[string]string
	Exec       string
	ResultPath string
	StartTime  int64
	TimeOut    int64
	Status     int
	ExitChan   chan bool
}

func (self *job) Run() {
	if len(self.RunEnv) > 0 {
		for key, value := range self.RunEnv {
			os.Setenv(key, value)
		}
	}

}

func (self *job) GetStartTime() int64 {
	return self.StartTime
}

var jobPool sync.Pool

func GetJob() *job {
	j := jobPool.Get()
	if j != nil {
		J, ok := j.(*job)
		if ok {
			return J
		}
	}
	return &job{ExitChan: make(chan bool, 1)}
}

func PutJob(j *job) {
	jobPool.Put(j)
}

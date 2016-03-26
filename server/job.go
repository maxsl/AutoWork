package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

var jobIDlock sync.RWMutex
var counter int64 = 0

func GetJobID(action string) string {
	jobIDlock.Lock()
	defer jobIDlock.Unlock()
	now := time.Now().Unix()
	if counter >= now {
		now = counter + 1
	}
	id := fmt.Sprintf("%s_%d", action, now)
	counter = now
	return id
}

var base64Encode *base64.Encoding = base64.RawStdEncoding

type Job struct {
	lock   *sync.RWMutex
	JobID  string `json:jobid`
	Action string `json:action`
	User   string `json:user`
	Body   string `json:body`
}

func (self *Job) Id() string {
	self.lock.RLocker()
	defer self.lock.RUnlock()
	return self.JobID
}

func (self *Job) Base64EncodeString() string {
	self.lock.Lock()
	defer self.lock.Unlock()

	buf, err := json.Marshal(*self)
	if err != nil {
		return ""
	}
	return base64Encode.EncodeToString(buf)
}

func (self *Job) String() string {
	self.lock.Lock()
	defer self.lock.Unlock()

	buf, err := json.Marshal(*self)
	if err != nil {
		return ""
	}
	return string(buf)
}

func (self *Job) Close() {
	putNewJob(self)
}

var jobPool sync.Pool

func getNewJob(action, user, body string) *Job {
	job := jobPool.Get()
	if job != nil {
		j, ok := job.(*Job)
		if ok {
			j.lock = new(sync.RWMutex)
			j.JobID = GetJobID(action)
			j.Action = action
			j.User = user
			j.Body = body
			return j
		}
	}
	return &Job{lock: new(sync.RWMutex), JobID: GetJobID(action), Action: action, User: user, Body: body}
}

func putNewJob(job *Job) {
	job.lock = nil
	jobPool.Put(job)
}

type ClientMsg struct {
	User       string   `json:user`
	Action     string   `json:action`
	Body       string   `json:body`
	Tag        string   `json:"tag,omitempty"`
	ServerList []string `json:"serverlist,omitempty`
}

func CreateJob(msg *ClientMsg) *Job {
	return getNewJob(msg.Action, msg.User, msg.Body)
}
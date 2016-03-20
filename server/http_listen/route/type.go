package route

type log interface {
	PrintfI(formate string, v ...interface{})
	PrintfW(formate string, v ...interface{})
	PrintfE(formate string, v ...interface{})
	PrintfF(formate string, v ...interface{})
}

const (
	JobWait = iota + 1
	JobRunning
	JobFinish
)

type JobResult struct {
	JobID  string `json:jobid`
	Action string `json:action`
	User   string `json:user`
	Result string `json:result`
	Tag    string `json:tag`
	Status int    `json:status`
}

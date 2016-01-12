package server

type StartConfig struct {
	Port string
	IP   string
}

var Config StartConfig

type Msg struct {
	Action  string `json:action`
	Address string `json:address`
	JobID   int64  `json:jobid`
	Remark  string `json:remark`
}

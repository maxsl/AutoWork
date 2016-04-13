package center

import (
	"flag"
	"strconv"
	"sync"
	"time"
)

type cfg struct {
	ListenIP     string `json:listenip`
	Relationship string `json:relationship`
	Debug        bool   `json:debug`
	Log          string `json:log`
	WhiteList    string `json:whitelist`
	TmpDir       string `json:tmpdir`
}

type Exec struct {
	JobId    string   `json:jobid`
	Size     int64    `json:size`
	Files    []string `json:files`
	AbsFiles []string `json:absfiles`
}

type Job struct {
	JobId  string   `json:jobid`
	Tag    int      `json:tag`
	Files  []string `json:files`
	Date   int64    `json:date`
	Regexp string   `json:regexp`
}

type SourceJob struct {
	Job    Job
	IPlist []string `json:iplist`
}

const (
	bodyType = "application/json"
)

var (
	ServerRelationship map[string]string = make(map[string]string)
	config             cfg               = cfg{ListenIP: ":2789", Relationship: "relationship.json",
		Debug: true, Log: "center.log", TmpDir: "tmp"}
	configFile string
)

func init() {
	flag.StringVar(&configFile, "-c", "center_cfg.json", "-c config-path")
	flag.Parse()
}

var unixTime int64 = time.Now().Unix()
var lock *sync.RWMutex = new(sync.RWMutex)

func getJobId() string {
	now := time.Now().Unix()
	lock.Lock()
	defer lock.Unlock()
	if now > unixTime {
		unixTime = now
		return strconv.FormatInt(now, 10)
	}
	unixTime += 1
	return strconv.FormatInt(unixTime, 10)
}

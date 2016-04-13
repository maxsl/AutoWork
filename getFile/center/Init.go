package center

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var Log *log.Logger

func Init() {
	buf, err := ioutil.ReadFile(configFile)
	if err == nil {
		err = json.Unmarshal(buf, &config)
	}
	if err != nil {
		println("init config faild", err.Error(), "use default config.")
	}
	config.TmpDir = strings.Replace(config.TmpDir, `\`, "/", -1)
	if !strings.HasSuffix(config.TmpDir, "/") {
		config.TmpDir += "/"
	}
	File, err := os.OpenFile(config.Log, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		Log = log.New(os.Stdout, "", log.LstdFlags)
		Log.Printf("init log path faild.error info:%v\nuse os.Stdout.\n", err)
	} else {
		Log = log.New(File, "", log.LstdFlags)
	}
	if config.Debug {
		Log.Println("config:", config)
	}
	buf, err = ioutil.ReadFile(config.Relationship)
	if err == nil {
		err = json.Unmarshal(buf, &ServerRelationship)
	}
	if err != nil {
		Log.Printf("init relationship faild,error info:%v\n", err)
	}
}

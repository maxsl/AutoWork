package fileagent

import (
	"log"
	"os"
)

var Log *log.Logger

func init() {
	File, err := os.OpenFile("run.log", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		println("create run.log file error")
		os.Exit(-1)
	}
	Log = log.New(File, "", log.LstdFlags)
}

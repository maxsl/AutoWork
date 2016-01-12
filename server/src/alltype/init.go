package alltype

import (
	"alltype/agent"
	"encoding/gob"
	"log"
	"os"
)

var File *os.File
var AgentMap map[string]agent.AgentInfo = make(map[string]agent.AgentInfo)

func init() {
	var err error
	File, err = os.OpenFile("AgentKey", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	if info, _ := File.Stat(); info.Size() <= 0 {
		return
	}
	err = gob.NewDecoder(File).Decode(&AgentMap)
	if err != nil {
		log.Fatalln(err)
	}
}

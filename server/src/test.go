package main

import (
	"alltype"
	"fmt"
	ftime "public_tool/time"
)

func main() {
	defer alltype.File.Close()
	fmt.Println(alltype.AgentMap)
	fmt.Println(ftime.LineTime())
}

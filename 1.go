package main

import (
	"encoding/base64"
	"fmt"

	"github.com/czxichen/AutoWork/server"
)

func main() {
	job := server.CreateJob("cmd", "root", "uname -a ")
	str := job.Base64EncodeString()
	fmt.Println(str)
	buf, err := base64.RawStdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(buf))
}

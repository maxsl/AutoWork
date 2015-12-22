package main

import (
	"fmt"
	"public_tool"
)

func main() {
	r, err := public_tool.NewRsaEncrypt()
	if err != nil {
		fmt.Println(err)
		return
	}
	buf, _ := r.Encrypt([]byte("1"))
	buf, _ = r.Decrypt(buf)
	fmt.Println(string(buf))
}

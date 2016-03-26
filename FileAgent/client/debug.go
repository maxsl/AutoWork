package client

import (
	"fmt"
	"io"
	"os"
	"runtime"
)

var (
	debug  bool      = true
	stdout io.Writer = os.Stdout
)

func currentLine(info interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		fmt.Fprintf(stdout, "File:%s Line:%d \nInfo: %v\n", file, line, info)
	}
}

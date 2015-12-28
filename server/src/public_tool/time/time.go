package public_tool

import (
	"fmt"
	"runtime"
	"time"
)

func FormatTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func LineTime() string {
	_, name, n, ok := runtime.Caller(1)
	if ok {
		return fmt.Sprintf("FileName:%s Line:%d", name, n)
	}
	return ""
}

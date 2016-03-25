package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	callback([]string{"string"})
}

func callback(result interface{}) {
	switch result.(type) {
	case string:
		fmt.Println(result)
	case []string:
		list, _ := result.([]string)
		for _, v := range list {
			fmt.Println(v)
		}
	}
}

func CurrentLine(info interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		fmt.Fprintf(os.Stdout, "File:%s Line:%d \nInfo: %v\n", file, line, info)
	}
}

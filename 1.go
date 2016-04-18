package main

import (
	"os"

	//	"github.com/czxichen/AutoWork/getFile"
	"github.com/czxichen/AutoWork/getFile/center"
)

func main() {
	if center.Daemon {
		if os.Getppid() != 1 {
			os.StartProcess(os.Args[0], os.Args, &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
			return
		} else {
			center.Server()
		}
	} else {
		center.Server()
	}
}

//func main() {
//	if getFile.Daemon {
//		if os.Getppid() != 1 {
//			os.StartProcess(os.Args[0], os.Args, &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
//			return
//		} else {
//			getFile.Server()
//		}
//	} else {
//		getFile.Server()
//	}
//}
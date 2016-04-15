package main

import (
	"os"

	//	"github.com/czxichen/AutoWork/getFile"
	"github.com/czxichen/AutoWork/getFile/center"
)

func main() {
	if os.Getppid() != 1 {
		os.StartProcess(os.Args[0], os.Args, &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
		return
	} else {
		center.Server()
	}
}

////func cmain() {
////	if os.Getppid() != 1 {
////		os.StartProcess(os.Args[0], os.Args, &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
////		return
////	} else {
////		getFile.Server()
////	}
////}

package file

import (
	"os"
	"testing"
)

func Test_SendFile(t *testing.T) {
	sFile,_:= os.Open("readFile.go")
	SendFile(os.Stdout,sFile)
	sFile.Close()
}
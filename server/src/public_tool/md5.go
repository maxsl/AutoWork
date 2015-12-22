package public_tool

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func Md5(File io.Reader) (string, error) {
	M := md5.New()
	io.Copy(M, File)
	return hex.EncodeToString(M.Sum([]byte{})), nil
}

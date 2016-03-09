package zip

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Compress interface {
	Close() error
	WriteHead(path string, info os.FileInfo) error
	Write(p []byte) (int, error)
}

func walk(path string, compresser Compress) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(path)
	}

	filepath.Walk(path, func(root string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		root = strings.Replace(root, "\\", "/", -1)
		fileroot := root
		if root = strings.TrimPrefix(root, path); root == "" {
			root = baseDir
		}
		err = compresser.WriteHead(root, info)
		if err != nil {
			return nil
		}
		F, err := os.Open(fileroot)
		if err != nil {
			return nil
		}
		io.Copy(compresser, F)
		F.Close()
		return nil
	})
	return nil
}

/*
func match(str string, regexpstr []*regexp.Regexp) bool {
	for _, v := range regexpstr {
		if !v.Match([]byte(str)) {
			continue
		}
		return true
	}
	return false
}

func getreg(regexpstr []string) []*regexp.Regexp {
	list := make([]*regexp.Regexp, 0, len(regexpstr))
	for _, v := range regexpstr {
		if !strings.HasPrefix(v, "^") {
			v = "^" + v
		}
		reg, err := regexp.Compile(v)
		if err != nil {
			continue
		}
		list = append(list, reg)
	}
	return list
}
*/

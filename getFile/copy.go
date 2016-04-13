package getFile

import (
	"io"
	"os"

	"github.com/czxichen/AutoWork/tools/zip"
)

type Exec struct {
	JobId    string   `json:jobid`
	Size     int64    `json:size`
	Files    []string `json:files`
	AbsFiles []string `json:absfiles`
}

func (self Exec) Copy() (string, error) {
	for _, path := range self.Files {
		dir, _ := SplitPath(path)
		tmp := config.TempPath + self.JobId + "/"
		os.MkdirAll(tmp+dir, 0644)
		copyFile(config.Path+path, tmp+path)
	}
	for _, path := range self.AbsFiles {
		dir, name := SplitPath(path)
		var dst string
		if dir == "" {
			dst = config.TempPath + self.JobId + "/"
		} else {
			dst = config.TempPath + self.JobId + dir + "/"
		}
		os.MkdirAll(dst, 0644)
		copyFile(path, dst+name)
	}
	p, err := self.zip()
	os.RemoveAll(config.TempPath + self.JobId)
	if err != nil {
		return "", err
	}
	return p, nil
}

func (self Exec) zip() (string, error) {
	dir := config.TempPath + self.JobId
	p := dir + ".zip"
	File, err := os.Create(p)
	if err != nil {
		return "", err
	}
	defer File.Close()
	z := zip.NewZipWriter(File)
	defer z.Close()
	err = z.Walk(dir)
	return self.JobId + ".zip", err
}

func copyFile(src, dst string) error {
	sFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sFile.Close()
	dFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dFile.Close()
	io.Copy(dFile, sFile)
	return nil
}

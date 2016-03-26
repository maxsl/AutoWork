package fileagent

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/czxichen/AutoWork/tools/zip"
)

type FilesInfo struct {
	Path     string   `json:path`
	Files    []string `json:files`
	AbsFiles []string `json:absfiles`
	Size     int64    `json:size`
	JobId    string   `json:jobid`
	Tag      string   `json:tag`
	Host     string   `json:host`
}

func (self *FilesInfo) Run() []string {
	switch self.Tag {
	case "compress":
		self.copy()
		dir, str := self.zip()
		os.RemoveAll(dir)
		if str != "" {
			return []string{str}
		}
	case "copyfiles":
		return self.copy()
	}
	return []string{}
}

func (self *FilesInfo) copy() []string {
	homeDir := tempDir + "/" + self.JobId + "/"
	var list []string
	for _, file := range self.AbsFiles {
		if path := resolvePath(file); file != "" {
			err := os.MkdirAll(homeDir+filepath.Dir(path), 0644)
			if err != nil {
				continue
			}
			name := filepath.Base(file)
			if name != path {
				path = homeDir + path + name
			}
			list = append(list, homeDir+path)
			err = copyFile(file, homeDir+path)
			if err != nil && debug {
				currentLine(err)
			}
			continue
		}
	}
	for _, file := range self.Files {
		err := os.MkdirAll(homeDir+filepath.Dir(file), 0644)
		if err != nil {
			continue
		}
		list = append(list, homeDir+file)
		err = copyFile(self.Path+file, homeDir+file)
		if err != nil && debug {
			currentLine(err)
		}
		continue
	}
	return list
}

func (self *FilesInfo) zip() (string, string) {
	dir := tempDir + "/" + self.JobId
	File, err := os.Create(tempDir + "/" + self.JobId + ".zip")
	if err != nil && debug {
		currentLine(err)
		return dir, ""
	}
	defer File.Close()
	z := zip.NewZipWriter(File)
	defer z.Close()
	err = z.Walk(dir)
	if err != nil && debug {
		currentLine(err)
	}
	return dir, dir + ".zip"
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

func resolvePath(path string) string {
	list := strings.Split(path, "/")
	if len(list) <= 1 {
		return ""
	}
	return strings.Join(list[1:], "/")
}

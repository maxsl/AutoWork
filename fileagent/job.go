package fileagent

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Job struct {
	JobId string   `json:id`
	Name  []string `json:name`
	Path  string   `json:path`
	Tag   string   `json:"tag"`
	temp  string   `json:tempfile`
}

func (self *Job) Init() {
	self.Path = filepath.ToSlash(self.Path)
	if !strings.HasSuffix(self.Path, "/") {
		self.Path = self.Path + "/"
	}
}
func (self *Job) GetFilesInfo() *FilesInfo {
	self.Init()
	var Files *FilesInfo = &FilesInfo{Path: self.Path, JobId: self.JobId, Tag: self.Tag}
	for _, file := range self.Name {
		file := filepath.ToSlash(file)
		if filepath.IsAbs(file) {
			info, err := os.Stat(file)
			if err != nil {
				fmt.Println(err)
				if debug {
					currentLine(err)
				}
				continue
			}
			Files.Size += info.Size()
			Files.AbsFiles = append(Files.AbsFiles, file)
			continue
		}
		info, err := os.Stat(self.Path + file)
		if err != nil {
			if debug {
				currentLine(err)
			}
			continue
		}
		if info.IsDir() {
			continue
		}
		Files.Size += info.Size()
		Files.Files = append(Files.Files, file)
	}
	return Files
}

func (self *Job) PrintFiles() map[string][]string {
	self.Init()
	fmt.Println(self.Path)
	var m map[string][]string = make(map[string][]string)
	for _, v := range self.Name {
		var list []string
		v = filepath.ToSlash(v)
		if filepath.IsAbs(v) {
			_, err := os.Stat(v)
			if err != nil {
				if debug {
					currentLine(err)
				}
				continue
			}
			m[v] = append(list, v)
			continue
		}
		path := self.Path + v
		info, err := os.Stat(path)
		if err != nil {
			if debug {
				currentLine(err)
			}
			continue
		}
		if !info.IsDir() {
			list = append(list, path)
			m[path] = list
			continue
		}
		l, err := ioutil.ReadDir(path)
		if err != nil {
			if debug {
				currentLine(err)
			}
			continue
		}
		for _, name := range l {
			if !name.IsDir() {
			}
			list = append(list, name.Name())
		}
		m[path] = list
	}
	return m
}

package findFile

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Event struct {
	Date     int      `json:date`
	NameReg  string   `json:Namereg`
	RootPath string   `json:rootpath`
	PathList []string `json:pathlist`
}

func (self *Event) init() {
	self.RootPath = filepath.ToSlash(self.RootPath)
	if !strings.HasSuffix(self.RootPath, "/") {
		self.RootPath += "/"
	}
}

func (self *Event) NewEvent() *job {
	self.init()
	var j job = job{Root: self.RootPath}
	for _, path := range self.PathList {
		path = filepath.ToSlash(path)
		if filepath.IsAbs(path) {
			info, err := os.Lstat(path)
			if err != nil {
				j.FaildPath = append(j.FaildPath, path)
				continue
			}
			if info.IsDir() {
				j.ABSDirectory = append(j.ABSDirectory, path)
			} else {
				j.ABSFile = append(j.ABSFile, path)
			}
		} else {
			info, err := os.Lstat(j.Root + path)
			if err != nil {
				j.FaildPath = append(j.FaildPath, path)
				continue
			}
			if info.IsDir() {
				j.Directory = append(j.Directory, path)
			} else {
				j.File = append(j.File, path)
			}
		}
	}
	return &j
}

type match struct {
	Unixtime int64
	Less     bool
	Reg      *regexp.Regexp
}

func (self match) Match(path string) (bool, int64) {
	info, err := os.Lstat(path)

	if err != nil {
		return false, 0
	}

	if self.Reg != nil && !self.Reg.MatchString(info.Name()) {
		return b, 0
	}

	if self.Unixtime != 0 {
		if self.Less {
			if info.ModTime().Unix() >= self.Unixtime {
				return b, 0
			}
		} else {
			if info.ModTime().Unix() <= self.Unixtime {
				return b, 0
			}
		}
	}
	return true, info.Size()
}

func newMatch(date int64, reg string) (m match) {
	if reg != "" {
		if strings.Index(reg, "*") == 0 {
			reg = "." + reg
		} else {
			reg = "^" + reg
		}
		reg += "$"
		Reg, err := regexp.Compile(reg)
		if err == nil {
			m.Reg = Reg
		}
	}
	if date != 0 {
		if date < 0 {
			m.Unixtime = time.Now().Unix() + date*24*60*60
		} else {
			m.Unixtime = time.Now().Unix() - date*24*60*60
			m.Less = true
		}
	}
	return
}

type job struct {
	m            match
	Root         string   `json:root`
	Size         int64    `json:size`
	File         []string `json:file`
	ABSFile      []string `json:absfile`
	Directory    []string `json:directory`
	ABSDirectory []string `json:absdirectory`
	FaildPath    []string `json:faildpath`
}

func (self *job) Walk() {

}

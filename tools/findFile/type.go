package findFile

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type match struct {
	unixtime int64
	less     bool
	reg      *regexp.Regexp
}

func (self match) MatchName(info os.FileInfo) bool {
	if self.reg != nil && !self.reg.MatchString(info.Name()) {
		return false
	}
	return true
}

func (self match) MatchDate(info os.FileInfo) bool {
	if self.unixtime != 0 {
		if self.less {
			if info.ModTime().Unix() >= self.unixtime {
				return false
			}
		} else {
			if info.ModTime().Unix() <= self.unixtime {
				return false
			}
		}
	}
	return true
}

func (self match) MatchInfo(info os.FileInfo) (bool, int64) {
	var b bool
	if b = self.MatchName(info); !b {
		return false, 0
	}
	if b = self.MatchDate(info); !b {
		return false, 0
	}
	return true, info.Size()
}

func (self match) MatchPath(path string) (bool, int64) {
	info, err := os.Lstat(path)

	if err != nil {
		return false, 0
	}

	return self.MatchInfo(info)
}

func (self match) Walk(dirPath string) (filelist []string, size int64) {
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if b, s := self.MatchInfo(info); b {
			//path = filepath.ToSlash(strings.TrimPrefix(path, dirPath))
			path = filepath.ToSlash(path)
			filelist = append(filelist, path)
			size += s
		}
		return nil
	})
	return
}

type Event struct {
	Id       string   `json:id`
	Date     int64    `json:date`
	NameReg  string   `json:Namereg`
	PathList []string `json:pathlist`
	RootPath string   `json:rootpath`
}

func (self Event) NewJob() job {
	self.RootPath = filepath.ToSlash(self.RootPath)
	if !strings.HasSuffix(self.RootPath, "/") {
		self.RootPath += "/"
	}
	var j job = job{Id: self.Id, Date: self.Date, NameReg: self.NameReg, Root: self.RootPath}
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
	return j
}

type jobResult struct {
	Id        string   `json:id`
	Size      int64    `json:size`
	RootPath  string   `json:rootpath`
	Files     []string `json:files`
	ABSFiles  []string `json:absfiles`
	FaildPath []string `json:faildpath`
}

type job struct {
	Id           string   `json:id`
	Date         int64    `json:date`
	NameReg      string   `json:namereg`
	Root         string   `json:root`
	File         []string `json:file`
	ABSFile      []string `json:absfile`
	Directory    []string `json:directory`
	ABSDirectory []string `json:absdirectory`
	FaildPath    []string `json:faildpath`
}

func (self job) Walk() jobResult {
	err := os.Chdir(self.Root)
	m := NewMatch(self.Date, self.NameReg)
	jr := jobResult{Id: self.Id, RootPath: self.Root}

	if err == nil {
		for _, path := range self.File {
			b, size := m.MatchPath(path)
			if b {
				jr.Files = append(jr.Files, path)
				jr.Size += size
			}
		}
		for _, dir := range self.Directory {
			list, size := m.Walk(dir)
			jr.Files = append(jr.Files, list...)
			jr.Size += size
		}
	}

	for _, path := range self.ABSFile {
		b, size := m.MatchPath(path)
		if b {
			jr.ABSFiles = append(jr.ABSFiles, path)
			jr.Size += size
		}
	}

	for _, dir := range self.ABSDirectory {
		list, size := m.Walk(dir)
		jr.ABSFiles = append(jr.ABSFiles, list...)
		jr.Size += size
	}

	return jr
}

package getFile

import (
	"os"
	"path/filepath"
	"strings"

	fp "github.com/czxichen/Goprograme/filepath"
)

const (
	Compress          = 1
	FullDirCompress   = 2
	RegAndDate        = 3
	FullDirRegAndDate = 4
)

type FilesInfo struct {
	JobId    string   `json:jobid`
	Tag      int      `json:tag`
	Date     int64    `json:date`
	Reg      string   `json:reg`
	Files    []string `json:files`
	AbsFiles []string `json:absfiles`
	Dirs     []string `json:dirs`
}

func (F FilesInfo) GetFilesInfo() (l Exec) {
	if config.Debug {
		Log.Println("FilesInfo:", F)
	}
	if F.Tag == 0 {
		return
	}
	l.JobId = F.JobId
	switch F.Tag {
	case Compress:
		l.Size += size(config.Path, F.Files)
		l.Size += size("", F.AbsFiles)
		l.AbsFiles = F.AbsFiles
		l.Files = F.Files
		return
	case FullDirCompress:
		for _, dir := range F.Dirs {
			files, length := walkDir(dir)
			l.Files = files
			l.Size += length
		}
		return
	case RegAndDate:
		if F.Date == 0 && len(F.Reg) == 0 {
			return
		}
		for _, dir := range F.Dirs {
			dir = config.Path + dir
			f := fp.FindFiles{dir, false, false}
			files, size, err := f.NewFind(F.Date, F.Reg)
			if err != nil {
				continue
			}
			for _, v := range files {
				l.Files = append(l.Files, strings.TrimPrefix(v, config.Path))
			}
			l.Size += size
		}
		return
	case FullDirRegAndDate:
		if F.Date == 0 && len(F.Reg) == 0 {
			return
		}
		for _, dir := range F.Dirs {
			dir = config.Path + dir
			f := fp.FindFiles{dir, true, false}
			files, size, err := f.NewFind(F.Date, F.Reg)
			if err != nil {
				continue
			}
			for _, v := range files {
				l.Files = append(l.Files, strings.TrimPrefix(v, config.Path))
			}
			l.Size += size
		}
		return
	}
	return
}

type Job struct {
	JobId  string   `json:id`
	Tag    int      `json:"tag"`
	Files  []string `json:files`
	Date   int64    `json:date`
	Regexp string   `json:regexp`
}

func (j Job) Start() FilesInfo {
	if j.Tag <= 0 && j.Tag > 4 || len(j.JobId) <= 0 {
		return FilesInfo{}
	}
	var fi FilesInfo = FilesInfo{JobId: j.JobId, Tag: j.Tag, Date: j.Date, Reg: j.Regexp}
	for _, file := range j.Files {
		file := filepath.ToSlash(file)
		if filepath.IsAbs(file) {
			Info, err := os.Lstat(file)
			if err != nil {
				continue
			}
			//不支持绝对路径的目录
			if Info.IsDir() {
				continue
			}
			fi.AbsFiles = append(fi.AbsFiles, file)
			continue
		}

		Info, err := os.Lstat(config.Path + file)
		if err != nil {
			continue
		}
		if Info.IsDir() {
			fi.Dirs = append(fi.Dirs, file)
			continue
		}
		fi.Files = append(fi.Files, file)
	}
	return fi
}

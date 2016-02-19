package handler

import (
	"net/http"
	"os/exec"
	"strings"

	"github.com/AutoWork/tools/md5"
	"github.com/AutoWork/tools/wget"
	"github.com/AutoWork/tools/zip"
)

var tmpfileDir string = "tmpfile/"

func Commond(w http.ResponseWriter, r *http.Request) {
	action := r.FormValue("action")
	switch action {
	case "custom":
		custom(w, r)
	case "system":
		system(w, r)
	default:
		Default(w, 601)
	}
}

func system(w http.ResponseWriter, r *http.Request) {
	cmd := r.FormValue("commond")
	list := strings.Split(cmd, " ")
	if len(list) <= 0 {
		Default(w, 602)
		return
	}
	c := exec.Command(list[0], list[1:]...)
	buf, err := c.Output()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(buf)
}

func custom(w http.ResponseWriter, r *http.Request) {
	cmd := r.FormValue("commond")
	switch cmd {
	case "wget":
		wgetFile(w, r)
	case "unzip":
		unzipFile(w, r)
	default:
		Default(w, 602)
	}
}

func checkmd5(w http.ResponseWriter, r *http.Request) {
	filepath := r.FormValue("filepath")
	str, err := md5.Md5(filepath)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write([]byte(str))
}

func unzipFile(w http.ResponseWriter, r *http.Request) {
	filepath := r.FormValue("filepath")
	if zip.CheckValidZip(filepath) {
		dscDir := r.FormValue("descdir")
		if dscDir == "" {
			Default(w, 602)
			return
		}
		err := zip.Unzip(filepath, dscDir, nil)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write([]byte("unzip ok"))
		return
	}
	Default(w, 603)
}

//post 文件下载路径,和要保存的名称.
func wgetFile(w http.ResponseWriter, r *http.Request) {
	fileurl := r.FormValue("fileurl")
	filename := r.FormValue("filename")
	if filename == "" || fileurl == "" {
		Default(w, 602)
		return
	}
	filename = tmpfileDir + filename
	_, err := wget.Wget(fileurl, filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	str, err := md5.Md5(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write([]byte(str))
}

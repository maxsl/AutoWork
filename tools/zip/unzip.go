package zip

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"strings"
)

type log interface {
	PrintfI(formate string, v ...interface{})
	PrintfE(formate string, v ...interface{})
}

func Unzip(filename, dir string, Log log) error {
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}
	File, err := zip.OpenReader(filename)
	if err != nil {
		return errors.New("Error Open zip faild: " + err.Error())
	}
	defer File.Close()
	for _, v := range File.File {
		err := createFile(v, dir)
		if err != nil {
			Log.PrintfE("unzip file err %v \n", err)
			continue
		}
		os.Chtimes(v.Name, v.ModTime(), v.ModTime())
		os.Chmod(v.Name, v.Mode())
		Log.PrintfI("unzip %s %s\n", filename, v.Name)
	}
	return nil
}

func createFile(v *zip.File, dscDir string) error {
	v.Name = dscDir + v.Name
	info := v.FileInfo()
	if info.IsDir() {
		err := os.MkdirAll(v.Name, v.Mode())
		if err != nil {
			return errors.New("Error Create direcotry" + v.Name + "faild: " + err.Error())
		}
		return nil
	}
	srcFile, err := v.Open()
	if err != nil {
		return errors.New("Error Read from zip faild: " + err.Error())
	}
	defer srcFile.Close()
	newFile, err := os.Create(v.Name)
	if err != nil {
		return errors.New("Error Create file faild: " + err.Error())
	}
	defer newFile.Close()
	io.Copy(newFile, srcFile)
	return nil
}

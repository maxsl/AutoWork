package zip

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"strings"
)

func Unzip(filename, dir string, Log func(format string, v ...interface{}) (int, error)) error {
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
			if Log != nil {
				Log("unzip file err %v \n", err)
			}
			return err
		}
		os.Chtimes(v.Name, v.ModTime(), v.ModTime())
		os.Chmod(v.Name, v.Mode())
		if Log != nil {
			Log("unzip %s %s\n", filename, v.Name)
		}
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

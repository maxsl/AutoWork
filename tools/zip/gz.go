package zip

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func Gzip(filepath, filename string, Log func(format string, v ...interface{}) (int, error)) error {
	File, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer File.Close()
	gw := gzip.NewWriter(File)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	return walk(filepath, tw, Log)
}

func walk(path string, tw *tar.Writer, Log func(format string, v ...interface{}) (int, error)) error {
	path = strings.Replace(path, "\\", "/", -1)
	info, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	index := strings.Index(path, "/")
	list := strings.Join(strings.Split(path, "/")[index:], "/")
	for _, v := range info {
		if v.IsDir() {
			head := tar.Header{Name: list + v.Name(), Typeflag: tar.TypeDir, ModTime: v.ModTime()}
			tw.WriteHeader(&head)
			if Log != nil {
				Log("create directory: %s\n", list+v.Name())
			}
			walk(path+v.Name(), tw, Log)
			continue
		}
		F, err := os.Open(path + v.Name())
		if err != nil {
			fmt.Println("打开文件%s失败.", err)
			continue
		}
		head := tar.Header{Name: list + v.Name(), Size: v.Size(), Mode: int64(v.Mode()), ModTime: v.ModTime()}
		tw.WriteHeader(&head)
		io.Copy(tw, F)
		F.Close()
		if Log != nil {
			Log("create directory: %s\n", list+v.Name())
		}
	}
	return nil
}

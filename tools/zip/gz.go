package zip

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func Gzip(filepath, filename string, exclude []string, Log func(format string, v ...interface{}) (int, error)) error {
	File, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer File.Close()
	gw := gzip.NewWriter(File)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	excludereg := getreg(exclude)

	return walk(filepath, tw, excludereg, Log)
}

func walk(path string, tw *tar.Writer, excludereg []*regexp.Regexp, Log func(format string, v ...interface{}) (int, error)) error {
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
		if match(v.Name(), excludereg) {
			continue
		}
		if v.IsDir() {
			head := tar.Header{Name: list + v.Name(), Typeflag: tar.TypeDir, ModTime: v.ModTime()}
			tw.WriteHeader(&head)
			if Log != nil {
				Log("create directory: %s\n", list+v.Name())
			}
			walk(path+v.Name(), tw, excludereg, Log)
			continue
		}
		F, err := os.Open(path + v.Name())
		if err != nil {
			fmt.Println("open file %s faild.", err)
			continue
		}
		head := tar.Header{Name: list + v.Name(), Size: v.Size(), Mode: int64(v.Mode()), ModTime: v.ModTime()}
		tw.WriteHeader(&head)
		io.Copy(tw, F)
		F.Close()
		if Log != nil {
			Log("create file: %s\n", list+v.Name())
		}
	}
	return nil
}

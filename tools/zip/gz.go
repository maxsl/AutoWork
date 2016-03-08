package zip

import (
	"archive/tar"
	"compress/gzip"
	"os"
)

type TgzWirter struct {
	tar *tar.Writer
	gz  *gzip.Writer
}

func NewTgzWirter(File *os.File) *TgzWirter {
	gw := gzip.NewWriter(File)
	tw := tar.NewWriter(gw)
	return &TgzWirter{tw, gw}
}

func (self *TgzWirter) Close() error {
	self.tar.Close()
	return self.gz.Close()
}

func (self *TgzWirter) WriteHead(path string, info os.FileInfo) error {
	if path == "." || path == ".." {
		return nil
	}
	head, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	head.Name = path
	if info.IsDir() {
		head.Name += "/"
	}
	return self.tar.WriteHeader(head)
}

func (self *TgzWirter) Write(p []byte) (int, error) {
	return self.tar.Write(p)
}

func (self *TgzWirter) Walk(source string) error {
	return walk(source, self)
}

/*
func Gzip(source, filename string) error {
	zipfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	gw := gzip.NewWriter(zipfile)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = path

		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.CopyN(tw, file, info.Size())
		return err
	})
	return err
}

func Gzip_exclude(filepath, filename string, exclude []string, Log func(format string, v ...interface{}) (int, error)) error {
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
		head, err := tar.FileInfoHeader(v, "")
		if err != nil {
			continue
		}
		head.Name = list + v.Name()
		if v.IsDir() {
			tw.WriteHeader(head)
			if Log != nil {
				Log("create directory: %s\n", list+v.Name())
			}
			walk(path+v.Name(), tw, excludereg, Log)
			continue
		}
		F, err := os.Open(path + v.Name())
		if err != nil {
			continue
		}
		tw.WriteHeader(head)
		io.Copy(tw, F)
		F.Close()
		if Log != nil {
			Log("create file: %s\n", list+v.Name())
		}
	}
	return nil
}
*/

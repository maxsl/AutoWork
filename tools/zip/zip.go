package zip

import (
	"archive/zip"
	"io"
	"os"
	"time"
)

func NewZipWriter(File *os.File) *ZipWrite {
	zipwrite := zip.NewWriter(File)
	return &ZipWrite{zone: 8, zw: zipwrite, file: File}
}

type ZipWrite struct {
	zone   int64
	zw     *zip.Writer
	writer io.Writer
	file   *os.File
}

func (self *ZipWrite) Close() error {
	return self.zw.Close()
}

func (self *ZipWrite) WriteHead(path string, info os.FileInfo) error {
	if path == "." || path == ".." {
		return nil
	}
	head, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		head.Method = zip.Deflate
	}
	head.Name = path
	if info.IsDir() {
		head.Name += "/"
	}
	head.SetModTime(time.Unix(info.ModTime().Unix()+self.zone*60*60, 0))
	write, err := self.zw.CreateHeader(head)
	if err != nil {
		return err
	}
	self.writer = write
	return nil
}

func (self *ZipWrite) Write(p []byte) (int, error) {
	return self.writer.Write(p)
}

func (self *ZipWrite) Walk(source string) error {
	return walk(source, self)
}

/*
const zone int64 = +8

func Zip(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			header.Method = zip.Deflate
		}
		header.SetModTime(time.Unix(info.ModTime().Unix()+(zone*60*60), 0))
		header.Name = path
		writer, err := archive.CreateHeader(header)
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
		_, err = io.Copy(writer, file)
		return err
	})
}
*/

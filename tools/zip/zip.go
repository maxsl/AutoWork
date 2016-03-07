package zip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"time"
)

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

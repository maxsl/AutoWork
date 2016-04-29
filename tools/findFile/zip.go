package findFile

import (
	"io"
	"os"
	"path/filepath"

	"github.com/czxichen/AutoWork/tools/zip"
)

func Compress(jr jobResult) error {
	sFile, err := os.Create("tmp/" + jr.Id + ".zip")
	if err != nil {
		return err
	}
	defer sFile.Close()

	zw := zip.NewZipWriter(sFile)
	defer zw.Close()

	for _, path := range jr.Files {
		zwFile(jr.RootPath, path, zw)
	}

	for _, path := range jr.ABSFiles {
		zwFile("", path, zw)
	}
	return nil
}

func zwFile(RootPath, path string, zw *zip.ZipWrite) error {
	d := filepath.Dir(path)

	info, err := os.Lstat(RootPath + d)
	if err != nil {
		return err
	}

	err = zw.WriteHead(d, info)
	if err != nil {
		return err
	}

	File, err := os.Open(RootPath + path)
	if err != nil {
		return err
	}

	info, err = File.Stat()
	if err != nil {
		return err
	}

	err = zw.WriteHead(path, info)
	if err != nil {
		return err
	}

	io.CopyN(zw, File, info.Size())

	return nil
}

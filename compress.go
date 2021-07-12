package goutils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Zip 压缩目录。
func Zip(srcFile string, destZip string, containSelf bool) error {
	zipfile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	srcFile = PathNormalized(srcFile)
	srcSelfDirName := filepath.Base(srcFile)

	err = filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		d := filepath.Dir(srcFile)

		if containSelf {
			header.Name = UnixPathNormalized(RemovePathSeparatorPrefix(strings.TrimPrefix(path, d)))
		} else {
			header.Name = UnixPathNormalized(RemovePathSeparatorPrefix(strings.TrimPrefix(path, fmt.Sprintf("%s%s%s", d, string(os.PathSeparator), srcSelfDirName))))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	return err
}

// Unzip 解压文件。
func Unzip(zipFile string, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	var decodeName string

	addFileFunc := func(f *zip.File, fpath string) error {
		if f.FileInfo().IsDir() {
			return os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}

		return nil
	}

	for _, f := range zipReader.File {
		if f.Flags == 0 {
			i := bytes.NewReader([]byte(f.Name))
			decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
			content, _ := ioutil.ReadAll(decoder)
			decodeName = string(content)
		} else {
			decodeName = f.Name
		}

		fpath := filepath.Join(destDir, decodeName)

		err := addFileFunc(f, fpath)

		if err != nil {
			return err
		}
	}
	return nil
}

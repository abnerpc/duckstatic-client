package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ZipFileName holds the default zipped filename
const ZipFileName = "static.zip"

// Zipit compress the content of source path in a zipped file
// The target parameter is used as the name for the first folder or file
func Zipit(source, target string) error {

	zipfile, err := os.Create(ZipFileName)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	var oneFile = false
	if info.IsDir() {
		if target != "" {
			baseDir = target
		} else {
			baseDir = filepath.Base(source)
		}
	} else {
		oneFile = true
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
			if oneFile && target != "" {
				fileName := header.Name
				header.Name = strings.Join([]string{target, strings.Split(fileName, ".")[1]}, ".")
			}
		}

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

	return err
}

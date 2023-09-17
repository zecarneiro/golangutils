package jnoronhautils

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type IFileInfo struct {
	File      string
	Directory string
}

type IDirectoryInfo struct {
	Files       []string
	Directories []string
}

const (
	FileTypeNone         = 0
	FileTypeDirectory    = 1
	FileTypeFile         = 2
	FileTypeSymbolicLink = 3
)

func FileType(fileName string) (int, error) {
	var typeFile int
	file, err := os.Open(fileName)
	if err != nil {
		return FileTypeNone, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return FileTypeNone, err
	}
	if fileInfo.IsDir() {
		typeFile = FileTypeDirectory
	} else {
		typeFile = FileTypeFile
	}
	defer file.Close()
	return typeFile, nil
}

func ReadDir(dir string) ([]IFileInfo, error) {
	var filesList []IFileInfo
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		ErrorLog(err.Error(), false)
		return []IFileInfo{}, err
	} else {
		for _, file := range files {
			if file.IsDir() {
				filesList = append(filesList, IFileInfo{Directory: file.Name()})
			} else {
				filesList = append(filesList, IFileInfo{File: file.Name()})
			}
		}
	}
	return filesList, nil
}

func ReadDirRecursive(dir string) (IDirectoryInfo, error) {
	files := IDirectoryInfo{Directories: []string{}, Files: []string{}}
	err := filepath.Walk(dir,
		func(path string, _ os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path != "." {
				info, err := FileType(path)
				if err != nil {
					return err
				}
				if info == FileTypeDirectory {
					files.Directories = append(files.Directories, path)
				} else {
					files.Files = append(files.Files, path)
				}
			}
			return nil
		},
	)
	if err != nil {
		return IDirectoryInfo{}, err
	}
	return files, nil
}

func IsDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// read in ONLY one file
	_, err = f.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

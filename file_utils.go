package jnoronha_golangutils

import (
	"io"
	"jnoronha_golangutils/entities"
	"os"
	"path/filepath"
)

func ResolvePath(path string) string {
	if len(path) > 0 {
		return filepath.FromSlash(path)
	}
	return ""
}

func ReadFile(file string) (string, error) {
	body, err := os.ReadFile(ResolvePath(file))
	if err != nil {
		return "", err
	}
	return string(body), err
}

func FileType(fileName string) (int, error) {
	var typeFile int
	file, err := os.Open(fileName)
	if err != nil {
		return entities.FileTypeNone, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return entities.FileTypeNone, err
	}
	if fileInfo.IsDir() {
		typeFile = entities.FileTypeDirectory
	} else {
		typeFile = entities.FileTypeFile
	}
	defer file.Close()
	return typeFile, nil
}

func ReadDir(dir string) ([]entities.FileInfo, error) {
	var filesList []entities.FileInfo
	files, err := os.ReadDir(dir)
	if err != nil {
		ErrorLog(err.Error(), false)
		return []entities.FileInfo{}, err
	} else {
		for _, file := range files {
			if file.IsDir() {
				filesList = append(filesList, entities.FileInfo{Directory: file.Name()})
			} else {
				filesList = append(filesList, entities.FileInfo{File: file.Name()})
			}
		}
	}
	return filesList, nil
}

func ReadDirRecursive(dir string) (entities.DirectoryInfo, error) {
	files := entities.DirectoryInfo{Directories: []string{}, Files: []string{}}
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
				if info == entities.FileTypeDirectory {
					files.Directories = append(files.Directories, path)
				} else {
					files.Files = append(files.Files, path)
				}
			}
			return nil
		},
	)
	if err != nil {
		return entities.DirectoryInfo{}, err
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

func FileExist(file string) bool {
	_, err := os.Stat(ResolvePath(file))
	return err == nil
}

func GetCurrentDir() string {
	path, err := os.Getwd()
	if err != nil {
		ErrorLog(err.Error(), false)
		path = ""
	}
	return path
}

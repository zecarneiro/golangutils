package file

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"golangutils/pkg/models"
	"golangutils/pkg/slice"
)

func GetCurrentDir() (string, error) {
	return os.Getwd()
}

func ReadDir(dir string) (models.FileInfo, error) {
	dir = ResolvePath(dir)
	var filesData models.FileInfo
	files, err := os.ReadDir(dir)
	if err != nil {
		return models.FileInfo{}, err
	} else {
		for _, file := range files {
			if file.IsDir() {
				filesData.Directories = append(filesData.Directories, file.Name())
			} else {
				filesData.Files = append(filesData.Files, file.Name())
			}
		}
	}
	filesData.Directories = slice.FilterArray(filesData.Directories, func(path string) bool { return path != dir })
	filesData.Files = slice.FilterArray(filesData.Files, func(path string) bool { return path != dir })
	return filesData, nil
}

func ReadDirRecursive(dir string) (models.FileInfo, error) {
	dir = ResolvePath(dir)
	files := models.FileInfo{Directories: []string{}, Files: []string{}}
	err := filepath.Walk(dir,
		func(path string, _ os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path != "." && path != ".." {
				info := Type(path)
				if info == Directory {
					files.Directories = append(files.Directories, path)
				} else {
					files.Files = append(files.Files, path)
				}
			}
			return nil
		},
	)
	if err != nil {
		return models.FileInfo{}, err
	}
	files.Directories = slice.FilterArray(files.Directories, func(path string) bool { return path != dir })
	files.Files = slice.FilterArray(files.Files, func(path string) bool { return path != dir })
	return files, nil
}

func IsDirEmpty(name string) (bool, error) {
	fs, err := os.Open(ResolvePath(name))
	if err != nil {
		return false, err
	}
	defer fs.Close()

	// read in ONLY one file
	_, err = fs.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func IsDir(file string) bool {
	return FileExist(file) && Type(file) == Directory
}

func CopyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo
	src = ResolvePath(src)
	dst = ResolvePath(dst)

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func GetExecutableDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return Dirname(execPath), nil
}

func CreateDirectory(dir string, recursive bool) error {
	dir = ResolvePath(dir)
	var err error
	if IsDir(dir) {
		err = nil
	} else {
		if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
			if recursive {
				err = os.MkdirAll(dir, os.ModePerm)
			} else {
				err = os.Mkdir(dir, os.ModePerm)
			}
		}
	}
	return err
}

func DeleteDirectory(directory string) error {
	directory = ResolvePath(directory)
	err := os.RemoveAll(directory)
	if err != nil {
		return err
	}
	return nil
}

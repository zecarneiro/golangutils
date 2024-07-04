package jnoronhautils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"jnoronhautils/entities"
	"log"
	"os"
	"path"
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

func ReadFileLineByLine(filePath string, callback func(string, error)) {
	file, err := os.Open(ResolvePath(filePath))
	if err != nil {
		callback("", err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		callback(line, nil)
	}
	if scanner.Err() != nil {
		callback("", scanner.Err())
	}
}

func WriteFile(file string, data string, isAppend bool) bool {
	var status = true
	var fileStream *os.File
	var err error
	file = ResolvePath(file)
	if isAppend && FileExist(file) {
		fileStream, err = os.OpenFile("lines", os.O_APPEND|os.O_WRONLY, 0644)
	} else {
		fileStream, err = os.Create(file)
	}
	if err != nil {
		ErrorLog(err.Error(), false)
		fileStream.Close()
		return false
	}
	if status {
		_, err = fmt.Fprintln(fileStream, data)
		if err != nil {
			ErrorLog(err.Error(), false)
			fileStream.Close()
			return false
		}
	}
	fileStream.Close()
	return true
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

func ReadDir(dir string) (entities.FileInfo, error) {
	var filesData entities.FileInfo
	files, err := os.ReadDir(dir)
	if err != nil {
		ErrorLog(err.Error(), false)
		return entities.FileInfo{}, err
	} else {
		for _, file := range files {
			if file.IsDir() {
				filesData.Directory = append(filesData.Directory, file.Name())
			} else {
				filesData.File = append(filesData.File, file.Name())
			}
		}
	}
	return filesData, nil
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

func ReadJsonFile[T any](jsonFile string) (T, error) {
	data, err := ReadFile(jsonFile)
	if err != nil {
		ErrorLog(err.Error(), false)
	}
	return StringToObject[T](data)
}

func WriteJsonFile[T any](jsonFile string, data T) bool {
	dataStr, err := ObjectToString(data)
	if err != nil {
		ErrorLog(err.Error(), false)
		return false
	}
	return WriteFile(jsonFile, dataStr, false)
}

func DeleteDirectory(directory string) bool {
	directory = ResolvePath(directory)
	err := os.RemoveAll(directory)
	if err != nil {
		ErrorLog(err.Error(), false)
		return false
	}
	return true
}

func DeleteFile(file string) bool {
	file = ResolvePath(file)
	err := os.Remove(file)
	if err != nil {
		ErrorLog(err.Error(), false)
		return false
	}
	return true
}

func CopyFile(src string, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo
	src = ResolvePath(src)
	dst = ResolvePath(dst)
	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
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

func GetExecutableDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		ErrorLog(err.Error(), false)
		dir = ""
	}
	return dir
}

func ReadFileInByte(filename string) []byte {
	file, err := os.Open(ResolvePath(filename))
	if err != nil {
		file.Close()
		ErrorLog(err.Error(), false)
		return []byte{}
	}
	defer file.Close()

	// Get the file size
	stat, err := file.Stat()
	if err != nil {
		ErrorLog(err.Error(), false)
		return []byte{}
	}

	// Read the file into a byte slice
	byteArr := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(byteArr)
	if err != nil && err != io.EOF {
		ErrorLog(err.Error(), false)
		return []byte{}
	}
	return byteArr
}

func CreateDirectory(dir string, recursive bool) {
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		var err error
		if recursive {
			err = os.MkdirAll(dir, os.ModePerm)
		} else {
			err = os.Mkdir(dir, os.ModePerm)
		}
		if err != nil {
			log.Println(err)
		}
	}
}

func GetDrives() (r []string) {
    for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ"{
		driveDir := string(drive)+":\\"
		f, err := os.Open(driveDir)
        if err == nil {
			r = append(r, driveDir)
            f.Close()
        }
    }
    return
}

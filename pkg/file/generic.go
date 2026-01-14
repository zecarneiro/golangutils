package file

import (
	"fmt"
	"os"
	"path/filepath"

	"golangutils/pkg/str"
)

func ResolvePath(path ...string) string {
	if len(path) > 0 {
		return filepath.FromSlash(JoinPath(path...))
	}
	return ""
}

func JoinPath(elem ...string) string {
	return filepath.Join(elem...)
}

func FileExist(file string) bool {
	_, err := os.Stat(ResolvePath(file))
	return err == nil
}

func GetDrives() (r []string) {
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		driveDir := string(drive) + ":\\"
		f, err := os.Open(driveDir)
		if err == nil {
			r = append(r, driveDir)
			f.Close()
		}
	}
	return
}

func Basename(path string) string {
	return filepath.Base(path)
}

func Dirname(path string) string {
	return filepath.Dir(path)
}

func Type(fileName string) FileType {
	var typeFile FileType
	file, err := os.Open(ResolvePath(fileName))
	if err != nil {
		return GetFileTypeFromValue(-1)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return GetFileTypeFromValue(-1)
	}
	if fileInfo.IsDir() {
		typeFile = GetFileTypeFromValue(Directory.Int())
	} else {
		typeFile = GetFileTypeFromValue(File.Int())
	}
	defer file.Close()
	return typeFile
}

func Move(src string, dstBaseDir string) error {
	if src == "." || src == ".." {
		return fmt.Errorf("invalid given source file/dir")
	}
	if src == "." || src == ".." || dstBaseDir == "." || dstBaseDir == ".." {
		return fmt.Errorf("invalid given destination file/dir")
	}
	src = ResolvePath(src)
	dstBaseDir = ResolvePath(dstBaseDir)
	if str.HasLastChar(src, "/") {
		src = str.DeleteLastChar(src)
	} else if str.HasLastChar(src, "\\") {
		src = str.DeleteLastChar(src)
	}
	if str.HasLastChar(dstBaseDir, "/") {
		dstBaseDir = str.DeleteLastChar(dstBaseDir)
	} else if str.HasLastChar(dstBaseDir, "\\") {
		dstBaseDir = str.DeleteLastChar(dstBaseDir)
	}
	if !FileExist(src) {
		return fmt.Errorf("%s: not found source file/dir", src)
	}
	basename := Basename(src)
	if IsFile(src) {
		err := CopyFile(src, ResolvePath(dstBaseDir, basename))
		if err != nil {
			return err
		}
		return DeleteFile(src)
	} else if IsDir(src) {
		err := CopyDir(src, ResolvePath(dstBaseDir, basename))
		if err != nil {
			return err
		}
		return DeleteDirectory(src)
	}
	return fmt.Errorf("%s: invalid type of file/dir", src)
}

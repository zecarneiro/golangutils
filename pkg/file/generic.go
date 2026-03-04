package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golangutils/pkg/enums"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"
)

func ResolvePath(path string) string {
	return filepath.FromSlash(path)
}

func JoinPath(elem ...string) string {
	return ResolvePath(filepath.Join(elem...))
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

func Type(fileName string) enums.FileType {
	var typeFile enums.FileType
	file, err := os.Open(ResolvePath(fileName))
	if err != nil {
		return enums.GetFileTypeFromValue(-1)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return enums.GetFileTypeFromValue(-1)
	}
	if fileInfo.IsDir() {
		typeFile = enums.GetFileTypeFromValue(enums.Directory.Int())
	} else {
		typeFile = enums.GetFileTypeFromValue(enums.File.Int())
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
		err := CopyFile(src, JoinPath(dstBaseDir, basename))
		if err != nil {
			return err
		}
		return DeleteFile(src)
	} else if IsDir(src) {
		err := CopyDir(src, JoinPath(dstBaseDir, basename))
		if err != nil {
			return err
		}
		return DeleteDirectory(src)
	}
	return fmt.Errorf("%s: invalid type of file/dir", src)
}

func GetFullPath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return path, err
	}
	absPath, err = filepath.EvalSymlinks(absPath)
	if err != nil {
		return path, err
	}
	return absPath, nil
}

func GetRelativePath(absPath string, basePath string) (string, error) {
	relPath := ResolvePath(absPath)
	basePath = ResolvePath(basePath)
	hasPrefix := false
	if (platform.IsWindows() && strings.HasPrefix(relPath, basePath+"\\")) || strings.HasPrefix(relPath, basePath+"/") {
		hasPrefix = true
	}
	if hasPrefix {
		relativePath, err := filepath.Rel(basePath, relPath)
		if err != nil {
			return relPath, err
		}
		relPath = relativePath
	}
	return relPath, nil
}

func FindMountPoint(path string) (string, error) {
	current := path
	for {
		parent := filepath.Dir(current)
		if parent == current {
			return current, nil
		}

		devCurrent, _ := GetDevice(current)
		devParent, _ := GetDevice(parent)

		if devCurrent != devParent {
			return current, nil
		}
		current = parent
	}
}
